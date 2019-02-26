package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkAutomountServiceAccountToken(result *Result) {
	// Check for use of deprecated service account name
	if result.DSA != "" {
		occ := Occurrence{
			id:      ErrorServiceAccountTokenDeprecated,
			kind:    Warn,
			message: "serviceAccount is a deprecated alias for ServiceAccountName, use that one instead",
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if labelExists, reason := getPodOverrideLabelReason(result, "allow-automount-service-account-token"); labelExists {
		if result.Token != nil && *result.Token {
			occ := Occurrence{
				id:       ErrorAutomountServiceAccountTokenTrueAllowed,
				kind:     Warn,
				message:  "Allowed setting automountServiceAccountToken to true",
				metadata: Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				id:       ErrorMisconfiguredKubeauditAllow,
				kind:     Warn,
				message:  "Allowed setting automountServiceAccountToken to true, but it is false or nil",
				metadata: Metadata{"Reason": prettifyReason(reason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
		return
	}

	if result.Token != nil && *result.Token && (result.SA == "" || result.SA == "default") {
		// automountServiceAccountToken = true, and serviceAccountName is blank (default: default)
		occ := Occurrence{
			id:      ErrorAutomountServiceAccountTokenTrueAndNoName,
			kind:    Error,
			message: "Default serviceAccount with token mounted. Please set automountServiceAccountToken to false",
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if result.Token == nil && result.SA == "" {
		// automountServiceAccountToken = nil (default: true), and serviceAccountName is blank (default: default)
		occ := Occurrence{
			id:      ErrorAutomountServiceAccountTokenNilAndNoName,
			kind:    Error,
			message: "Default serviceAccount with token mounted. Please set automountServiceAccountToken to false",
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
}

func auditAutomountServiceAccountToken(resource Resource) (results []Result) {
	result, err, warn := newResultFromResourceWithServiceAccountInfo(resource)
	if warn != nil {
		log.Warn(warn)
		return
	}
	if err != nil {
		log.Error(err)
		return
	}

	checkAutomountServiceAccountToken(result)
	if len(result.Occurrences) > 0 {
		results = append(results, *result)
	}
	return
}

var satCmd = &cobra.Command{
	Use:   "sat",
	Short: "Audit automountServiceAccountToken = true pods against an empty (default) service account",
	Long: `This command determines which pods are running with
autoMountServiceAcccountToken = true and default service account names.

An ERROR log is generated when a container matches one of the fol:
  automountServiceAccountToken = true and serviceAccountName is blank (default: default)
  automountServiceAccountToken = nil  and serviceAccountName is blank (default: default)

A WARN log is generated when a pod is found using Pod.Spec.DeprecatedServiceAccount
Fix this by updating serviceAccount to serviceAccountName in your .yamls

Example usage:
kubeaudit sat`,
	Run: runAudit(auditAutomountServiceAccountToken),
}

func init() {
	RootCmd.AddCommand(satCmd)
}
