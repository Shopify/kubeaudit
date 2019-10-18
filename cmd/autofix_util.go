package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Shopify/kubeaudit/scheme"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

// src https://github.com/go-yaml/yaml/blob/v3/resolve.go#L70
const (
	nullTag      = "!!null"
	boolTag      = "!!bool"
	strTag       = "!!str"
	intTag       = "!!int"
	floatTag     = "!!float"
	timestampTag = "!!timestamp"
	seqTag       = "!!seq"
	mapTag       = "!!map"
	binaryTag    = "!!binary"
	mergeTag     = "!!merge"
)

func getAuditFunctions() []interface{} {
	return []interface{}{
		auditAllowPrivilegeEscalation, auditReadOnlyRootFS, auditRunAsNonRoot,
		auditAutomountServiceAccountToken, auditPrivileged, auditCapabilities,
		auditAppArmor, auditSeccomp, auditNetworkPolicies, auditNamespaces,
	}
}

func fixPotentialSecurityIssue(resource Resource, result Result) Resource {
	resource = prepareResourceForFix(resource, result)

	for _, occurrence := range result.Occurrences {
		switch occurrence.id {
		case ErrorAllowPrivilegeEscalationNil, ErrorAllowPrivilegeEscalationTrue:
			resource = fixAllowPrivilegeEscalation(&result, resource, occurrence)
		case ErrorCapabilityNotDropped:
			resource = fixCapabilityNotDropped(&result, resource, occurrence)
		case ErrorCapabilityAdded:
			resource = fixCapabilityAdded(&result, resource, occurrence)
		case ErrorPrivilegedNil, ErrorPrivilegedTrue:
			resource = fixPrivileged(&result, resource, occurrence)
		case ErrorReadOnlyRootFilesystemFalse, ErrorReadOnlyRootFilesystemNil:
			resource = fixReadOnlyRootFilesystem(&result, resource, occurrence)
		case ErrorRunAsNonRootPSCTrueFalseCSCFalse, ErrorRunAsNonRootPSCNilCSCNil, ErrorRunAsNonRootPSCFalseCSCNil:
			resource = fixRunAsNonRoot(&result, resource, occurrence)
		case ErrorServiceAccountTokenDeprecated:
			resource = fixDeprecatedServiceAccount(resource)
		case ErrorAutomountServiceAccountTokenTrueAndNoName, ErrorAutomountServiceAccountTokenNilAndNoName:
			resource = fixServiceAccountToken(&result, resource)
		case ErrorAppArmorAnnotationMissing, ErrorAppArmorDisabled:
			resource = fixAppArmor(resource)
		case ErrorSeccompAnnotationMissing, ErrorSeccompDeprecated, ErrorSeccompDeprecatedPod, ErrorSeccompDisabled,
			ErrorSeccompDisabledPod:
			resource = fixSeccomp(resource)
		case ErrorMissingDefaultDenyIngressNetworkPolicy, ErrorMissingDefaultDenyEgressNetworkPolicy, ErrorMissingDefaultDenyIngressAndEgressNetworkPolicy:
			resource = fixNetworkPolicy(resource, occurrence)
		case ErrorNamespaceHostIPCTrue, ErrorNamespaceHostNetworkTrue, ErrorNamespaceHostPIDTrue:
			resource = fixNamespace(&result, resource)
		}
	}
	return resource
}

