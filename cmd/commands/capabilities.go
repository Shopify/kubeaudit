package commands

import (
	"fmt"

	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/spf13/cobra"
)

var capabilitiesConfig capabilities.Config

const capsAddFlagName = "allow-add-list"

var capabilitiesCmd = &cobra.Command{
	Use:     "capabilities",
	Aliases: []string{"caps"},
	Short:   "Audit containers not dropping ALL capabilities",
	Long: fmt.Sprintf(`This command determines which pods either have capabilities added or not set to ALL:
An ERROR result is generated when a pod does not have drop ALL specified or when a capability is added. In case 
you need specific capabilities you can add them with the '--allow-add-list' flag, so kubeaudit will not report errors.

Example usage:
kubeaudit capabilities
kubeaudit capabilities --allow-add-list "%s"`, "CHOWN"),
	Run: func(cmd *cobra.Command, args []string) {
		runAudit(capabilities.New(capabilitiesConfig))(cmd, args)
	},
}

func setCapabilitiesFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceVar(&capabilitiesConfig.AllowAddList, capsAddFlagName, capabilities.DefaultAllowAddList,
		"Comma separated list of added capabilities that can be ignored by kubeaudit reports")
}

func init() {
	RootCmd.AddCommand(capabilitiesCmd)
	setCapabilitiesFlags(capabilitiesCmd)
}
