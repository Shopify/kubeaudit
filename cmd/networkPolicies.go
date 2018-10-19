package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	networking "k8s.io/api/networking/v1"
)

func getNamespacesMap(namespaceList *NamespaceListV1) map[string]bool {
	// TODO filter namespaces based on annotations
	nsMap := map[string]bool{}

	for _, ns := range namespaceList.Items {
		nsMap[ns.Name] = true
	}

	return nsMap
}

func printResult(auditedNamspaces map[string]bool) {
	for ns, missingDefaultDenyNetPol := range auditedNamspaces {
		if missingDefaultDenyNetPol {
			log.WithField("KubeType", "netpol").WithField("Namespace", ns).Error("Missing default deny NeworkPolicy isolation")
		} else {
			log.WithField("KubeType", "netpol").WithField("Namespace", ns).Info("Has default deny NeworkPolicy isolation")
		}
	}
}

// Currently only checks ingress rules
func checkIfDefaultDenyPolicy(netpol networking.NetworkPolicy) bool {
	// No PodSelector is set via MatchLabels -> Catch all pods
	if len(netpol.Spec.PodSelector.MatchLabels) > 0 {
		return false
	}

	// No PodSelector is set via MatchExpressions -> catch all Pods
	if len(netpol.Spec.PodSelector.MatchExpressions) > 0 {
		return false
	}

	// No Ingress rule is defined -> denie all ingress traffic
	if len(netpol.Spec.Ingress) > 0 {
		return false
	}

	// TODO check if ingress or egress rule? if netpol.Spec.PolicyTypes
	// TODO evalute also egress?
	return true
}

func checkNamespaceNetworkPolicies(namespaceList *NamespaceListV1, netPols *NetworkPolicyListV1) map[string]bool {
	nsMap := getNamespacesMap(namespaceList)

	for _, netPol := range netPols.Items {
		// If not set check if default deny policy
		if nsMap[netPol.Namespace] {
			if checkIfDefaultDenyPolicy(netPol) {
				nsMap[netPol.Namespace] = false
			}
		}

		for _, ingress := range netPol.Spec.Ingress {
			if (len(ingress.From)) == 0 {
				log.WithField("KubeType", "netpol").
					WithField("Namespace", netPol.Namespace).
					Warn("Has allow all Networkpolicy: ", netPol.Namespace, "/", netPol.Name)
			}
		}
	}

	return nsMap
}

var npCmd = &cobra.Command{
	Use:   "np",
	Short: "Audit namespace network policies",
	Long: `This command determines whether or not a namespace has
a default deny NetworkPolicy.

An INFO log is given when a namespace has a default deny NetworkPolicy
An WARN log is given whan a namespace contains a default allow NetworkPolicy
An ERROR log is given when a namespace does not have a default deny NetworkPolicy

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
		namespaces, err := getNamespaces(kube)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		printResult(checkNamespaceNetworkPolicies(namespaces, netPols))
	},
}

func init() {
	RootCmd.AddCommand(npCmd)
}
