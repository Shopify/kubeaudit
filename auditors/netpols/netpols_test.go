package netpols

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditDefaultDenyNetworkPolicies(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
	}{
		{"namespace_missing_default_deny_netpol.yml", []string{MissingDefaultDenyIngressAndEgressNetworkPolicy}},
		{"namespace_missing_default_deny_egress_netpol.yml", []string{MissingDefaultDenyEgressNetworkPolicy}},
		{"namespace_missing_default_deny_ingress_netpol.yml", []string{MissingDefaultDenyIngressNetworkPolicy}},
		{"namespace_has_default_deny_netpol.yml", nil},
		{"namespace_has_default_deny_and_allow_all_netpol.yml", []string{AllowAllIngressNetworkPolicyExists, AllowAllEgressNetworkPolicyExists}},
		{"namespace_missing_default_deny_netpol_allowed.yml", []string{override.GetOverriddenResultName(MissingDefaultDenyIngressAndEgressNetworkPolicy)}},
		{"namespace_missing_default_deny_egress_netpol_allowed.yml", []string{override.GetOverriddenResultName(MissingDefaultDenyEgressNetworkPolicy)}},
		{"namespace_missing_default_deny_ingress_netpol_allowed.yml", []string{override.GetOverriddenResultName(MissingDefaultDenyIngressNetworkPolicy)}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, fixtureDir, tt.file, New(), tt.expectedErrors)
		})
	}
}