func prepareResourceForFix(resource Resource, result Result) Resource {
	needSecurityContextDefined := []int{ErrorAllowPrivilegeEscalationNil, ErrorAllowPrivilegeEscalationTrue,
		ErrorPrivilegedNil, ErrorPrivilegedTrue, ErrorReadOnlyRootFilesystemFalse, ErrorReadOnlyRootFilesystemNil,
		ErrorRunAsNonRootPSCTrueFalseCSCFalse, ErrorRunAsNonRootPSCNilCSCNil, ErrorRunAsNonRootPSCFalseCSCNil, ErrorServiceAccountTokenDeprecated,
		ErrorAutomountServiceAccountTokenTrueAndNoName, ErrorAutomountServiceAccountTokenNilAndNoName,
		ErrorCapabilityNotDropped, ErrorCapabilityAdded, ErrorMisconfiguredKubeauditAllow}
	needCapabilitiesDefined := []int{ErrorCapabilityNotDropped, ErrorCapabilityAdded, ErrorMisconfiguredKubeauditAllow}

	// Set of errors to fix
	errors := make(map[int]bool)
	for _, occurrence := range result.Occurrences {
		errors[occurrence.id] = true
	}

	for _, err := range needSecurityContextDefined {
		if _, ok := errors[err]; ok {
			resource = fixSecurityContextNil(resource)
			break
		}
	}

	for _, err := range needCapabilitiesDefined {
		if _, ok := errors[err]; ok {
			resource = fixCapabilitiesNil(resource)
			break
		}
	}

	return resource
}

func fix(resources []Resource) (fixedResources []Resource, extraResources []Resource) {
	for _, resource := range resources {
		if !IsSupportedResourceType(resource) {
			fixedResources = append(fixedResources, resource)
			continue
		}
		results := mergeAuditFunctions(getAuditFunctions())(resource)
		for _, result := range results {
			if IsNamespaceType(resource) {
				extraResource := fixPotentialSecurityIssue(resource, result)
				// If return resource from fixPotentialSecurityIssue is Namespace type then we don't have to add extra resources for it.
				if !IsNamespaceType(extraResource) {
					extraResources = append(extraResources, extraResource)
				}
			} else {
				resource = fixPotentialSecurityIssue(resource, result)
			}
		}
		fixedResources = append(fixedResources, resource)
	}
	return
}

// deepEqual recursively compares two values but ignores map and array child order and comments. For example the
// following values are considered to be equal:
//
//     []yaml.SequenceItem{{Value: yaml.MapSlice{
// 	       {Key: "k", Value: "v", Comment: "c"},
// 	       {Key: "k2", Value: "v2", Comment: "c2"},
//     }}}
//
//     []yaml.SequenceItem{{Value: yaml.MapSlice{
//          {Key: "k2", Value: "v2"},
//          {Key: "k", Value: "v"},
//      }}}
func deepEqual(val1, val2 *yaml.Node) bool {
	if val1.Kind != val2.Kind {
		return false
	}

	switch val1.Kind {
	case yaml.ScalarNode:
		return equalScalar(val1, val2)
	case yaml.MappingNode:
		return equalMap(val1, val2)
	case yaml.SequenceNode:
		return equalSequence(val1, val2)
	}

	return false
}

func equalScalar(val1, val2 *yaml.Node) bool {
	if val1.Kind != yaml.ScalarNode || val2.Kind != yaml.ScalarNode {
		return false
	}
	return val1.Tag == val2.Tag && val1.Value == val2.Value
}

func equalSequence(seq1, seq2 *yaml.Node) bool {
	if seq1.Kind != yaml.SequenceNode || seq2.Kind != yaml.SequenceNode {
		return false
	}

	content1 := seq1.Content
	content2 := seq2.Content
	if len(content1) != len(content2) {
		return false
	}

	for _, val1 := range content1 {
		if !isItemInSequence("", val1, seq2) {
			return false
		}
	}

	return true
}

func equalMap(map1, map2 *yaml.Node) bool {
	if map1.Kind != yaml.MappingNode || map2.Kind != yaml.MappingNode {
		return false
	}

	content1 := map1.Content
	content2 := map2.Content
	if len(content1) != len(content2) {
		return false
	}

	for i := 0; i < len(content1); i += 2 {
		key1 := content1[i]
		index2 := findKeyInMap(key1, map2)
		if index2 == -1 {
			return false
		}

		value1 := content1[i+1]
		value2 := content2[index2+1]
		if !deepEqual(value1, value2) {
			return false
		}
	}

	return true
}

