package cmd

import (
	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
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

		kubeconfig := rootConfig.kubeConfig
		if rootConfig.localMode {
			kubeconfig = os.Getenv("HOME") + "/.kube/config"
		}

		kube, err := kubeClient(kubeconfig)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		printKubernetesVersion(kube)
	},
}
