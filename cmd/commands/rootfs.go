package commands

import (
	"github.com/Shopify/kubeaudit/auditors/rootfs"
	"github.com/spf13/cobra"
)

var readonlyfsCmd = &cobra.Command{
	Use:   "rootfs",
	Short: "Audit containers not using a read only root filesystems",
	Long: `This command determines which containers do not have a read only root file system.

An ERROR result is generated when a container does not have 'readOnlyRootFilesystem = true' in its SecurityContext.

Example usage:
kubeaudit rootfs`,
	Run: runAudit(rootfs.New()),
}

func init() {
	RootCmd.AddCommand(readonlyfsCmd)
}
