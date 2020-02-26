package commands

import (
	"github.com/Shopify/kubeaudit/auditors/privesc"
	"github.com/spf13/cobra"
)

var allowPrivilegeEscalationCmd = &cobra.Command{
	Use:     "privesc",
	Aliases: []string{"allowpe"},
	Short:   "Audit containers that allow privilege escalation",
	Long: `This command determines which containers allow privilege escalation.

An ERROR result is generated when a container does not have 'allowPrivilegeEscalation = false' in its
  SecurityContext.

Example usage:
kubeaudit privesc`,
	Run: runAudit(privesc.New()),
}

func init() {
	RootCmd.AddCommand(allowPrivilegeEscalationCmd)
}
