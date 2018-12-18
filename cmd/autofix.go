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

func mergeComments(commentFile, fixedFile string) ([]byte, error) {
	origData, err := ioutil.ReadFile(commentFile)
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

	// Merge fixed YAML and comments
	mergedYaml := mergeMapSlice(origYaml, fixedYaml)

	// Put back post-doc comments
	mergedYaml = append(mergedYaml, comments...)

	// Convert MapSlice to byte array
	data, err := yaml.Marshal(&mergedYaml)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Recursively merge fixedSlice into commentSlice.
// Keys which exist in commentSlice but not fixedSlice are removed from commentSlice.
// Keys which exist in fixedSlice but not commentSlice are added to commentSlice.
// If keys exist in both fixedSlice and commentSlice then the value in commentSlice is either replaced with the value
// in fixedSlice (if the value is simple) or, if the values are MapSlices, they are merged recursively.
func mergeMapSlice(commentSlice, fixedSlice yaml.MapSlice) yaml.MapSlice {
	// Remove items from the original which are not present in the fixed yaml
	for i := 0; i < len(commentSlice); i++ {
		item := commentSlice[i]
		if _, ok := item.Key.(yaml.PreDoc); item.Key == nil || ok {
			continue
		}
		if findKeyInYaml(item.Key, fixedSlice) == -1 {
			commentSlice = append(commentSlice[:i], commentSlice[i+1:]...)
			i--
		}
	}

	// Modify or add items from the fixed yaml which are not in the original
	for i := 0; i < len(fixedSlice); i++ {
		item := fixedSlice[i]
		index := findKeyInYaml(item.Key, commentSlice)
		if index == -1 {
			commentSlice = append(commentSlice, fixedSlice[i])
			continue
		} else if _, ok := item.Value.(yaml.MapSlice); ok {
			if _, ok = commentSlice[i].Value.(yaml.MapSlice); ok {
				commentSlice[i].Value = mergeMapSlice(commentSlice[i].Value.(yaml.MapSlice), item.Value.(yaml.MapSlice))
				continue
			}
		}
		commentSlice[i].Value = item.Value
	}

	return commentSlice

}

func findKeyInYaml(key interface{}, slice yaml.MapSlice) int {
	for i, item := range slice {
		if item.Key != nil && item.Key == key {
			return i
		}
	}
	return -1
}

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

	fixedYaml, err := mergeComments(rootConfig.manifest, tmpFile.Name())
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
