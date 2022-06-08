package commands

import (
	"github.com/Shopify/kubeaudit/auditors/deprecatedapis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deprecatedapisConfig deprecatedapis.Config

const (
	currentVersionFlagName  = "current-k8s-version"
	targetedVersionFlagName = "targeted-k8s-version"
)

var deprecatedapisCmd = &cobra.Command{
	Use:   "deprecatedapis",
	Short: "Audit resource API version deprecations",
	Long: `This command determines which resource is defined with a deprecated API version.

An ERROR result is generated for API version not available in the targeted version
A WARN result is generated for API version deprecated in the current version
An INFO result is generated for API version not yet deprecated in the current version

Example usage:
kubeaudit deprecatedapis
kubeaudit deprecatedapis --current-k8s-version 1.22 --targeted-k8s-version 1.24`,
	Run: func(cmd *cobra.Command, args []string) {
		auditor, err := deprecatedapis.New(deprecatedapisConfig)
		if err != nil {
			log.Fatal("failed to create deprecatedapis auditor")
		}
		runAudit(auditor)(cmd, args)
	},
}

func setdeprecatedapisFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&deprecatedapisConfig.CurrentVersion, currentVersionFlagName, "", "Kubernetes current version (eg 1.22)")
	cmd.Flags().StringVar(&deprecatedapisConfig.TargetedVersion, targetedVersionFlagName, "", "Kubernetes version to migrate to (eg 1.24)")
}

func init() {
	RootCmd.AddCommand(deprecatedapisCmd)
	setdeprecatedapisFlags(deprecatedapisCmd)
}
