package cmd

import (
	log "github.com/sirupsen/logrus"
)

const Version = "0.1.0"

func printKubeauditVersion() {
	log.WithFields(log.Fields{
		"Version": Version,
	}).Info("Kubeaudit")
}