// isFirstLineSeparator returns true if the first line in the manifest file is a yaml separator
func isFirstLineSeparatorOrComment(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineStr := scanner.Text()
		if lineStr == "---" || lineStr[0] == '#' {
			return true
		}
		return false
	}
	return false
}

// isCommentSlice returns true if the byteslice contains only yaml comments
func isCommentSlice(b []byte) bool {
	lineSlice := bytes.Split(b, []byte("/n"))
	for _, line := range lineSlice {
		if len(line) > 0 && !strings.HasPrefix(string(line), "#") {
			return false
		}
	}
	return true
}

// equalValueForKey returns true if map1 and map2 have the same key-value pair for the given key
func equalValueForKey(findKey string, map1, map2 *yaml.Node) bool {
	if map1.Kind != yaml.MappingNode || map2.Kind != yaml.MappingNode {
		return false
	}

	if val1, index1 := findValInMap(findKey, map1); index1 != -1 {
		if val2, index2 := findValInMap(findKey, map2); index2 != -1 {
			return deepEqual(val1, val2)
		}
	}
	return false
}

// mergeYAML takes the file name of an autofixed YAML file (fixedFile) and the file name of the original YAML file
// (origFile) and merges fixedFile into origFile such that the resulting byte array is autofixed YAML but with the
// same order and comments as the original.
func mergeYAML(origFile, fixedFile string) ([]byte, error) {
	var origYaml yaml.Node
	var fixedYaml yaml.Node

	origData, err := ioutil.ReadFile(origFile)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(origData, &origYaml); err != nil {
		return nil, err
	}
	if len(origYaml.Content) != 1 {
		return nil, fmt.Errorf("expected original yaml document to have one child but got %v", len(origYaml.Content))
	}
	if origYaml.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected mapping node as child of original yaml document node but got %v", origYaml.Content[0].Kind)
	}

	fixedData, err := ioutil.ReadFile(fixedFile)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(fixedData, &fixedYaml); err != nil {
		return nil, err
	}
	if len(fixedYaml.Content) != 1 {
		return nil, fmt.Errorf("expected fixed yaml document to have one child but got %v", len(fixedYaml.Content))
	}
	if fixedYaml.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected mapping node as child of fixed yaml document node but got %v", fixedYaml.Content[0].Kind)
	}

	// Create a new document node that contains the merged maps for the original and fixed yaml
	mergedYaml := shallowCopyNode(&origYaml)
	mergedYaml.Content = []*yaml.Node{
		mergeMaps(origYaml.Content[0], fixedYaml.Content[0]),
	}

	// Convert YAML to byte array
	data := bytes.NewBuffer(nil)
	encoder := yaml.NewEncoder(data)
	encoder.SetIndent(2)
	defer encoder.Close()
	if err = encoder.Encode(mergedYaml); err != nil {
		return nil, fmt.Errorf("error marshaling merged yaml: %v", err)
	}

	return data.Bytes(), nil
}

// mergeMaps recursively merges orig and fixed.
// Key-value pairs which exist in orig but not fixed are excluded (as determined by matching the key).
// Key-value pairs which exist in fixed but not orig are included.
// If keys exist in both orig and fixed then the key-value pair from fixed is used unless both values are complex
// (maps or sequences), in which case they are merged recursively.
func mergeMaps(orig, fixed *yaml.Node) *yaml.Node {
	merged := shallowCopyNode(orig)
	origContent := orig.Content
	fixedContent := fixed.Content

	// Drop items from original if they are not in fixed
	for i := 0; i < len(origContent); i += 2 {
		origKey := origContent[i]
		if isKeyInMap(origKey, fixed) {
			origVal := origContent[i+1]
			merged.Content = append(merged.Content, origKey)
			merged.Content = append(merged.Content, origVal)
		}
	}

	// Update or add items from the fixed yaml which are not in the original
	for i := 0; i < len(fixedContent); i += 2 {
		fixedKey := fixedContent[i]
		fixedVal := fixedContent[i+1]
		if mergedKeyIndex := findKeyInMap(fixedKey, merged); mergedKeyIndex == -1 {
			// Add item
			merged.Content = append(merged.Content, fixedKey)
			merged.Content = append(merged.Content, fixedVal)
		} else {
			// Update item
			mergedValIndex := mergedKeyIndex + 1
			mergedVal := merged.Content[mergedValIndex]

			if fixedVal.Kind != mergedVal.Kind {
				merged.Content[mergedValIndex] = fixedVal
				continue
			}

			switch fixedVal.Kind {
			case yaml.ScalarNode:
				merged.Content[mergedValIndex].Value = fixedVal.Value
			case yaml.MappingNode:
				merged.Content[mergedValIndex] = mergeMaps(mergedVal, fixedVal)
			case yaml.SequenceNode:
				merged.Content[mergedValIndex] = mergeSequences(fixedKey.Value, mergedVal, fixedVal)
			default:
				log.Error("Unexpected yaml node kind", fixedVal.Kind)
			}
		}
	}

	return merged
}

