package commands

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/spf13/cobra"
)

var capabilitiesConfig customCapabilitiesConfig

type customCapabilitiesConfig struct {
	dropList string
}

func (conf customCapabilitiesConfig) ToConfig() capabilities.Config {
	return capabilities.Config{
		DropList: strings.Split(capabilitiesConfig.dropList, " "),
	}
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
kubeaudit capabilities --drop "%s"`, formatDropList(), strings.Join(capabilities.DefaultDropList[:3], " ")),
	Run: runAudit(capabilities.New(capabilitiesConfig.ToConfig())),
}

func formatDropList() string {
	var buffer bytes.Buffer
	for _, cap := range capabilities.DefaultDropList {
		buffer.WriteString("\n- ")
		buffer.WriteString(cap)
	}
	return buffer.String()
}

func init() {
	RootCmd.AddCommand(capabilitiesCmd)
	setCapabilitiesFlags(capabilitiesCmd)
}

func setCapabilitiesFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&capabilitiesConfig.dropList, "drop", "d", strings.Join(capabilities.DefaultDropList, " "),
		"List of capabilities that should be dropped")
}
