package netpols

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const fixtureDir = "fixtures"

func TestAuditDefaultDenyNetworkPolicies(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
	}{
		{"namespace-missing-default-deny-netpol.yml", []string{MissingDefaultDenyIngressAndEgressNetworkPolicy}},
		{"namespace-missing-default-deny-egress-netpol.yml", []string{MissingDefaultDenyEgressNetworkPolicy}},
		{"namespace-missing-default-deny-ingress-netpol.yml", []string{MissingDefaultDenyIngressNetworkPolicy}},
		{"namespace-has-default-deny-netpol.yml", nil},
		{"namespace-has-default-deny-and-allow-all-netpol.yml", []string{AllowAllIngressNetworkPolicyExists, AllowAllEgressNetworkPolicyExists}},
		{"namespace-missing-default-deny-netpol-allowed.yml", []string{override.GetOverriddenResultName(MissingDefaultDenyIngressAndEgressNetworkPolicy)}},
		{"namespace-missing-default-deny-egress-netpol-allowed.yml", []string{override.GetOverriddenResultName(MissingDefaultDenyEgressNetworkPolicy)}},
		{"namespace-missing-default-deny-ingress-netpol-allowed.yml", []string{override.GetOverriddenResultName(MissingDefaultDenyIngressNetworkPolicy)}},
	}

	for _, tc := range cases {
		// This line is needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			test.AuditManifest(t, fixtureDir, tc.file, New(), tc.expectedErrors)
			test.AuditLocal(t, fixtureDir, tc.file, New(), strings.Split(tc.file, ".")[0], tc.expectedErrors)
		})
	}
}
