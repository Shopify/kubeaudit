package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	networking "k8s.io/api/networking/v1"
)

// isNetworkPolicyType checks if the NetworkPolicy is from the provided type e.g. egress
func isNetworkPolicyType(netPol networking.NetworkPolicy, netPolType string) bool {
	for _, polType := range netPol.Spec.PolicyTypes {
		if string(polType) == netPolType {
			return true
		}
	}
	return false
}

// checkIfDefaultDenyPolicy checks if the policy contains a deny all ingress / egress rules
func checkIfDefaultDenyPolicy(netPol networking.NetworkPolicy) (bool, bool) {
	hasDenyAllIngressRule, hasDenyAllEgressRule := false, false

	// No PodSelector is set via MatchLabels -> Catch all pods
	if len(netPol.Spec.PodSelector.MatchLabels) > 0 {
		return hasDenyAllIngressRule, hasDenyAllEgressRule
	}

	// No PodSelector is set via MatchExpressions -> Catch all Pods
	if len(netPol.Spec.PodSelector.MatchExpressions) > 0 {
		return hasDenyAllIngressRule, hasDenyAllEgressRule
	}

	// No ingress is defined -> Deny all ingress traffic
	if len(netPol.Spec.Ingress) == 0 && isNetworkPolicyType(netPol, "Ingress") {
		hasDenyAllIngressRule = true
	}

	// No Egress rule is defined -> Deny all egress traffic
	if len(netPol.Spec.Egress) == 0 && isNetworkPolicyType(netPol, "Egress") {
		hasDenyAllEgressRule = true
	}

	return hasDenyAllIngressRule, hasDenyAllEgressRule
}

func checkNamespaceNetworkPolicies(netPols *NetworkPolicyListV1, result *Result, nsName string) {
	hasDenyAllIngressRule, hasDenyAllEgressRule := false, false

	for _, netPol := range netPols.Items {
		// If not set check if default deny policy
		if !hasDenyAllIngressRule || !hasDenyAllEgressRule {
			resDenyAllIngress, resDenyAllEgress := checkIfDefaultDenyPolicy(netPol)
			// If hasDenyAllIngressRule is not already set use the result from above
			// We need this extra step because the policies could be splitted over
			// two network policies
			if !hasDenyAllIngressRule {
				hasDenyAllIngressRule = resDenyAllIngress
			}

			// Same as for hasDenyAllIngressRule
			if !hasDenyAllEgressRule {
				hasDenyAllEgressRule = resDenyAllEgress
			}
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

	egressLabelExists, egressReason := getNamespaceOverrideLabelReason(result, nsName, "egress")
	ingressLabelExists, ingressReason := getNamespaceOverrideLabelReason(result, nsName, "ingress")

	if egressLabelExists && ingressLabelExists {
		if !hasDenyAllEgressRule && !hasDenyAllIngressRule {
			occ := Occurrence{
				container: "",
				id:        ErrorMissingDefaultDenyIngressAndEgressNetworkPolicyAllowed,
				kind:      Warn,
				message:   "Allowed Namespace is missing a default deny ingress and default deny egress NetworkPolicy",
				metadata:  Metadata{"Reason": prettifyReason("Ingress: " + ingressReason + " Egress: " + egressReason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else {
		if !hasDenyAllEgressRule && !hasDenyAllIngressRule {
			occ := Occurrence{
				container: "",
				id:        ErrorMissingDefaultDenyIngressAndEgressNetworkPolicy,
				kind:      Error,
				message:   "Namespace is missing a default deny ingress and default deny egress NetworkPolicy",
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	}

	if ingressLabelExists && !egressLabelExists {
		if !hasDenyAllIngressRule && hasDenyAllEgressRule {
			occ := Occurrence{
				container: "",
				id:        ErrorMissingDefaultDenyIngressNetworkPolicyAllowed,
				kind:      Warn,
				message:   "Allowed Namespace is missing a default deny ingress NetworkPolicy",
				metadata:  Metadata{"Reason": prettifyReason(ingressReason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else {
		if !hasDenyAllIngressRule && hasDenyAllEgressRule {
			occ := Occurrence{
				container: "",
				id:        ErrorMissingDefaultDenyIngressNetworkPolicy,
				kind:      Error,
				message:   "Namespace is missing a default deny ingress NetworkPolicy",
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	}

	if !ingressLabelExists && egressLabelExists {
		if hasDenyAllIngressRule && !hasDenyAllEgressRule {
			occ := Occurrence{
				container: "",
				id:        ErrorMissingDefaultDenyEgressNetworkPolicyAllowed,
				kind:      Warn,
				message:   "Allowed Namespace is missing a default deny ingress NetworkPolicy",
				metadata:  Metadata{"Reason": prettifyReason(egressReason)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	} else {
		if hasDenyAllIngressRule && !hasDenyAllEgressRule {
			occ := Occurrence{
				container: "",
				id:        ErrorMissingDefaultDenyEgressNetworkPolicy,
				kind:      Error,
				message:   "Namespace is missing a default deny ingress NetworkPolicy",
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	}

	if hasDenyAllIngressRule && hasDenyAllEgressRule {
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

func getNetworkPoliciesResources(namespace string) (netPolList *NetworkPolicyListV1, err error) {
	// Prevent the return of a nil value
	netPolList = &NetworkPolicyListV1{}
	if rootConfig.manifest != "" {
		resources, err := getKubeResourcesManifest(rootConfig.manifest)
		if err != nil {
			return netPolList, err
		}

		for _, resource := range resources {
			switch kubeType := resource.(type) {
			case *NetworkPolicyV1:
				if kubeType.ObjectMeta.Namespace == namespace {
					netPolList.Items = append(netPolList.Items, *kubeType)
				}
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
	return netPolList, nil
}

func getNamespaceName(resource Resource) string {
	name := ""
	ns, ok := resource.(*NamespaceV1)
	if ok {
		name = ns.ObjectMeta.Name
	}
	return name
}

func auditNetworkPolicies(resource Resource) (results []Result) {
	nsName := getNamespaceName(resource)

	// We found no namespace
	if nsName == "" {
		return
	}

	// iterate over namespaces not netpol --> actually an namespace check not an netpol check
	result, err, warn := newResultFromResource(resource)
	if warn != nil {
		log.Warn(warn)
		return
	}
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

	checkNamespaceNetworkPolicies(netPols, result, nsName)
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
