package cmd

import (
	"github.com/spf13/cobra"
)

func checkRunAsNonRoot(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{id: ErrorSecurityContextNIL, kind: Error, message: "SecurityContext not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if container.SecurityContext.RunAsNonRoot == nil {
		occ := Occurrence{id: ErrorRunAsNonRootNIL, kind: Error, message: "RunAsNonRoot is not set, which results in root user being allowed!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if *container.SecurityContext.RunAsNonRoot == false {
		occ := Occurrence{id: ErrorRunAsNonRootFalse, kind: Error, message: "RunAsNonRoot is set to false (root user allowed), please set to true!"}
		result.Occurrences = append(result.Occurrences, occ)
	}
}

func auditRunAsNonRoot(items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkRunAsNonRoot(container, result)
			if result != nil && len(result.Occurrences) > 0 {
				results = append(results, *result)
				break
			}
		}
	}
	return
}

// runAsNonRootCmd represents the runAsNonRoot command
var runAsNonRootCmd = &cobra.Command{
	Use:   "nonroot",
	Short: "Audit containers running as root",
	Long: `This command determines which containers in a kubernetes cluster
are running as root (uid=0).

A PASS is given when a container runs as a uid greater than 0
A FAIL is generated when a container runs as root

Example usage:
kubeaudit runAsNonRoot`,
	Run: runAudit(auditRunAsNonRoot),
}

func init() {
	RootCmd.AddCommand(runAsNonRootCmd)
}