// mergeSequences recursively merges orig and fixed.
// Items which exist in orig but not fixed are excluded.
// Items which exist in fixed but not orig are included.
// If items exist in both orig and fixed then the item from fixed is used unless both items are complex
// (maps or sequences), in which case they are merged recursively.
func mergeSequences(sequenceKey string, orig, fixed *yaml.Node) *yaml.Node {
	merged := shallowCopyNode(orig)
	origContent := orig.Content
	fixedContent := fixed.Content

	// Drop items from original if they are not in fixed
	for _, origItem := range origContent {
		if isItemInSequence(sequenceKey, origItem, fixed) {
			merged.Content = append(merged.Content, origItem)
		}
	}

	// Update or add items from the fixed yaml which are not in the original
	for _, fixedItem := range fixedContent {
		if mergedItemIndex := findItemInSequence(sequenceKey, fixedItem, merged); mergedItemIndex == -1 {
			// Add item
			merged.Content = append(merged.Content, fixedItem)
		} else {
			// Update item
			mergedItem := merged.Content[mergedItemIndex]
			switch {
			case fixedItem.Kind != mergedItem.Kind:
				merged.Content[mergedItemIndex] = fixedItem
			case fixedItem.Kind == yaml.MappingNode:
				merged.Content[mergedItemIndex] = mergeMaps(mergedItem, fixedItem)
			case fixedItem.Kind == yaml.SequenceNode:
				merged.Content[mergedItemIndex] = mergeSequences(sequenceKey, mergedItem, fixedItem)
			}
		}
	}
	return merged
}

// isKeyInMap returns true if findKey is a child of mapNode
func isKeyInMap(findKey *yaml.Node, mapNode *yaml.Node) bool {
	return findKeyInMap(findKey, mapNode) != -1
}

// findKeyInMap returns the index of findKey in mapNode's list of children, or -1 if it isn't found
func findKeyInMap(findKey *yaml.Node, mapNode *yaml.Node) int {
	if mapNode.Kind != yaml.MappingNode {
		return -1
	}

	children := mapNode.Content
	for i := 0; i < len(children); i += 2 {
		key := children[i]
		if deepEqual(key, findKey) {
			return i
		}
	}

	return -1
}

// findValInMap returns the child of mapNode which is value corresponding to the given key, and its index
func findValInMap(key string, mapNode *yaml.Node) (*yaml.Node, int) {
	findKey := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: key,
		Tag:   strTag,
	}

	keyIndex := findKeyInMap(findKey, mapNode)
	if keyIndex == -1 {
		return nil, -1
	}

	valIndex := keyIndex + 1
	return mapNode.Content[valIndex], valIndex
}

// isItemInSequence returns true if findVal is a child of sequenceNode
func isItemInSequence(sequenceKey string, findVal *yaml.Node, sequenceNode *yaml.Node) bool {
	return findItemInSequence(sequenceKey, findVal, sequenceNode) != -1
}

