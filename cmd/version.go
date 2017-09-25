package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

const Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of kubeaudit",
	Long:  `This just prints the version of kubeaudit`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"Version": Version,
		}).Info("Kubeaudit")
	},
}
