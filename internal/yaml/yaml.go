package yaml

import (
	"bytes"
	"fmt"

	log "github.com/sirupsen/logrus"
	goyaml "gopkg.in/yaml.v3"
)

// src https://github.com/go-yaml/yaml/blob/v3/resolve.go#L70
const (
	strTag = "!!str"
	seqTag = "!!seq"
	mapTag = "!!map"
)

// Merge merges the original YAML with the fixed YAML such that the resulting YAML is autofixed but with the
// same order and comments as the original.
func Merge(origData, fixedData []byte) ([]byte, error) {
	origYaml, err := unmarshal(origData)
	if err != nil {
		return nil, err
	}

	fixedYaml, err := unmarshal(fixedData)
	if err != nil {
		return nil, err
	}

	// Create a new document node that contains the merged maps for the original and fixed yaml
	mergedYaml := shallowCopyNode(origYaml)
	mergedYaml.Content = []*goyaml.Node{
		mergeMaps(origYaml.Content[0], fixedYaml.Content[0]),
	}

	return marshal(mergedYaml)
}

func unmarshal(data []byte) (*goyaml.Node, error) {
	var node goyaml.Node

	if err := goyaml.Unmarshal(data, &node); err != nil {
		return nil, err
	}
	if len(node.Content) != 1 {
		return nil, fmt.Errorf("expected original yaml document to have one child but got %v", len(node.Content))
	}
	if node.Content[0].Kind != goyaml.MappingNode {
		return nil, fmt.Errorf("expected mapping node as child of original yaml document node but got %v", node.Content[0].Kind)
	}

	return &node, nil
}

func marshal(node *goyaml.Node) ([]byte, error) {
	data := bytes.NewBuffer(nil)
	encoder := goyaml.NewEncoder(data)
	defer encoder.Close()
	encoder.SetIndent(2)
	if err := encoder.Encode(node); err != nil {
		return nil, fmt.Errorf("error marshaling merged yaml: %v", err)
	}
	return data.Bytes(), nil
}

