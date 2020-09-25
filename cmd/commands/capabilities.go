package commands

import (
	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/spf13/cobra"
)

var capabilitiesConfig capabilities.Config

var capabilitiesCmd = &cobra.Command{
	Use:     "capabilities",
	Aliases: []string{"caps"},
	Short:   "Audit containers not dropping ALL capabilities",
	Long: `This command determines which pods either have capabilities added or not set to ALL:
An ERROR result is generated when a pod does not have drop ALL specified or when a capability is added.

Example usage:
kubeaudit capabilities`,
	Run: runAudit(capabilities.New(capabilitiesConfig)),
}

func init() {
	RootCmd.AddCommand(capabilitiesCmd)
}
