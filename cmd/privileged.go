package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkPrivileged(container ContainerV1, result *Result) {
	if container.SecurityContext == nil || container.SecurityContext.Privileged == nil {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorPrivilegedNil,
			kind:      Warn,
			message:   "Privileged defaults to false, which results in non privileged, which is okay.",
		}
		result.Occurrences = append(result.Occurrences, occ)
	} else if reason := result.Labels["audit.kubernetes.io/allow-privileged"]; reason != "" {
		if *container.SecurityContext.Privileged == true {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorPrivilegedTrueAllowed,
				kind:      Warn,
				message:   "Allowed setting privileged to true",
				metadata:  Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorMisconfiguredKubeauditAllow,
				kind:      Warn,
				message:   "Allowed setting privileged to true, but privileged is false or nil",
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if *container.SecurityContext.Privileged == true {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorPrivilegedTrue,
			kind:      Error,
			message:   "Privileged set to true! Please change it to false!",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	return
}

func auditPrivileged(resource Resource) (results []Result) {
	for _, container := range getContainers(resource) {
		result, err := newResultFromResource(resource)
		if err != nil {
			log.Error(err)
			return
		}

		checkPrivileged(container, result)
		if len(result.Occurrences) > 0 {
			results = append(results, *result)
		}
	}
	return
}

var privileged = &cobra.Command{
	Use:   "priv",
	Short: "Audit containers running as privileged",
	Long: `This command determines which containers in a kubernetes cluster
are running as privileged.

A PASS is given when a container runs in a non-privileged mode
A FAIL is generated when a container runs in a privileged mode

Example usage:
kubeaudit priv`,
	Run: runAudit(auditPrivileged),
}

func init() {
	RootCmd.AddCommand(privileged)
}
