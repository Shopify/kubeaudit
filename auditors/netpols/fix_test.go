package netpols

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestFixDefaultDenyNetworkPolicies(t *testing.T) {
	cases := []struct {
		file string
		// There should be a catch-all network policy (either newly added or pre-existing)
		expectedHasDefaultDenyPolicy bool
		// The catch-all policy should deny all ingress
		expectedDenyAllIngress bool
		// The catch-all policy should deny all egress
		expectedDenyAllEgress bool
	}{
		{"namespace_missing_default_deny_netpol.yml", true, true, true},
		{"namespace_missing_default_deny_egress_netpol.yml", true, true, true},
		{"namespace_missing_default_deny_ingress_netpol.yml", true, true, true},
		{"namespace_has_default_deny_netpol.yml", true, true, true},
		{"namespace_has_default_deny_and_allow_all_netpol.yml", true, true, true},
		{"namespace_missing_default_deny_netpol_allowed.yml", false, false, false},
		{"namespace_missing_default_deny_egress_netpol_allowed.yml", true, true, false},
		{"namespace_missing_default_deny_ingress_netpol_allowed.yml", true, false, true},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			assert := assert.New(t)
			resources, _ := test.FixSetup(t, fixtureDir, tt.file, New())
			networkPolicies := getNetworkPolicies(resources, "default")
			hasCatchAllNetPol, networkPolicy := hasCatchAllNetworkPolicy(networkPolicies)
			assert.Equal(tt.expectedHasDefaultDenyPolicy, hasCatchAllNetPol)
			assert.Equal(tt.expectedDenyAllIngress, hasDenyAllIngress(networkPolicy))
			assert.Equal(tt.expectedDenyAllEgress, hasDenyAllEgress(networkPolicy))
		})
	}
}
