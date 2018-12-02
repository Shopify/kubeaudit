package cmd

import (
	"github.com/spf13/cobra"
)

var allAuditFunctions = []interface{}{
	auditAllowPrivilegeEscalation, auditReadOnlyRootFS, auditRunAsNonRoot,
	auditAutomountServiceAccountToken, auditPrivileged, auditCapabilities,
	auditLimits, auditImages, auditAppArmor, auditSeccomp, auditNetworkPolicies,
}

var auditAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Run all audits",
	Long: `Run all audits

Example usage:
kubeaudit all -f /path/to/yaml`,
	Run: runAudit(mergeAuditFunctions(allAuditFunctions)),
}

func init() {
	RootCmd.AddCommand(auditAllCmd)
}
