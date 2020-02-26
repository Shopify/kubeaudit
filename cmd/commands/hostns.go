package commands

import (
	"github.com/Shopify/kubeaudit/auditors/hostns"
	"github.com/spf13/cobra"
)

var hostnsCmd = &cobra.Command{
	Use:     "hostns",
	Aliases: []string{"namespaces"},
	Short:   "Audit pods with hostNetwork, hostIPC or hostPID enabled",
	Long: `This command determines which pods are running with hostNetwork, hostIPC or hostPID set to 'true'.
	
An ERROR result is generated when a pod has at least one of hostNetwork, hostIPC or hostPID set to 'true'.

Example usage:
kubeaudit hostns`,
	Run: runAudit(hostns.New()),
}

func init() {
	RootCmd.AddCommand(hostnsCmd)
}
