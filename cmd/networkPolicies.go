package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func checkNamespaceNetworkPolicies(netPols *NetworkPolicyList) {
	badNetPols := []NetworkPolicy{}

	for _, netPol := range netPols.Items {
		for _, ingress := range netPol.Spec.Ingress {
			if (len(ingress.From)) == 0 {
				log.WithField("type", "netpol").
					Warn("Default allow mode on ", netPol.Namespace, "/", netPol.Name)
			}
		}
	}

	if len(badNetPols) != 0 {
		for _, netPol := range badNetPols {
			log.WithField("type", "netpol").Error(netPol.Name)
		}
	}
}

var npCmd = &cobra.Command{
	Use:   "np",
	Short: "Audit namespace network policies",
	Long: `This command determines whether or not a namespace contains
a NetworkPolicy isolation annotation.

An INFO log is given when a namespace has NetworkPolicy isolation
An ERROR log is given when a namespace does not have NetworkPolicy isolation

Example usage:
kubeaudit np`,
	Run: func(cmd *cobra.Command, args []string) {
		kube, err := kubeClient()
		if err != nil {
			log.Error(err)
		}

		if rootConfig.json {
			log.SetFormatter(&log.JSONFormatter{})
		}
		netPols := getNetworkPolicies(kube)
		checkNamespaceNetworkPolicies(netPols)
	},
}

func init() {
	RootCmd.AddCommand(npCmd)
}
