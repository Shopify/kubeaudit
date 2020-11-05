package commands

import (
	"github.com/Shopify/kubeaudit/auditors/limits"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var limitsConfig limits.Config

const (
	limitMemoryFlagName = "memory"
	limitCpuFlagName    = "cpu"
)

var limitsCmd = &cobra.Command{
	Use:   "limits",
	Short: "Audit containers exceeding a specified CPU or memory limit",
	Long: `This command determines which containers exceed the specified CPU and memory limits, or have no limits configured.

A WARN result is generated for each of the following cases:
  - The CPU limit is unset or exceeds the specified CPU limit
  - The memory limit is unset or exceeds the specified memory limit

Example usage:
kubeaudit limits
kubeaudit limits --cpu 500m --memory 256Mi`,
	Run: func(cmd *cobra.Command, args []string) {
		auditor, err := limits.New(limitsConfig)
		if err != nil {
			log.Fatal("failed to create limits auditor")
		}
		runAudit(auditor)(cmd, args)
	},
}

func setLimitsFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&limitsConfig.CPU, limitCpuFlagName, "", "Max CPU limit")
	cmd.Flags().StringVar(&limitsConfig.Memory, limitMemoryFlagName, "", "Max memory limit")
}

func init() {
	RootCmd.AddCommand(limitsCmd)
	setLimitsFlags(limitsCmd)
}
