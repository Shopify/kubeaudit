package commands

import (
	"github.com/Shopify/kubeaudit/auditors/nonroot"
	"github.com/spf13/cobra"
)

var runAsNonRootCmd = &cobra.Command{
	Use:   "nonroot",
	Short: "Audit containers running as root",
	Long: `This command determines which containers are running as root (uid=0).

An ERROR result is generated when container does not have 'runAsNonRoot = true' in either its container
  SecurityContext or its pod SecurityContext.

Example usage:
kubeaudit nonroot`,
	Run: runAudit(nonroot.New()),
}

func init() {
	RootCmd.AddCommand(runAsNonRootCmd)
}
