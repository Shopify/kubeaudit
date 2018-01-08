package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

func getAuditFunctions() []interface{} {
	return []interface{}{
		auditAllowPrivilegeEscalation, auditReadOnlyRootFS, auditRunAsNonRoot,
		auditAutomountServiceAccountToken, auditPrivileged, auditCapabilities,
	}
}

func runAllAudits(resource k8sRuntime.Object) (results []Result) {
	for _, function := range getAuditFunctions() {
		for _, result := range getResults([]k8sRuntime.Object{resource}, function) {
			results = append(results, result)
		}
	}
	return results
}

func fixPotentialSecurityIssue(resource k8sRuntime.Object, result Result) k8sRuntime.Object {
	resource = fixSecurityContextNIL(resource)
	resource = fixCapabilitiesNIL(resource)
	for _, occurrence := range result.Occurrences {
		switch occurrence.id {
		case ErrorAllowPrivilegeEscalationNIL, ErrorAllowPrivilegeEscalationTrue:
			resource = fixAllowPrivilegeEscalation(resource)
		case ErrorCapabilityNotDropped:
			resource = fixCapabilityNotDropped(resource, occurrence)
		case ErrorCapabilityAdded:
			resource = fixCapabilityAdded(resource, occurrence)
		case ErrorPrivilegedNIL, ErrorPrivilegedTrue:
			resource = fixPrivilegeEscalation(resource)
		case ErrorReadOnlyRootFilesystemFalse, ErrorReadOnlyRootFilesystemNIL:
			resource = fixReadOnlyRootFilesystem(resource)
		case ErrorRunAsNonRootFalse, ErrorRunAsNonRootNIL:
			resource = fixRunAsNonRoot(resource)
		case ErrorServiceAccountTokenDeprecated:
			resource = fixDeprecatedServiceAccount(resource)
		case ErrorAutomountServiceAccountTokenTrueAndNoName, ErrorAutomountServiceAccountTokenNILAndNoName:
			resource = fixServiceAccountToken(resource)
		}
	}
	return resource
}

func fix(resources []k8sRuntime.Object) (fixedResources []k8sRuntime.Object) {
	for _, resource := range resources {
		results := runAllAudits(resource)
		for _, result := range results {
			resource = fixPotentialSecurityIssue(resource, result)
		}
		fixedResources = append(fixedResources, resource)
	}
	return
}

func autofix(*cobra.Command, []string) {
	resources, err := getKubeResourcesManifest(rootConfig.manifest)
	if err != nil {
		log.Error(err)
	}
	fixedResources := fix(resources)
	err = writeManifestFile(fixedResources, rootConfig.manifest)
	if err != nil {
		return
	}
}

var autofixCmd = &cobra.Command{
	Use:   "autofix",
	Short: "Automagically fixes a manifest to be secure",
	Long: `"autofix" will examine a manifest file and automagically fill in the blanks for you leave your yaml file more secure than it found it.

Example usage:
kubeaudit autofix -f /path/to/yaml`,
	Run: autofix,
}

func init() {
	RootCmd.AddCommand(autofixCmd)
}
