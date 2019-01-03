package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkReadOnlyRootFS(container ContainerV1, result *Result) {
	if reason := result.Labels["audit.kubernetes.io/allow-read-only-root-filesystem-false"]; reason != "" {
		if container.SecurityContext == nil || container.SecurityContext.ReadOnlyRootFilesystem == nil || *container.SecurityContext.ReadOnlyRootFilesystem == false {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorReadOnlyRootFilesystemFalseAllowed,
				kind:      Warn,
				message:   "Allowed setting readOnlyRootFilesystem to false",
				metadata:  Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorMisconfiguredKubeauditAllow,
				kind:      Warn,
				message:   "Allowed setting readOnlyRootFilesystem to false, but it is set to true",
				metadata:  Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if container.SecurityContext == nil || container.SecurityContext.ReadOnlyRootFilesystem == nil {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorReadOnlyRootFilesystemNil,
			kind:      Error,
			message:   "ReadOnlyRootFilesystem not set which results in a writable rootFS, please set to true",
		}
		result.Occurrences = append(result.Occurrences, occ)
	} else if !*container.SecurityContext.ReadOnlyRootFilesystem {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorReadOnlyRootFilesystemFalse,
			kind:      Error,
			message:   "ReadOnlyRootFilesystem set to false, please set to true",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	return
}

func auditReadOnlyRootFS(resource Resource) (results []Result) {
	for _, container := range getContainers(resource) {
		result, err := newResultFromResource(resource)
		if err != nil {
			log.Error(err)
			return
		}

		checkReadOnlyRootFS(container, result)
		if len(result.Occurrences) > 0 {
			results = append(results, *result)
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
kubeaudit rootfs`,
	Run: runAudit(auditReadOnlyRootFS),
}

func init() {
	RootCmd.AddCommand(readonlyfsCmd)
}
