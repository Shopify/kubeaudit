package cmd

import (
	"os"

	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

// Version is the semantic versioning number for kubeaudit.
const Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of kubeaudit",
	Long:  `This just prints the version of kubeaudit`,
	Run: func(cmd *cobra.Command, args []string) {
		ver, err := version.NewVersion(Version)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		log.WithFields(log.Fields{
			"Version": ver,
		}).Info("Kubeaudit")

		kube, err := kubeClient()
		if err != nil {
			log.Warn("Could not get kubernetes server version.")
			return
		}
		printKubernetesVersion(kube)
	},
}