// findItemInSequence returns the index of the child in sequenceNode which "matches" findVal. See sequenceItemMatch for
// what is considered a "match". Returns -1 if there is no match found.
func findItemInSequence(sequenceKey string, findVal *yaml.Node, sequenceNode *yaml.Node) int {
	children := sequenceNode.Content
	for i, val := range children {
		if sequenceItemMatch(sequenceKey, val, findVal) {
			return i
		}
	}
	return -1
}

var identifyingKey = map[string]string{
	"allowedFlexVolumes": "driver",     // PodSecurityPolicySpec.allowedFlexVolumes : AllowedFlexVolume.driver
	"allowedHostPaths":   "pathPrefix", // PodSecurityPolicySpec.allowedHostPaths : AllowedHostPath.pathPrefix
	// StorageClass.allowedTopologies : TopologySelectorTerm.matchLabelExpressions
	"allowedTopologies":    "matchLabelExpressions",
	"clusterRoleSelectors": "matchExpressions", // AggregationRule.clusterRoleSelectors : LabelSelector.matchExpressions
	"containers":           "name",             // PodSpec.contaienrs : Container.name
	"egress":               "ports",            // NetworkPolicySpec.egress : NetworkPolicyEgressRule.ports
	"env":                  "name",             // Container.env : EnvVar.name
	"hostAliases":          "ip",               // PodSpec.hostAliases : HostAlias.ip
	// Assumes it is not possible to add multiple values for the same header, ie.
	//     httpHeaders:
	//         - name: header1
	//           value: value1
	//         - name: header1
	//           value: value2
	// This restriction is not documented so the assumption may be incorrect
	// See https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#httpheader-v1-core
	"httpHeaders": "name", // HTTPGetAction.httpHeaders : HTTPHeader.name
	// PodSpec.imagePullSecrets : LocalObjectReference.name
	// ServiceAccount.imagePullSecrets : LocalObjectReference.name
	"imagePullSecrets": "name",
	"initContainers":   "name", // PodSpec.initContainers : Container.name
	// LabelSelector.matchExpressions : LabelSelectorRequirement.key
	// NodeSelectorTerm.matchExpressions : NodeSelectorRequirement.key
	"matchExpressions": "key",
	"matchFields":      "key",  // NodeSelectorTerm.matchFields : NodeSelectorRequirement.key
	"options":          "name", // PodDNSConfig.options : PodDNSConfigOption.name
	// TopologySelectorTerm.matchLabelExpressions : TopologySelectorLabelRequirement.key
	"matchLabelExpressions": "key",
	"pending":               "name",          // Initializers.pending : Initializer.name
	"readinessGates":        "conditionType", // PodSpec.readinessGates : PodReadinessGate.conditionType
	// PodAffinity.requiredDuringSchedulingIgnoredDuringExecution : PodAffinityTerm.labelSelector
	// PodAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution : PodAffinityTerm.labelSelector
	"requiredDuringSchedulingIgnoredDuringExecution": "labelSelector",
	"secrets": "name", // ServiceAccount.secrets : ObjectReference.name
	// ClusterRoleBinding.subjects : Subject.name
	// RoleBinding.subjects : Subject.name
	"subjects":      "name",
	"subsets":       "addresses",  // Endpoints.subsets : EndpointSubset.addresses
	"sysctls":       "name",       // PodSecurityContext.sysctls : Sysctl.name
	"taints":        "key",        // NodeSpec.taints : Taint.key
	"volumeDevices": "devicePath", // Container.volumeDevices : VolumeDevice.devicePath
	"volumeMounts":  "mountPath",  // Container.volumeMounts : VolumeMount.mountPath
	"volumes":       "name",       // PodSpec.volumes : Volume.name
}

