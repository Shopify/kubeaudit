package netpols

import "github.com/Shopify/kubeaudit/pkg/k8s"

const AllNamespaces = ""

func getNetworkPolicies(resources []k8s.Resource, namespace string) (networkPolicies []*k8s.NetworkPolicyV1) {
	for _, resource := range resources {
		networkPolicy, ok := resource.(*k8s.NetworkPolicyV1)
		if ok && (namespace == AllNamespaces || getResourceNamespace(resource) == namespace) {
			networkPolicies = append(networkPolicies, networkPolicy)
		}
	}
	return
}

// isNetworkPolicyType checks if the NetworkPolicy applies to the specified policy type (Ingress or Egress)
func isNetworkPolicyType(netPol *k8s.NetworkPolicyV1, netPolType string) bool {
	for _, polType := range netPol.Spec.PolicyTypes {
		if string(polType) == netPolType {
			return true
		}
	}
	return false
}

func getResourceNamespace(resource k8s.Resource) string {
	switch kubeType := resource.(type) {
	case *k8s.NamespaceV1:
		return kubeType.ObjectMeta.Name
	case *k8s.NetworkPolicyV1:
		return kubeType.ObjectMeta.Namespace
	}
	return ""
}

func allIngressTrafficAllowed(networkPolicy *k8s.NetworkPolicyV1) bool {
	for _, ingress := range networkPolicy.Spec.Ingress {
		if (len(ingress.From)) == 0 {
			return true
		}
	}
	return false
}

func allEgressTrafficAllowed(networkPolicy *k8s.NetworkPolicyV1) bool {
	for _, egress := range networkPolicy.Spec.Egress {
		if (len(egress.To)) == 0 {
			return true
		}
	}
	return false
}

func hasCatchAllNetworkPolicy(networkPolicies []*k8s.NetworkPolicyV1) (bool, *k8s.NetworkPolicyV1) {
	for _, networkPolicy := range networkPolicies {
		if isCatchAllNetworkPolicy(networkPolicy) {
			return true, networkPolicy
		}
	}

	return false, nil
}

func isCatchAllNetworkPolicy(networkPolicy *k8s.NetworkPolicyV1) bool {
	// No PodSelector is set via MatchLabels -> Catch all pods
	if len(networkPolicy.Spec.PodSelector.MatchLabels) > 0 {
		return false
	}

	// No PodSelector is set via MatchExpressions -> Catch all Pods
	if len(networkPolicy.Spec.PodSelector.MatchExpressions) > 0 {
		return false
	}

	return true
}

func hasDenyAllIngress(networkPolicies []*k8s.NetworkPolicyV1) bool {
	for _, networkPolicy := range networkPolicies {
		if networkPolicy == nil {
			continue
		}

		if len(networkPolicy.Spec.Ingress) != 0 {
			continue
		}

		if !isNetworkPolicyType(networkPolicy, Ingress) {
			continue
		}

		if isCatchAllNetworkPolicy(networkPolicy) {
			return true
		}
	}

	return false
}

func hasDenyAllEgress(networkPolicies []*k8s.NetworkPolicyV1) bool {
	for _, networkPolicy := range networkPolicies {
		if networkPolicy == nil {
			continue
		}

		if len(networkPolicy.Spec.Egress) != 0 {
			continue
		}

		if !isNetworkPolicyType(networkPolicy, Egress) {
			continue
		}

		if isCatchAllNetworkPolicy(networkPolicy) {
			return true
		}
	}

	return false
}
