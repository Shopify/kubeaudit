package cmd

import (
	"github.com/spf13/cobra"

	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

func checkRunAsNonRoot(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{
			id:      ErrorSecurityContextNIL,
			kind:    Error,
			message: "SecurityContext not set, please set it!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	} else if reason := result.Labels["kubeaudit.allow.runAsRoot"]; reason != "" {
		if container.SecurityContext.RunAsNonRoot == nil || *container.SecurityContext.RunAsNonRoot == false {
			occ := Occurrence{
				id:       ErrorRunAsNonRootFalseAllowed,
				kind:     Warn,
				message:  "Allowed setting RunAsNonRoot to false",
				metadata: Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				id:       ErrorMisconfiguredKubeauditAllow,
				kind:     Warn,
				message:  "Allowed setting RunAsNonRoot to false, but it is set to true",
				metadata: Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if container.SecurityContext.RunAsNonRoot == nil {
		occ := Occurrence{
			id:      ErrorRunAsNonRootNIL,
			kind:    Error,
			message: "RunAsNonRoot is not set, which results in root user being allowed!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	} else if *container.SecurityContext.RunAsNonRoot == false {
		occ := Occurrence{
			id:      ErrorRunAsNonRootFalse,
			kind:    Error,
			message: "RunAsNonRoot is set to false (root user allowed), please set to true!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	return
}

func auditRunAsNonRoot(resource k8sRuntime.Object) (results []Result) {
	for _, container := range getContainers(resource) {
		result := newResultFromResource(resource)
		checkRunAsNonRoot(container, &result)
		if len(result.Occurrences) > 0 {
			results = append(results, result)
			break
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
