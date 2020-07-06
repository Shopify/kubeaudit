package commands

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/spf13/cobra"
)

var capabilitiesConfig capabilities.Config

func formatDropList() string {
	var buffer bytes.Buffer
	for _, cap := range capabilities.DefaultDropList {
		buffer.WriteString("\n- ")
		buffer.WriteString(cap)
	}
	return buffer.String()
}

var capabilitiesCmd = &cobra.Command{
	Use:     "capabilities",
	Aliases: []string{"caps"},
	Short:   "Audit containers not dropping capabilities",
	Long: fmt.Sprintf(`This command determines which pods have capabilities which they should not according to
the drop list. If no drop list is provided, the following capabilities are dropped:
%s

An ERROR result is generated when a pod has a capability which is on the drop list.

Example usage:
kubeaudit capabilities
kubeaudit capabilities --drop "%s"`, formatDropList(), strings.Join(capabilities.DefaultDropList[:3], ",")),
	Run: func(cmd *cobra.Command, args []string) {
		runAudit(capabilities.New(capabilitiesConfig))(cmd, args)
	},
}

func setCapabilitiesFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceVarP(&capabilitiesConfig.DropList, "drop", "d", capabilities.DefaultDropList,
		"List of capabilities that should be dropped")
}

func init() {
	RootCmd.AddCommand(capabilitiesCmd)
	setCapabilitiesFlags(capabilitiesCmd)
}