// sequenceItemMatch returns true if item1 and item2 are a match, false otherwise. In order to determine whether
// sequence items match (and should be merged) we determine the "identifying key" for the sequence item, and if both
// sequence items have the same key-value pair for the "identifying key" then they are a match. The sequenceKey
// is the key for which the array items are the value. ie:
//     sequenceKey:
//     - item1
//     - item2
func sequenceItemMatch(sequenceKey string, item1, item2 *yaml.Node) bool {
	if item1.Kind != item2.Kind {
		return false
	}

	if sequenceKey == "" || item1.Kind != yaml.MappingNode {
		return deepEqual(item1, item2)
	}

	if idKey, ok := identifyingKey[sequenceKey]; ok {
		return equalValueForKey(idKey, item1, item2)
	}

	switch sequenceKey {
	// EndpointSubset.addresses : EndpointAddress.[hostname OR ip]
	// EndpointSubset.notReadyAddresses : EndpointAddress.[hostname OR ip]
	case "addresses", "notReadyAddresses":
		if equalValueForKey("hostname", item1, item2) {
			return true
		}
		return equalValueForKey("ip", item1, item2)

	// Container.envFrom : EnvFromSource.[configMapRef OR secretRef].name
	case "envFrom":
		if val1, index1 := findValInMap("configMapRef", item1); index1 != -1 {
			if val2, index2 := findValInMap("configMapRef", item2); index2 != -1 {
				return equalValueForKey("name", val1, val2)
			}
		}
		if val1, index1 := findValInMap("secretRef", item1); index1 != -1 {
			if val2, index2 := findValInMap("secretRef", item2); index2 != -1 {
				return equalValueForKey("name", val1, val2)
			}
		}
		return false

	// NetworkPolicySpec.ingress : NetworkPolicyIngressRule.[ports OR from]
	case "ingress":
		if equalValueForKey("ports", item1, item2) {
			return true
		}
		return equalValueForKey("from", item1, item2)

	// ConfigMapProjection.items : KeyToPath.key
	// ConfigMapVolumeSource.items : KeyToPath.key
	// DownwardAPIVolumeSource.items : DownwardAPIVolumeFile.path
	// SecretSecretProjection.items : KeyToPath.key
	// SecretVolumeSource.items : KeyToPath.key
	case "items":
		// ConfigMapVolumeSource.items : KeyToPath.key
		// SecretVolumeSource.items : KeyToPath.key
		if equalValueForKey("key", item1, item2) {
			return true
		}
		// DownwardAPIVolumeSource.items : DownwardAPIVolumeFile.path
		return equalValueForKey("path", item1, item2)

	// NodeSelector.nodeSelectorTerms : NodeSelectorTerm.[matchExpressions OR matchFields]
	case "nodeSelectorTerms":
		// This is a bit of a complicated case.
		// See https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#nodeselector-v1-core
		// For now, only match if there is an exact match for the complex value of either the "matchExpressions" or
		// "matchFields" fields.
		if equalValueForKey("matchExpressions", item1, item2) {
			return true
		}
		return equalValueForKey("matchFields", item1, item2)

	// ObjectMeta.ownerReferences : OwnerReference.[uid OR name]
	case "ownerReferences":
		if equalValueForKey("uid", item1, item2) {
			return true
		}
		return equalValueForKey("name", item1, item2)

	// NodeAffinity.preferredDuringSchedulingIgnoredDuringExecution : PreferredSchedulingTerm.preference
	// PodAffinity.preferredDuringSchedulingIgnoredDuringExecution : WeightedPodAffinityTerm.podAffinityTerm
	// PodAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution : WeightedPodAffinityTerm.podAffinityTerm
	case "preferredDuringSchedulingIgnoredDuringExecution":
		// This is a bit of a complicated case as the values are very nested and because the same identifying key is
		// used for two different array types (PreferredSchedulingTerm and WeightedPodAffinityTerm).
		// See https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#nodeaffinity-v1-core
		// and https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#podaffinity-v1-core
		// For now, only match if there is an exact match for the complex value of the "preference" or
		// "podAffinityTerm" field.
		// The value for the "weight" field can be updated.

		// NodeAffinity.preferredDuringSchedulingIgnoredDuringExecution : PreferredSchedulingTerm.preference
		if equalValueForKey("preference", item1, item2) {
			return true
		}
		// PodAffinity.preferredDuringSchedulingIgnoredDuringExecution : WeightedPodAffinityTerm.podAffinityTerm
		// PodAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution : WeightedPodAffinityTerm.podAffinityTerm
		return equalValueForKey("podAffinityTerm", item1, item2)

	// Container.ports : ContainerPort.containerPort
	// EndpointSubset.ports : EndpointPort.port
	// ServiceSpec.ports : ServicePort.port
	case "ports":
		// Container.ports : ContainerPort.containerPort
		if equalValueForKey("containerPort", item1, item2) {
			return true
		}
		// EndpointSubset.ports : EndpointPort.port
		// ServiceSpec.ports : ServicePort.port
		return equalValueForKey("port", item1, item2)

	// ClusterRole.rules : PolicyRule.resources
	// IngressSpec.rules : IngressRule.host
	// Role.rules : PolicyRule.resources
	case "rules":
		// ClusterRole.rules : PolicyRule.resources
		// Role.rules : PolicyRule.resources
		if equalValueForKey("resources", item1, item2) {
			return true
		}
		// IngressSpec.rules : IngressRule.host
		if equalValueForKey("host", item1, item2) {
			return true
		}
		return deepEqual(item1, item2)

	// ProjectedVolumeSource.sources
	case "sources":
		// ProjectedVolumeSource.sources : VolumeProjection.configMap.name
		if val1, index1 := findValInMap("configMap", item1); index1 != -1 {
			if val2, index2 := findValInMap("configMap", item2); index2 != -1 {
				return equalValueForKey("name", val1, val2)
			}
			return false
		}
		// ProjectedVolumeSource.sources : VolumeProjection.downwardAPI.items
		if val1, index1 := findValInMap("downwardAPI", item1); index1 != -1 {
			if val2, index2 := findValInMap("downwardAPI", item2); index2 != -1 {
				return equalValueForKey("items", val1, val2)
			}
			return false
		}
		// ProjectedVolumeSource.sources : VolumeProjection.secret.name
		if val1, index1 := findValInMap("secret", item1); index1 != -1 {
			if val2, index2 := findValInMap("secret", item2); index2 != -1 {
				return equalValueForKey("name", val1, val2)
			}
			return false
		}
		// ProjectedVolumeSource.sources : VolumeProjection.serviceAccountToken.name
		if val1, index1 := findValInMap("serviceAccountToken", item1); index1 != -1 {
			if val2, index2 := findValInMap("serviceAccountToken", item2); index2 != -1 {
				return equalValueForKey("path", val1, val2)
			}
		}
		return false

	// IngressSpec.tls : IngressTLS.[secretName OR hosts]
	case "tls":
		if equalValueForKey("secretName", item1, item2) {
			return true
		}
		return equalValueForKey("hosts", item1, item2)

	// StatefulSetSpec.volumeClaimTemplates : PersistentVolumeClaim.metadata.name
	case "volumeClaimTemplates":
		if val1, index1 := findValInMap("metadata", item1); index1 != -1 {
			if val2, index2 := findValInMap("metadata", item2); index2 != -1 {
				return equalValueForKey("name", val1, val2)
			}
		}
		return false
	}

	// FSGroupStrategyOptions.ranges : IDRange
	// RunAsGroupStrategyOptions.ranges : IDRange
	// RunAsUserStrategyOptions.ranges : IDRange
	// SupplementalGroupsStrategyOptions.ranges : IDRange
	// PodSecurityPolicySpec.hostPorts : HostPortRange
	// PodSpec.tolerations : Toleration
	return deepEqual(item1, item2)
}

