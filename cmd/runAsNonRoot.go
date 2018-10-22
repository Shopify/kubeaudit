package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkRunAsNonRoot(container ContainerV1, result *Result) {
	if reason := result.Labels["audit.kubernetes.io/allow-run-as-root"]; reason != "" {
		if container.SecurityContext == nil || container.SecurityContext.RunAsNonRoot == nil || *container.SecurityContext.RunAsNonRoot == false {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorRunAsNonRootFalseAllowed,
				kind:      Warn,
				message:   "Allowed setting RunAsNonRoot to false",
				metadata:  Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorMisconfiguredKubeauditAllow,
				kind:      Warn,
				message:   "Allowed setting RunAsNonRoot to false, but it is set to true",
				metadata:  Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if container.SecurityContext == nil || container.SecurityContext.RunAsNonRoot == nil {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorRunAsNonRootNil,
			kind:      Error,
			message:   "RunAsNonRoot is not set, which results in root user being allowed!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	} else if *container.SecurityContext.RunAsNonRoot == false {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorRunAsNonRootFalse,
			kind:      Error,
			message:   "RunAsNonRoot is set to false (root user allowed), please set to true!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	return
}

func auditRunAsNonRoot(resource Resource) (results []Result) {
	for _, container := range getContainers(resource) {
		result, err := newResultFromResource(resource)
		if err != nil {
			log.Error(err)
			return
		}

		checkRunAsNonRoot(container, result)
		if len(result.Occurrences) > 0 {
			results = append(results, *result)
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
kubeaudit nonroot`,
	Run: runAudit(auditRunAsNonRoot, getResources),
}

func init() {
	RootCmd.AddCommand(runAsNonRootCmd)
}