// mergeMaps recursively merges orig and fixed.
// Key-value pairs which exist in orig but not fixed are excluded (as determined by matching the key).
// Key-value pairs which exist in fixed but not orig are included.
// If keys exist in both orig and fixed then the key-value pair from fixed is used unless both values are complex
// (maps or sequences), in which case they are merged recursively.
func mergeMaps(orig, fixed *goyaml.Node) *goyaml.Node {
	merged := shallowCopyNode(orig)
	origContent := orig.Content
	fixedContent := fixed.Content

	// Drop items from original if they are not in fixed
	for i := 0; i < len(origContent); i += 2 {
		origKey := origContent[i]
		if isKeyInMap(origKey, fixed) {
			origVal := origContent[i+1]
			merged.Content = append(merged.Content, origKey, origVal)
		}
	}

	// Update or add items from the fixed yaml which are not in the original
	for i := 0; i < len(fixedContent); i += 2 {
		fixedKey := fixedContent[i]
		fixedVal := fixedContent[i+1]
		if mergedKeyIndex := findKeyInMap(fixedKey, merged); mergedKeyIndex == -1 {
			// Add item
			merged.Content = append(merged.Content, fixedKey, fixedVal)
		} else {
			// Update item
			mergedValIndex := mergedKeyIndex + 1
			mergedVal := merged.Content[mergedValIndex]

			if fixedVal.Kind != mergedVal.Kind {
				merged.Content[mergedValIndex] = fixedVal
				continue
			}

			switch fixedVal.Kind {
			case goyaml.ScalarNode:
				merged.Content[mergedValIndex].Value = fixedVal.Value
			case goyaml.MappingNode:
				merged.Content[mergedValIndex] = mergeMaps(mergedVal, fixedVal)
			case goyaml.SequenceNode:
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
func mergeSequences(sequenceKey string, orig, fixed *goyaml.Node) *goyaml.Node {
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
			case fixedItem.Kind == goyaml.MappingNode:
				merged.Content[mergedItemIndex] = mergeMaps(mergedItem, fixedItem)
			case fixedItem.Kind == goyaml.SequenceNode:
				merged.Content[mergedItemIndex] = mergeSequences(sequenceKey, mergedItem, fixedItem)
			}
		}
	}
	return merged
}

// deepEqual recursively compares two values but ignores map and array child order and comments. For example the
// following values are considered to be equal:
//
//     []goyaml.SequenceItem{{Value: goyaml.MapSlice{
// 	       {Key: "k", Value: "v", Comment: "c"},
// 	       {Key: "k2", Value: "v2", Comment: "c2"},
//     }}}
//
//     []goyaml.SequenceItem{{Value: goyaml.MapSlice{
//          {Key: "k2", Value: "v2"},
//          {Key: "k", Value: "v"},
//      }}}
func deepEqual(val1, val2 *goyaml.Node) bool {
	if val1.Kind != val2.Kind {
		return false
	}

	switch val1.Kind {
	case goyaml.ScalarNode:
		return equalScalar(val1, val2)
	case goyaml.MappingNode:
		return equalMap(val1, val2)
	case goyaml.SequenceNode:
		return equalSequence(val1, val2)
	}

	return false
}

func equalScalar(val1, val2 *goyaml.Node) bool {
	if val1.Kind != goyaml.ScalarNode || val2.Kind != goyaml.ScalarNode {
		return false
	}
	return val1.Tag == val2.Tag && val1.Value == val2.Value
}

func equalSequence(seq1, seq2 *goyaml.Node) bool {
	if seq1.Kind != goyaml.SequenceNode || seq2.Kind != goyaml.SequenceNode {
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

func equalMap(map1, map2 *goyaml.Node) bool {
	if map1.Kind != goyaml.MappingNode || map2.Kind != goyaml.MappingNode {
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

// equalValueForKey returns true if map1 and map2 have the same key-value pair for the given key
func equalValueForKey(findKey string, map1, map2 *goyaml.Node) bool {
	if map1.Kind != goyaml.MappingNode || map2.Kind != goyaml.MappingNode {
		return false
	}

	if val1, index1 := findValInMap(findKey, map1); index1 != -1 {
		if val2, index2 := findValInMap(findKey, map2); index2 != -1 {
			return deepEqual(val1, val2)
		}
	}
	return false
}

// isKeyInMap returns true if findKey is a child of mapNode
func isKeyInMap(findKey *goyaml.Node, mapNode *goyaml.Node) bool {
	return findKeyInMap(findKey, mapNode) != -1
}

// findKeyInMap returns the index of findKey in mapNode's list of children, or -1 if it isn't found
func findKeyInMap(findKey *goyaml.Node, mapNode *goyaml.Node) int {
	if mapNode.Kind != goyaml.MappingNode {
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
func findValInMap(key string, mapNode *goyaml.Node) (*goyaml.Node, int) {
	findKey := &goyaml.Node{
		Kind:  goyaml.ScalarNode,
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
func isItemInSequence(sequenceKey string, findVal *goyaml.Node, sequenceNode *goyaml.Node) bool {
	return findItemInSequence(sequenceKey, findVal, sequenceNode) != -1
}

// findItemInSequence returns the index of the child in sequenceNode which "matches" findVal. See sequenceItemMatch for
// what is considered a "match". Returns -1 if there is no match found.
func findItemInSequence(sequenceKey string, findVal *goyaml.Node, sequenceNode *goyaml.Node) int {
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
	"containers":           "name",             // PodSpec.containers : Container.name
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
func sequenceItemMatch(sequenceKey string, item1, item2 *goyaml.Node) bool {
	if item1.Kind != item2.Kind {
		return false
	}

	if sequenceKey == "" || item1.Kind != goyaml.MappingNode {
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
		// ProjectedVolumeSource.sources : VolumeProjection.serviceAccountToken.path
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

// Returns a new *goyaml.Node with the same values as the original node except for the Content field, which is initialized
// to an empty array
func shallowCopyNode(orig *goyaml.Node) *goyaml.Node {
	return &goyaml.Node{
		Kind:        orig.Kind,
		Style:       orig.Style,
		Tag:         orig.Tag,
		Value:       orig.Value,
		Anchor:      orig.Anchor,
		Alias:       orig.Alias,
		Content:     []*goyaml.Node{},
		HeadComment: orig.HeadComment,
		LineComment: orig.LineComment,
		FootComment: orig.FootComment,
		Line:        orig.Line,
		Column:      orig.Column,
	}
}
