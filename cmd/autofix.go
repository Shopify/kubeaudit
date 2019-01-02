package cmd

import (
	"io/ioutil"
	"os"

	"github.com/Shopify/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getAuditFunctions() []interface{} {
	return []interface{}{
		auditAllowPrivilegeEscalation, auditReadOnlyRootFS, auditRunAsNonRoot,
		auditAutomountServiceAccountToken, auditPrivileged, auditCapabilities,
		auditAppArmor, auditSeccomp,
	}
}

func fixPotentialSecurityIssue(resource Resource, result Result) Resource {
	resource = prepareResourceForFix(resource, result)
	for _, occurrence := range result.Occurrences {
		switch occurrence.id {
		case ErrorAllowPrivilegeEscalationNil, ErrorAllowPrivilegeEscalationTrue:
			resource = fixAllowPrivilegeEscalation(resource, occurrence)
		case ErrorCapabilityNotDropped:
			resource = fixCapabilityNotDropped(resource, occurrence)
		case ErrorCapabilityAdded:
			resource = fixCapabilityAdded(resource, occurrence)
		case ErrorPrivilegedNil, ErrorPrivilegedTrue:
			resource = fixPrivileged(resource, occurrence)
		case ErrorReadOnlyRootFilesystemFalse, ErrorReadOnlyRootFilesystemNil:
			resource = fixReadOnlyRootFilesystem(resource, occurrence)
		case ErrorRunAsNonRootFalse, ErrorRunAsNonRootNil:
			resource = fixRunAsNonRoot(resource, occurrence)
		case ErrorServiceAccountTokenDeprecated:
			resource = fixDeprecatedServiceAccount(resource)
		case ErrorAutomountServiceAccountTokenTrueAndNoName, ErrorAutomountServiceAccountTokenNilAndNoName:
			resource = fixServiceAccountToken(resource)
		case ErrorAppArmorAnnotationMissing, ErrorAppArmorDisabled:
			resource = fixAppArmor(resource)
		case ErrorSeccompAnnotationMissing, ErrorSeccompDeprecated, ErrorSeccompDeprecatedPod, ErrorSeccompDisabled,
			ErrorSeccompDisabledPod:
			resource = fixSeccomp(resource)
		}
	}
	return resource
}

