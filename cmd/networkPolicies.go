package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	networking "k8s.io/api/networking/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

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

func checkNamespaceNetworkPolicies(netPols *NetworkPolicyListV1, result *Result) {
	// TODO check if any netpol is default Policy
	hasDefaultDeny := false
	if len(netPols.Items) == 0 {
		// TODO error no defaultPolicy

		return
	}

	for _, netPol := range netPols.Items {
		// If not set check if default deny policy
		if !hasDefaultDeny {
			hasDefaultDeny = checkIfDefaultDenyPolicy(netPol)
		}

		for _, ingress := range netPol.Spec.Ingress {
			if (len(ingress.From)) == 0 {
				/*
					occ := Occurrence{
						container: container.Name,
						id:        ErrorRunAsNonRootFalse,
						kind:      Error,
						message:   "RunAsNonRoot is set to false (root user allowed), please set to true!",
					}
					result.Occurrences = append(result.Occurrences, occ)
				*/

				log.WithField("KubeType", "netpol").
					WithField("Namespace", netPol.Namespace).
					Warn("Has allow all Networkpolicy: ", netPol.Namespace, "/", netPol.Name)
			}
		}
	}

	if !hasDefaultDeny {
		//TODO
	}

	return
}

func getNamespaceResources() (resources []k8sRuntime.Object, err error) {
	kube, err := kubeClient()
	if err != nil {
		return
	}

	nsList, err := getNamespaces(kube)
	if err != nil {
		return
	}
	for _, resource := range nsList.Items {
		resources = append(resources, resource.DeepCopyObject())
	}
	return
}

//TODO can we set/get only specific resources?
func auditNetworkPolicies(resource k8sRuntime.Object) (results []Result) {
	/*kube, err := kubeClient()
	if err != nil {
		log.Error(err)
		return
	}

	/*if rootConfig.json {
		log.SetFormatter(&log.JSONFormatter{})
	}
	netPols := getNetworkPolicies(kube)
	namespaces, err := getNamespaces(kube)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}*/

	//TODO fetch network policies
	// iterate over namespaces not net pol --> actually an namespace check not an netpol check
	log.Info(resource)

	result, err := newResultFromResource(resource)
	if err != nil {
		log.Error(err)
		return
	}

	checkNamespaceNetworkPolicies(nil, result)
	if len(result.Occurrences) > 0 {
		results = append(results, *result)
	}

	return
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
	Run: runAudit(auditNetworkPolicies, getNamespaceResources),
}

func init() {
	RootCmd.AddCommand(npCmd)
}
