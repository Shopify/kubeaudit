package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkAllowPrivilegeEscalation(container ContainerV1, result *Result) {
	if reason := result.Labels["audit.kubernetes.io/allow-privilege-escalation"]; reason == "" {
		if container.SecurityContext == nil || container.SecurityContext.AllowPrivilegeEscalation == nil {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorAllowPrivilegeEscalationNil,
				kind:      Error,
				message:   "AllowPrivilegeEscalation not set which allows privilege escalation, please set to false",
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else if *container.SecurityContext.AllowPrivilegeEscalation == true {
			occ := Occurrence{
				container: container.Name,
				id:        ErrorAllowPrivilegeEscalationTrue,
				kind:      Error,
				message:   "AllowPrivilegeEscalation set to true, please set to false",
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if container.SecurityContext == nil || container.SecurityContext.AllowPrivilegeEscalation == nil || *container.SecurityContext.AllowPrivilegeEscalation == true {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorAllowPrivilegeEscalationTrueAllowed,
			kind:      Warn,
			message:   "Allowed AllowPrivilegeEscalation to be set as true",
			metadata:  Metadata{"Reason": prettifyReason(reason)},
		}
		result.Occurrences = append(result.Occurrences, occ)
	} else {
		occ := Occurrence{
			container: container.Name,
			id:        ErrorMisconfiguredKubeauditAllow,
			kind:      Warn,
			message:   "Allowed setting AllowPrivilegeEscalation to true, but it is set to false",
			metadata:  Metadata{"Reason": prettifyReason(reason)},
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	return
}

func auditAllowPrivilegeEscalation(resource Resource) (results []Result) {
	for _, container := range getContainers(resource) {
		result, err := newResultFromResource(resource)
		if err != nil {
			log.Error(err)
			return
		}

		checkAllowPrivilegeEscalation(container, result)
		if len(result.Occurrences) > 0 {
			results = append(results, *result)
			break
		}
	}
	return
}

var allowPrivilegeEscalationCmd = &cobra.Command{
	Use:   "allowpe",
	Short: "Audit containers that allow privilege escalation",
	Long: `This command determines which containers in a kubernetes cluster allow privilege escalation.

A PASS is given when a container does not allow privilege escalation
A FAIL is generated when a container allows privilege escalation

Example usage:
kubeaudit allowpe`,
	Run: runAudit(auditAllowPrivilegeEscalation),
}

func init() {
	RootCmd.AddCommand(allowPrivilegeEscalationCmd)
}
