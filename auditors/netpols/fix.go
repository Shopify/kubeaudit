package netpols

import (
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit/pkg/k8s"
)

const DefaultDenyNetworkPolicyName = "default-deny"

type fixByAddingNetworkPolicy struct {
	policyList []string
	namespace  string
}

func (f *fixByAddingNetworkPolicy) Plan() string {
	return fmt.Sprintf("Create a new NetworkPolicy resource which denies all %s traffic", strings.Join(f.policyList, " and "))
}

func (f *fixByAddingNetworkPolicy) Apply(resource k8s.Resource) []k8s.Resource {
	return []k8s.Resource{newDefaultDenyNetworkPolicy(f.namespace, f.policyList)}
}

type fixByAddingPolicyToNetPol struct {
	networkPolicy *k8s.NetworkPolicyV1
	policyType    string
}

func (f *fixByAddingPolicyToNetPol) Plan() string {
	return fmt.Sprintf("Add the '%s' policy type to the network policy", f.policyType)
}

func (f *fixByAddingPolicyToNetPol) Apply(resource k8s.Resource) []k8s.Resource {
	f.networkPolicy.Spec.PolicyTypes = append(f.networkPolicy.Spec.PolicyTypes, k8s.PolicyTypeV1(f.policyType))
	return nil
}

func newDefaultDenyNetworkPolicy(namespace string, policyList []string) k8s.Resource {
	policies := make([]k8s.PolicyTypeV1, 0, len(policyList))
	for _, policy := range policyList {
		policies = append(policies, k8s.PolicyTypeV1(policy))
	}

	networkPolicy := &k8s.NetworkPolicyV1{
		ObjectMeta: k8s.ObjectMetaV1{
			Name:      DefaultDenyNetworkPolicyName,
			Namespace: namespace,
		},
		Spec: k8s.NetworkPolicySpecV1{
			PolicyTypes: policies,
		},
	}

	networkPolicy.Kind = "NetworkPolicy"
	networkPolicy.APIVersion = "networking.k8s.io/v1"

	return networkPolicy
}