// SplitYamlResource splits the yaml file into byte slices for each resource in the yaml file and checks if the first resource
// is only comments, in which case it deletes the first resource in the slice and adds the comment to the final file and updates toAppend flag
func splitYamlResources(filename string, toWriteFile string) (splitDecoded [][]byte, toAppend bool, err error) {
	buf, err := ioutil.ReadFile(rootConfig.manifest)

	if err != nil {
		log.Error("File not found")
		return
	}
	splitDecoded = bytes.Split(buf, []byte("---"))
	if err != nil {
		log.Error(err)
		return nil, false, err
	}
	if len(splitDecoded) != 0 {
		if len(splitDecoded[0]) == 0 {
			splitDecoded = splitDecoded[1:]
		}
		decoder := scheme.Codecs.UniversalDeserializer()
		_, _, err := decoder.Decode(splitDecoded[0], nil, nil)
		// if Decode returns err, then it means that splitDecoded[0] is only comments(pre doc) in this case remove this resource from slice and write to file
		if err != nil {
			err = writeManifestFile(splitDecoded[0], toWriteFile, false)
			if err != nil {
				return nil, false, err
			}
			splitDecoded = splitDecoded[1:]
			return splitDecoded, true, nil
		}
	}

	return splitDecoded, false, nil
}

