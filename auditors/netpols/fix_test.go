package netpols

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestFixDefaultDenyNetworkPolicies(t *testing.T) {
	cases := []struct {
		file                   string
		expectedDenyAllIngress bool
		expectedDenyAllEgress  bool
	}{
		{"namespace-missing-default-deny-netpol.yml", true, true},
		{"namespace-missing-default-deny-egress-netpol.yml", true, true},
		{"namespace-missing-default-deny-ingress-netpol.yml", true, true},
		{"namespace-has-default-deny-netpol.yml", true, true},
		{"namespace-has-default-deny-and-allow-all-netpol.yml", true, true},
		{"namespace-missing-default-deny-netpol-allowed.yml", false, false},
		{"namespace-missing-default-deny-egress-netpol-allowed.yml", true, false},
		{"namespace-missing-default-deny-ingress-netpol-allowed.yml", false, true},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			assert := assert.New(t)
			resources, _ := test.FixSetup(t, fixtureDir, tc.file, New())
			networkPolicies := getNetworkPolicies(resources, strings.Split(tc.file, ".")[0])
			assert.Equal(tc.expectedDenyAllIngress, hasDenyAllIngress(networkPolicies))
			assert.Equal(tc.expectedDenyAllEgress, hasDenyAllEgress(networkPolicies))
		})
	}
}
