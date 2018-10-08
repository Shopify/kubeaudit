package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getNamespacesMap(namespaceList *NamespaceListV1) map[string]bool {
	nsMap := map[string]bool{}

	for _, ns := range namespaceList.Items {
		nsMap[ns.Name] = false
	}

	return nsMap
}

func checkNamespaceNetworkPolicies(namespaceList *NamespaceListV1, netPols *NetworkPolicyListV1) {
	badNetPols := []NetworkPolicy{}
	nsMap := getNamespacesMap(namespaceList)

	for _, netPol := range netPols.Items {
		// namespace has a networkPolicy
		// TODO check is default deny Policy is included
		nsMap[netPol.Namespace] = true
		log.Info(netPol)
		for _, ingress := range netPol.Spec.Ingress {
			if (len(ingress.From)) == 0 {
				log.WithField("KubeType", "netpol").
					Warn("Default allow mode on ", netPol.Namespace, "/", netPol.Name)
			}
		}
	}

	for ns, hasNetPol := range nsMap {
		if hasNetPol {
			log.WithField("KubeType", "netpol").WithField("Namespace", ns).Info("Has NeworkPolicy isolation")
		} else {
			log.WithField("KubeType", "netpol").WithField("Namespace", ns).Error("Missing NeworkPolicy isolation")
		}
	}

	// TODO do we need this? badNetPols is never set
	if len(badNetPols) != 0 {
		for _, netPol := range badNetPols {
			log.WithField("KubeType", "netpol").Error(netPol.Name)
		}
	}
}

var npCmd = &cobra.Command{
	Use:   "np",
	Short: "Audit namespace network policies",
	Long: `This command determines whether or not a namespace contains
a NetworkPolicy isolation.

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
		namespaces := getNamespaces(kube)
		checkNamespaceNetworkPolicies(namespaces, netPols)
	},
}

func init() {
	RootCmd.AddCommand(npCmd)
}
