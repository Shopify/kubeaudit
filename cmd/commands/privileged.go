package commands

import (
	"github.com/Shopify/kubeaudit/auditors/privileged"
	"github.com/spf13/cobra"
)

var privilegedCmd = &cobra.Command{
	Use:     "privileged",
	Aliases: []string{"priv"},
	Short:   "Audit containers running as privileged",
	Long: `This command determines which containers are running as privileged.

An ERROR result is generated when a container has 'privileged = true' in its SecurityContext.

A WARN result is generated a when a container has 'privileged = nil' in its SecurityContext. 'privileged'
  defaults to 'true' so this is ok, but it should be explicitly set to 'true'.

Example usage:
kubeaudit priv`,
	Run: runAudit(privileged.New()),
}

func init() {
	RootCmd.AddCommand(privilegedCmd)
}
