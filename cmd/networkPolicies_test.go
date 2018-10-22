package cmd

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckNamespaceNetworkPolicies(t *testing.T) {
	namespaceWithoutNetPol := apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "HasNoDefaultDeny",
		},
	}
	namespaceWithAllowAllNetPol := apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "HasNoDefaultDeny",
		},
	}
	namespaceWithDenyAllNetPol := apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "HasDefaultDeny",
		},
	}

	namespaceList := &NamespaceList{
		Items: []apiv1.Namespace{
			namespaceWithoutNetPol,
			namespaceWithAllowAllNetPol,
			namespaceWithDenyAllNetPol,
		},
	}

	defaultDenyNetPol := networking.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-deny",
			Namespace: namespaceWithDenyAllNetPol.ObjectMeta.Name,
		},
		Spec: networking.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress:     []networking.NetworkPolicyIngressRule{},
			Egress:      []networking.NetworkPolicyEgressRule{},
			PolicyTypes: []networking.PolicyType{
				"Ingress",
			},
		},
	}

	defaultAllowNetPol := networking.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-deny",
			Namespace: namespaceWithAllowAllNetPol.ObjectMeta.Name,
		},
		Spec: networking.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress: []networking.NetworkPolicyIngressRule{
				networking.NetworkPolicyIngressRule{
					Ports: []networking.NetworkPolicyPort{},
					From:  []networking.NetworkPolicyPeer{},
				},
			},
			Egress: []networking.NetworkPolicyEgressRule{},
			PolicyTypes: []networking.PolicyType{
				"Ingress",
			},
		},
	}

	netPolList := &NetworkPolicyList{
		Items: []networking.NetworkPolicy{
			defaultDenyNetPol,
			defaultAllowNetPol,
		},
	}

	expected := map[string]bool{
		namespaceWithDenyAllNetPol.ObjectMeta.Name:  false,
		namespaceWithAllowAllNetPol.ObjectMeta.Name: true,
		namespaceWithoutNetPol.ObjectMeta.Name:      true,
	}

	result := checkNamespaceNetworkPolicies(namespaceList, netPolList)

	for namespace, expectedResult := range expected {
		if result[namespace] != expectedResult {
			t.Logf("Expected %t for namespace: %s got: %t", result[namespace], namespace, expectedResult)
			t.Fail()
		}
	}
}

func TestCheckIfDefaultDenyPolicy(t *testing.T) {
	defaultDenyNetPol := networking.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-deny",
			Namespace: "default",
		},
		Spec: networking.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress:     []networking.NetworkPolicyIngressRule{},
			Egress:      []networking.NetworkPolicyEgressRule{},
			PolicyTypes: []networking.PolicyType{
				"Ingress",
			},
		},
	}

	defaultAllowNetPol := networking.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-deny",
			Namespace: "default",
		},
		Spec: networking.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress: []networking.NetworkPolicyIngressRule{
				networking.NetworkPolicyIngressRule{
					Ports: []networking.NetworkPolicyPort{},
					From:  []networking.NetworkPolicyPeer{},
				},
			},
			Egress: []networking.NetworkPolicyEgressRule{},
			PolicyTypes: []networking.PolicyType{
				"Ingress",
			},
		},
	}

	if !checkIfDefaultDenyPolicy(defaultDenyNetPol) {
		t.Logf("Expected NetPol: %v to be a default-deny policy", defaultDenyNetPol)
		t.Fail()
	}

	if checkIfDefaultDenyPolicy(defaultAllowNetPol) {
		t.Logf("Expected NetPol: %v to be a default-allow policy", defaultDenyNetPol)
		t.Fail()
	}
}
