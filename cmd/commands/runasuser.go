package commands

import (
	"github.com/Shopify/kubeaudit/auditors/runasuser"
	"github.com/spf13/cobra"
)

var runAsUserCmd = &cobra.Command{
	Use:   "runasuser",
	Short: "Audit containers not overriding user ID with a non root user",
	Long: `This command determines which containers are not overriding the image user ID with a non root user (uid>0).

An ERROR result is generated when container does not have 'runAsUser > 0' in either its container
  SecurityContext or its pod SecurityContext.

Example usage:
kubeaudit runasuser`,
	Run: runAudit(runasuser.New()),
}

func init() {
	RootCmd.AddCommand(runAsUserCmd)
}
