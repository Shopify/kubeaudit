package cmd

import (
	"github.com/spf13/cobra"

	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

func checkReadOnlyRootFS(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{id: ErrorSecurityContextNIL, kind: Error, message: "SecurityContext not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if container.SecurityContext.ReadOnlyRootFilesystem == nil {
		occ := Occurrence{id: ErrorReadOnlyRootFilesystemNIL, kind: Error, message: "ReadOnlyRootFilesystem not set which results in a writable rootFS, please set to true"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if !*container.SecurityContext.ReadOnlyRootFilesystem {
		occ := Occurrence{id: ErrorReadOnlyRootFilesystemFalse, kind: Error, message: "ReadOnlyRootFilesystem set to false, please set to true"}
		result.Occurrences = append(result.Occurrences, occ)
	}
}

func auditReadOnlyRootFS(resource k8sRuntime.Object) (results []Result) {
	for _, container := range getContainers(resource) {
		result := newResultFromResource(resource)
		checkReadOnlyRootFS(container, &result)
		if len(result.Occurrences) > 0 {
			results = append(results, result)
			break
		}
	}
	return
}

var readonlyfsCmd = &cobra.Command{
	Use:   "rootfs",
	Short: "Audit containers with read only root filesystems",
	Long: `This command determines which containers in a kubernetes cluster
have their filesystems set to read only.

A PASS is given when a container has a read only root filesystem
A FAIL is given when a container does not have a read only root filesystem

Example usage:
kubeaudit runAsNonRoot`,
	Run: runAudit(auditReadOnlyRootFS),
}

func init() {
	RootCmd.AddCommand(readonlyfsCmd)
}
