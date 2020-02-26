package commands

import (
	"github.com/Shopify/kubeaudit/auditors/mountds"
	"github.com/spf13/cobra"
)

var mountdsCmd = &cobra.Command{
	Use:   "mountds",
	Short: "Audit containers that mount /var/run/docker.sock",
	Long: `This command determines which containers mount /var/run/docker.sock. 

A WARN result is generated when a container mounts '/var/run/docker.sock'.

Example usage:
kubeaudit mountds`,
	Run: runAudit(mountds.New()),
}

func init() {
	RootCmd.AddCommand(mountdsCmd)
}
