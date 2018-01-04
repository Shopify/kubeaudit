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

func fixSecurityContextNIL(resource k8sRuntime.Object) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext = &SecurityContext{
			Capabilities: &Capabilities{
				Drop: []Capability{"AUDIT_WRITE", "CHOWN", "DAC_OVERRIDE", "FOWNER",
					"FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE", "NET_RAW", "SETFCAP",
					"SETGID", "SETUID", "SETPCAP", "SYS_CHROOT"},
			},
			Privileged: newFalse(),
			// These two aren't handled yet because we don't have detection for them yet
			//SELinuxOptions *SELinuxOptions
			//RunAsUser *int64
			RunAsNonRoot:             newTrue(),
			ReadOnlyRootFilesystem:   newTrue(),
			AllowPrivilegeEscalation: newFalse(),
		}
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixPrivilegeEscalation(resource k8sRuntime.Object) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext.Privileged = newFalse()
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixAllowPrivilegeEscalation(resource k8sRuntime.Object) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext.AllowPrivilegeEscalation = newFalse()
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixReadOnlyRootFilesystem(resource k8sRuntime.Object) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext.ReadOnlyRootFilesystem = newTrue()
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixRunAsNonRoot(resource k8sRuntime.Object) k8sRuntime.Object {
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext.RunAsNonRoot = newTrue()
		containers = append(containers, container)
	}
	return setContainers(resource, containers)
}

func fixServiceAccountToken(resource k8sRuntime.Object) k8sRuntime.Object {
	return setASAT(resource, false)
}

func fixDeprecatedServiceAccount(resource k8sRuntime.Object) k8sRuntime.Object {
	return disableDSA(resource)
}

func fixPotentialSecurityIssues(resource k8sRuntime.Object, result Result) k8sRuntime.Object {
	for _, occurrence := range result.Occurrences {
		switch occurrence.id {
		case ErrorAllowPrivilegeEscalationNIL, ErrorAllowPrivilegeEscalationTrue:
			resource = fixAllowPrivilegeEscalation(resource)
		case ErrorCapabilitiesNIL:
			resource = fixCapabilitiesNIL(resource)
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
		case ErrorSecurityContextNIL:
			resource = fixSecurityContextNIL(resource)
		case ErrorServiceAccountTokenDeprecated:
			resource = fixDeprecatedServiceAccount(resource)
		case ErrorAutomountServiceAccountTokenTrueAndNoName, ErrorAutomountServiceAccountTokenNILAndNoName:
			resource = fixServiceAccountToken(resource)
		case ErrorAllowPrivilegeEscalationTrueAllowed,
			ErrorAutomountServiceAccountTokenTrueAllowed, ErrorCapabilityAllowed,
			ErrorMisconfiguredKubeauditAllow, ErrorPrivilegedTrueAllowed,
			ErrorReadOnlyRootFilesystemFalseAllowed, ErrorRunAsNonRootFalseAllowed:
			informAboutAllowed()
		}
	}
	return resource
}

func informAboutAllowed() {
	// TODO do we want to hae this
	log.Info("something wasn't fixed because it wasn't broken")
}

func autofix(*cobra.Command, []string) {
	resources, err := getKubeResourcesManifest(rootConfig.manifest)
	if err != nil {
		log.Error(err)
	}
	var fixedResources []k8sRuntime.Object
	for _, resource := range resources {
		results := runAllAudits(resource)
		for _, result := range results {
			fixedResources = append(fixedResources, fixPotentialSecurityIssues(resource, result))
		}
	}
	err = writeManifestFile(rootConfig.manifest, fixedResources)
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
