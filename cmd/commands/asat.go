package commands

import (
	"github.com/Shopify/kubeaudit/auditors/asat"
	"github.com/spf13/cobra"
)

var asatCmd = &cobra.Command{
	Use:     "asat",
	Aliases: []string{"sat"},
	Short:   "Audit pods using an automatically mounted default service account",
	Long: `This command determines which pods are running with
autoMountServiceAcccountToken = true (or nil) and using a default service account.

An ERROR result is generated when a container matches one of the following:
  automountServiceAccountToken = true and serviceAccountName is blank (default service account)
  automountServiceAccountToken = nil (defaults to true) and serviceAccountName is blank (default service account)

A WARN result is generated when a pod is found using the deprecated 'serviceAccount' field.

Example usage:
kubeaudit asat`,
	Run: runAudit(asat.New()),
}

func init() {
	RootCmd.AddCommand(asatCmd)
}
