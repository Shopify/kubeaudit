package cmd

import (
	"github.com/spf13/cobra"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

func checkAllowPrivilegeEscalation(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{
			id:      ErrorSecurityContextNIL,
			kind:    Error,
			message: "SecurityContext not set, please set it!",
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if reason := result.Labels["kubeaudit.allow.privilegeEscalation"]; reason == "" {
		if container.SecurityContext.AllowPrivilegeEscalation == nil {
			occ := Occurrence{
				id:      ErrorAllowPrivilegeEscalationNIL,
				kind:    Error,
				message: "AllowPrivilegeEscalation not set which allows privilege escalation, please set to false",
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else if *container.SecurityContext.AllowPrivilegeEscalation == true {
			occ := Occurrence{
				id:      ErrorAllowPrivilegeEscalationTrue,
				kind:    Error,
				message: "AllowPrivilegeEscalation set to true, please set to false",
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else if container.SecurityContext.AllowPrivilegeEscalation == nil || *container.SecurityContext.AllowPrivilegeEscalation == true {
		occ := Occurrence{
			id:       ErrorAllowPrivilegeEscalationTrueAllowed,
			kind:     Warn,
			message:  "Allowed AllowPrivilegeEscalation to be set as true",
			metadata: Metadata{"Reason": prettifyReason(reason)},
		}
		result.Occurrences = append(result.Occurrences, occ)
	} else {
		occ := Occurrence{
			id:       ErrorMisconfiguredKubeauditAllow,
			kind:     Warn,
			message:  "Allowed setting AllowPrivilegeEscalation to true, but it is set to false",
			metadata: Metadata{"Reason": prettifyReason(reason)},
		}
		result.Occurrences = append(result.Occurrences, occ)
	}
	return
}

func auditAllowPrivilegeEscalation(resource k8sRuntime.Object) (results []Result) {
	for _, container := range getContainers(resource) {
		result := newResultFromResource(resource)
		checkAllowPrivilegeEscalation(container, &result)
		if len(result.Occurrences) > 0 {
			results = append(results, result)
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