func prepareResourceForFix(resource Resource, result Result) Resource {
	needSecurityContextDefined := []int{ErrorAllowPrivilegeEscalationNil, ErrorAllowPrivilegeEscalationTrue,
		ErrorPrivilegedNil, ErrorPrivilegedTrue, ErrorReadOnlyRootFilesystemFalse, ErrorReadOnlyRootFilesystemNil,
		ErrorRunAsNonRootFalse, ErrorRunAsNonRootNil, ErrorServiceAccountTokenDeprecated,
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

func fix(resources []Resource) (fixedResources []Resource) {
	for _, resource := range resources {
		results := mergeAuditFunctions(getAuditFunctions())(resource)
		for _, result := range results {
			resource = fixPotentialSecurityIssue(resource, result)
		}
		fixedResources = append(fixedResources, resource)
	}
	return
}

// mergeYAML takes the file name of an autofixed YAML file (fixedFile) and the file name of the original YAML file
// (origFile) and merges fixedFile into origFile such that the resulting byte array is autofixed YAML but with the
// same order and comments as the original.
func mergeYAML(origFile, fixedFile string) ([]byte, error) {
	origData, err := ioutil.ReadFile(origFile)
	if err != nil {
		return nil, err
	}
	origYaml, err := yaml.CommentUnmarshal(origData)
	if err != nil {
		return nil, err
	}

	fixedData, err := ioutil.ReadFile(fixedFile)
	if err != nil {
		return nil, err
	}
	fixedYaml, err := yaml.CommentUnmarshal(fixedData)
	if err != nil {
		return nil, err
	}

	// Take out post-doc comments
	commentStart := len(origYaml)
	for origYaml[commentStart-1].Key == nil && len(origYaml[commentStart-1].Comment) > 0 {
		commentStart--
	}
	comments := make(yaml.MapSlice, 0, len(origYaml)-commentStart)
	comments = append(comments, origYaml[commentStart:]...)
	origYaml = origYaml[:commentStart]

	// Merge fixed YAML into original YAML
	mergedYaml := mergeMapSlices(origYaml, fixedYaml)

	// Put back post-doc comments
	mergedYaml = append(mergedYaml, comments...)

	// Convert YAML to byte array
	data, err := yaml.Marshal(&mergedYaml)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Recursively merge fixedSlice and origSlice.
// Keys which exist in origSlice but not fixedSlice are excluded.
// Keys which exist in fixedSlice but not origSlice are included.
// If keys exist in both fixedSlice and origSlice then the value from fixedSlice is used unless both values are complex
// (MapSlices or SequenceItem arrays), in which case they are merged recursively.
func mergeMapSlices(origSlice, fixedSlice yaml.MapSlice) yaml.MapSlice {
	var mergedSlice yaml.MapSlice

	// Keep comments, and items which are present in both the original and fixed yaml
	for _, item := range origSlice {
		if _, ok := item.Key.(yaml.PreDoc); item.Key == nil || ok || findKeyInMapSlice(item.Key, fixedSlice) != -1 {
			mergedSlice = append(mergedSlice, item)
			continue
		}
	}

	// Update or add items from the fixed yaml which are not in the original
	for _, fixedItem := range fixedSlice {
		mergedItemIndex := findKeyInMapSlice(fixedItem.Key, mergedSlice)
		if mergedItemIndex == -1 {
			mergedSlice = append(mergedSlice, fixedItem)
			continue
		}

		mergedItem := &mergedSlice[mergedItemIndex]
		if _, ok := fixedItem.Value.(yaml.MapSlice); ok {
			if _, ok = mergedItem.Value.(yaml.MapSlice); ok {
				mergedItem.Value = mergeMapSlices(mergedItem.Value.(yaml.MapSlice), fixedItem.Value.(yaml.MapSlice))
				continue
			}
		}
		if _, ok := fixedItem.Value.([]yaml.SequenceItem); ok {
			if _, ok = mergedItem.Value.([]yaml.SequenceItem); ok {
				sequenceKey := mergedItem.Key.(string)
				fixedSeq := fixedItem.Value.([]yaml.SequenceItem)
				origSeq := mergedItem.Value.([]yaml.SequenceItem)
				mergedItem.Value = mergeSequences(sequenceKey, origSeq, fixedSeq)
				continue
			}
		}
		mergedItem.Value = fixedItem.Value
	}

	return mergedSlice
}

// Returns the index of the MapItem in MapSlice which has the given key
func findKeyInMapSlice(key interface{}, slice yaml.MapSlice) int {
	for i, item := range slice {
		if item.Key != nil && item.Key == key {
			return i
		}
	}
	return -1
}

// Recursively merge fixedSlice and origSlice.
// Values which exist in origSlice but not fixedSlice are excluded.
// Values which exist in fixedSlice but not origSlice are included.
// If values exist in both fixedSlice and origSlice then the value from fixedSlice is used unless both values are
// complex (MapSlices or SequenceItem arrays), in which case they are merged recursively.
func mergeSequences(sequenceKey string, origSlice, fixedSlice []yaml.SequenceItem) []yaml.SequenceItem {
	var mergedSlice []yaml.SequenceItem

	// Keep comments, and items which are present in both the original and fixed yaml
	for _, item := range origSlice {
		if (item.Value == nil && len(item.Comment) > 0) || findItemInSequence(sequenceKey, item, fixedSlice) != -1 {
			mergedSlice = append(mergedSlice, item)
		}
	}

	// Update or add items from the fixed yaml which are not in the original
	for _, fixedItem := range fixedSlice {
		mergedItemIndex := findItemInSequence(sequenceKey, fixedItem, mergedSlice)
		if mergedItemIndex == -1 {
			mergedSlice = append(mergedSlice, fixedItem)
			continue
		}

		mergedItem := &mergedSlice[mergedItemIndex]
		if _, ok := fixedItem.Value.(yaml.MapSlice); ok {
			if _, ok = mergedItem.Value.(yaml.MapSlice); ok {
				mergedItem.Value = mergeMapSlices(mergedItem.Value.(yaml.MapSlice), fixedItem.Value.(yaml.MapSlice))
				continue
			}
		}
		mergedItem.Value = fixedItem.Value
	}

	return mergedSlice
}

// Returns the index of the item in slice which "matches" val. See sequenceItemMatch for what is considered a "match".
func findItemInSequence(sequenceKey string, val yaml.SequenceItem, slice []yaml.SequenceItem) int {
	for i, item := range slice {
		if item.Value != nil && sequenceItemMatch(sequenceKey, val, item) {
			return i
		}
	}
	return -1
}

var identifyingKey = map[string]string{
	"containers":    "name",          // Container
	"hostAliases":   "ip",            // HostAlias
	"env":           "name",          // EnvVar
	"ports":         "containerPort", // ContainerPort
	"volumeDevices": "name",          // VolumeDevice
	"volumeMounts":  "name",          // VolumeMount
}

// In order to determine whether sequence items match (and should be merged) we determine the "identifying key" for the
// sequence item, and if both sequence items have the same key-value pair for the "identifying key" then they are a match.
func sequenceItemMatch(sequenceKey string, item1, item2 yaml.SequenceItem) bool {
	if val1, ok := item1.Value.(string); ok {
		if val2, ok := item2.Value.(string); ok {
			return val1 == val2
		}
		return false
	}
	val1, ok := item1.Value.(yaml.MapSlice)
	if !ok {
		return false
	}
	val2, ok := item2.Value.(yaml.MapSlice)
	if !ok {
		return false
	}

	switch sequenceKey {
	// EnvFromSource
	case "envFrom":
		map1 := item1.Value.(yaml.MapSlice)
		map2 := item2.Value.(yaml.MapSlice)
		if findKeyInMapSlice("configMapRef", map1) != -1 && findKeyInMapSlice("configMapRef", map2) != -1 {
			return mapPairMatch("name", map1, map2)
		}
		if findKeyInMapSlice("secretRef", map1) != -1 && findKeyInMapSlice("secretRef", map2) != -1 {
			return mapPairMatch("name", map1, map2)
		}
		return false
	}

	return mapPairMatch(identifyingKey[sequenceKey], val1, val2)
}

// Returns true if map1 and map2 have the same key-value pair for the given key
// Assumes that the value at the given key is a primitive type
func mapPairMatch(key string, map1, map2 yaml.MapSlice) bool {
	index1 := findKeyInMapSlice(key, map1)
	if index1 == -1 {
		return false
	}
	index2 := findKeyInMapSlice(key, map2)
	if index2 == -1 {
		return false
	}
	return map1[index1].Value == map2[index2].Value
}

// The fix function does not preserve comments (because kubernetes resources do not support comments) so we convert
// both the original manifest file and the fixed manifest file into MapSlices (an array representation of a map which
// preserves the order of the keys) using the Shopify/yaml fork of go-yaml/yaml (the fork adds comment support) and
// then merge the fixed MapSlice back into the original MapSlice so that we get the comments back.
func autofix(*cobra.Command, []string) {
	resources, err := getKubeResourcesManifest(rootConfig.manifest)
	if err != nil {
		log.Error(err)
	}

	fixedResources := fix(resources)

	tmpFile, err := ioutil.TempFile("", "kubeaudit_autofix")
	if err != nil {
		log.Error(err)
	}
	defer os.Remove(tmpFile.Name())

	err = writeManifestFile(fixedResources, tmpFile.Name())
	if err != nil {
		log.Error(err)
	}

	fixedYaml, err := mergeYAML(rootConfig.manifest, tmpFile.Name())
	if err != nil {
		log.Error(err)
	}

	err = ioutil.WriteFile(rootConfig.manifest, fixedYaml, 0644)
	if err != nil {
		log.Error(err)
	}
}

var autofixCmd = &cobra.Command{
	Use:   "autofix",
	Short: "Automagically fixes a manifest to be secure",
	Long: `"autofix" will examine a manifest file and automagically fill in the blanks to leave your yaml file more secure than it found it

Example usage:
kubeaudit autofix -f /path/to/yaml`,
	Run: autofix,
}

func init() {
	RootCmd.AddCommand(autofixCmd)
}
