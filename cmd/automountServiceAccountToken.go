package cmd

import (
	"github.com/spf13/cobra"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

func checkAutomountServiceAccountToken(result *Result) {
	// Check for use of deprecated service account name
	if result.DSA != "" {
		occ := Occurrence{id: ErrorServiceAccountTokenDeprecated, kind: Warn, message: "serviceAccount is a depreciated alias for ServiceAccountName, use that one instead"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if result.Token != nil && *result.Token && result.SA == "" {
		// automountServiceAccountToken = true, and serviceAccountName is blank (default: default)
		occ := Occurrence{id: ErrorServiceAccountTokenTrueAndNoName, kind: Error, message: "Default serviceAccount with token mounted. Please set AutomountServiceAccountToken to false"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if result.Token == nil && result.SA == "" {
		// automountServiceAccountToken = nil (default: true), and serviceAccountName is blank (default: default)
		occ := Occurrence{id: ErrorServiceAccountTokenNILAndNoName, kind: Error, message: "Default serviceAccount with token mounted. Please set AutomountServiceAccountToken to false"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
}

func auditAutomountServiceAccountToken(resource k8sRuntime.Object) (results []Result) {
	result := newResultFromResourceWithServiceAccountInfo(resource)
	checkAutomountServiceAccountToken(&result)
	if len(result.Occurrences) > 0 {
		results = append(results, result)
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
kubeaudit rbac sat`,
	Run: runAudit(auditAutomountServiceAccountToken),
}

func init() {
	RootCmd.AddCommand(satCmd)
}
