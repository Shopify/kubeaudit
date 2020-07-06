package commands

import (
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/spf13/cobra"
)

var appArmorCmd = &cobra.Command{
	Use:   "apparmor",
	Short: "Audit containers running without AppArmor",
	Long: `This command determines which containers are running without AppArmor enabled.

An ERROR result is generated when a container has AppArmor disabled or misconfigured.

Example usage:
kubeaudit apparmor`,
	Run: runAudit(apparmor.New()),
}

func init() {
	RootCmd.AddCommand(appArmorCmd)
}
