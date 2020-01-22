package netpols

import (
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit/k8stypes"
)

const DefaultDenyNetworkPolicyName = "default-deny"

type fixByAddingNetworkPolicy struct {
	policyList []string
	namespace  string
}

func (f *fixByAddingNetworkPolicy) Plan() string {
	var policyType string
	if len(f.policyList) == 2 {
		policyType = fmt.Sprintf("%s and %s", strings.ToLower(f.policyList[0]), strings.ToLower(f.policyList[1]))
	} else {
		policyType = strings.ToLower(f.policyList[0])
	}
	return fmt.Sprintf("Create a new NetworkPolicy resource which denies all %s traffic", policyType)
}

func (f *fixByAddingNetworkPolicy) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	return []k8stypes.Resource{newDefaultDenyNetworkPolicy(f.namespace, f.policyList)}
}

type fixByAddingPolicyToNetPol struct {
	networkPolicy *k8stypes.NetworkPolicyV1
	policyType    string
}

func (f *fixByAddingPolicyToNetPol) Plan() string {
	return fmt.Sprintf("Add the '%s' policy type to the network policy", f.policyType)
}

func (f *fixByAddingPolicyToNetPol) Apply(resource k8stypes.Resource) []k8stypes.Resource {
	f.networkPolicy.Spec.PolicyTypes = append(f.networkPolicy.Spec.PolicyTypes, k8stypes.PolicyTypeV1(f.policyType))
	return nil
}

func newDefaultDenyNetworkPolicy(namespace string, policyList []string) k8stypes.Resource {
	policies := make([]k8stypes.PolicyTypeV1, 0, len(policyList))
	for _, policy := range policyList {
		policies = append(policies, k8stypes.PolicyTypeV1(policy))
	}

	networkPolicy := &k8stypes.NetworkPolicyV1{
		ObjectMeta: k8stypes.ObjectMetaV1{
			Name:      DefaultDenyNetworkPolicyName,
			Namespace: namespace,
		},
		Spec: k8stypes.NetworkPolicySpecV1{
			PolicyTypes: policies,
		},
	}

	networkPolicy.Kind = "NetworkPolicy"
	networkPolicy.APIVersion = "networking.k8s.io/v1"

	return networkPolicy
}
