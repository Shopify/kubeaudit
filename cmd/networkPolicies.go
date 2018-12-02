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

	// No Ingress rule is defined -> deny all ingress traffic
	if len(netpol.Spec.Ingress) > 0 {
		return false
	}

	//  No Egress rule is defined -> deny all egress traffic
	if len(netpol.Spec.Egress) > 0 {
		return false
	}

	return true
}

func checkNamespaceNetworkPolicies(netPols *NetworkPolicyListV1, result *Result) {
	// TODO check if any netpol is default Policy
	hasDefaultDeny := false

	for _, netPol := range netPols.Items {
		// If not set check if default deny policy
		if !hasDefaultDeny {
			hasDefaultDeny = checkIfDefaultDenyPolicy(netPol)
		}

		for _, ingress := range netPol.Spec.Ingress {
			// Allow all ingress traffic
			if (len(ingress.From)) == 0 {
				occ := Occurrence{
					container: "",
					id:        WarningAllowAllIngressNetworkPolicyExists,
					kind:      Warn,
					message:   "Found allow all ingress traffic NetworkPolicy",
					metadata: Metadata{
						"PolicyName": netPol.ObjectMeta.Name,
					},
				}
				result.Occurrences = append(result.Occurrences, occ)
			}
		}

		for _, egress := range netPol.Spec.Egress {
			// Allow all egress traffic
			if (len(egress.To)) == 0 {
				occ := Occurrence{
					container: "",
					id:        WarningAllowAllEgressNetworkPolicyExists,
					kind:      Warn,
					message:   "Found allow all egress traffic NetworkPolicy",
					metadata: Metadata{
						"PolicyName": netPol.ObjectMeta.Name,
					},
				}
				result.Occurrences = append(result.Occurrences, occ)
			}
		}
	}

	if !hasDefaultDeny {
		occ := Occurrence{
			container: "",
			id:        ErrorMissingDefaultDenyNetworkPolicy,
			kind:      Error,
			message:   "Namespace is missing a default deny NetworkPolicy",
		}
		result.Occurrences = append(result.Occurrences, occ)
	} else {
		occ := Occurrence{
			container: "",
			id:        InfoDefaultDenyNetworkPolicyExists,
			kind:      Info,
			message:   "Namespace has a default deny NetworkPolicy",
		}
		result.Occurrences = append(result.Occurrences, occ)
	}

	return
}

func getNetworkPoliciesResources(namespace string) (netPolList *NetworkPolicyList, err error) {
	// Prevent the return of a nil value
	netPolList = &NetworkPolicyList{}
	if rootConfig.manifest != "" {
		resources, err := getKubeResourcesManifest(rootConfig.manifest)
		if err != nil {
			return netPolList, nil
		}

		for _, resource := range resources {
			switch kubeType := resource.(type) {
			case *NetworkPolicy:
				netPolList.Items = append(netPolList.Items, *kubeType)
			}
		}

		return netPolList, nil
	}

	currentRootNamespace := rootConfig.namespace
	kube, err := kubeClient()
	if err != nil {
		return netPolList, err
	}

	if namespace != "" {
		rootConfig.namespace = namespace
	}

	netPolList = getNetworkPolicies(kube)

	rootConfig.namespace = currentRootNamespace
	return netPolList, err
}

func getNamespaceName(resource k8sRuntime.Object) (ns string) {
	switch kubeType := resource.(type) {
	case *Namespace:
		ns = kubeType.ObjectMeta.Name
	}

	return ns
}

func auditNetworkPolicies(resource k8sRuntime.Object) (results []Result) {
	nsName := getNamespaceName(resource)
	// We found no namespace
	if nsName == "" {
		return
	}

	// iterate over namespaces not netpol --> actually an namespace check not an netpol check
	result, err := newResultFromResource(resource)
	if err != nil {
		log.Error(err)
		return
	}

	// Fetch NetworkPolicies for the current namespace
	netPols, err := getNetworkPoliciesResources(nsName)
	if err != nil {
		log.Error(err)
		return
	}

	checkNamespaceNetworkPolicies(netPols, result)
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
	Run: runAudit(auditNetworkPolicies),
}

func init() {
	RootCmd.AddCommand(npCmd)
}