func cleanupManifest(origFile string, finalData []byte) ([]byte, error) {
	objectMetacreationTs := []byte("\n  creationTimestamp: null\n")
	specTemplatecreationTs := []byte("\n      creationTimestamp: null\n")
	jobSpecTemplatecreationTs := []byte("\n          creationTimestamp: null\n")
	nullStatus := []byte("\nstatus: {}\n")
	nullReplicaStatus := []byte("status:\n  replicas: 0\n")
	nullLBStatus := []byte("status:\n  loadBalancer: {}\n")
	nullMetaStatus := []byte("\n    status: {}\n")

	var hasObjectMetacreationTs, hasSpecTemplatecreationTs, hasJobSpecTemplatecreationTs, hasNullStatus,
		hasNullReplicaStatus, hasNullLBStatus, hasNullMetaStatus bool

	if origFile != "" {
		origData, err := ioutil.ReadFile(origFile)
		if err != nil {
			return nil, err
		}
		hasObjectMetacreationTs = bytes.Contains(origData, objectMetacreationTs)
		hasSpecTemplatecreationTs = bytes.Contains(origData, specTemplatecreationTs)
		hasJobSpecTemplatecreationTs = bytes.Contains(origData, jobSpecTemplatecreationTs)

		hasNullStatus = bytes.Contains(origData, nullStatus)
		hasNullReplicaStatus = bytes.Contains(origData, nullReplicaStatus)
		hasNullLBStatus = bytes.Contains(origData, nullLBStatus)
		hasNullMetaStatus = bytes.Contains(origData, nullMetaStatus)

	} // null value is false in case of origFile

	if !hasObjectMetacreationTs {
		finalData = bytes.Replace(finalData, objectMetacreationTs, []byte("\n"), -1)
	}
	if !hasSpecTemplatecreationTs {
		finalData = bytes.Replace(finalData, specTemplatecreationTs, []byte("\n"), -1)
	}
	if !hasJobSpecTemplatecreationTs {
		finalData = bytes.Replace(finalData, jobSpecTemplatecreationTs, []byte("\n"), -1)
	}
	if !hasNullStatus {
		finalData = bytes.Replace(finalData, nullStatus, []byte("\n"), -1)
	}
	if !hasNullReplicaStatus {
		finalData = bytes.Replace(finalData, nullReplicaStatus, []byte("\n"), -1)
	}
	if !hasNullLBStatus {
		finalData = bytes.Replace(finalData, nullLBStatus, []byte("\n"), -1)
	}
	if !hasNullMetaStatus {
		finalData = bytes.Replace(finalData, nullMetaStatus, []byte("\n"), -1)
	}

	return finalData, nil
}

// Returns a new *yaml.Node with the same values as the original node except for the Content field, which is initialized
// to an empty array
func shallowCopyNode(orig *yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:        orig.Kind,
		Style:       orig.Style,
		Tag:         orig.Tag,
		Value:       orig.Value,
		Anchor:      orig.Anchor,
		Alias:       orig.Alias,
		Content:     []*yaml.Node{},
		HeadComment: orig.HeadComment,
		LineComment: orig.LineComment,
		FootComment: orig.FootComment,
		Line:        orig.Line,
		Column:      orig.Column,
	}
}
