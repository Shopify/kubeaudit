package commands

import (
	"github.com/Shopify/kubeaudit/auditors/netpols"
	"github.com/spf13/cobra"
)

var netpolsCmd = &cobra.Command{
	Use:     "netpols",
	Aliases: []string{"np"},
	Short:   "Audit namespaces that do not have a default deny network policy",
	Long: `This command determines which namespaces do not have a default deny NetworkPolicy.

An ERROR result is generated for each of the followign cases:
  - A namespace does not have a default deny-all-ingress NetworkPolicy
  - A namespace does not have a default deny-all-egress NetworkPolicy

A WARN result is generated for each of the following cases:
  - A namespace has a default allow-all-ingress NetworkPolicy
  - A namespace has a default allow-all-egress NetworkPolicy


Example usage:
kubeaudit netpols`,
	Run: runAudit(netpols.New()),
}

func init() {
	RootCmd.AddCommand(netpolsCmd)
}
