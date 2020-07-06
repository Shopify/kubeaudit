package capabilities

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditCapabilities(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"capabilities-nil.yml", fixtureDir, []string{CapabilityNotDropped}},
		{"capabilities-added.yml", fixtureDir, []string{CapabilityAdded}},
		{"capabilities-added-not-dropped.yml", fixtureDir, []string{CapabilityAdded, CapabilityNotDropped}},
		{"capabilities-some-allowed.yml", fixtureDir, []string{
			CapabilityAdded,
			override.GetOverriddenResultName(CapabilityAdded),
			override.GetOverriddenResultName(CapabilityNotDropped),
		}},
		{"capabilities-some-dropped.yml", fixtureDir, []string{CapabilityNotDropped}},
		{"capabilities-dropped-all.yml", fixtureDir, []string{}},
		{"capabilities-redundant-override.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"capabilities-some-allowed-multi-containers-all-labels.yml", fixtureDir, []string{
			CapabilityAdded,
			override.GetOverriddenResultName(CapabilityAdded),
			override.GetOverriddenResultName(CapabilityNotDropped),
		}},
		{"capabilities-some-allowed-multi-containers-some-labels.yml", fixtureDir, []string{
			CapabilityAdded,
			CapabilityNotDropped,
			override.GetOverriddenResultName(CapabilityAdded),
			override.GetOverriddenResultName(CapabilityNotDropped),
		}},
		{"capabilities-some-allowed-multi-containers-mix-labels.yml", fixtureDir, []string{
			CapabilityAdded,
			override.GetOverriddenResultName(CapabilityAdded),
			override.GetOverriddenResultName(CapabilityNotDropped),
		}},
	}

	for _, tc := range cases {
		// This line is needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			test.AuditManifest(t, tc.fixtureDir, tc.file, New(Config{}), tc.expectedErrors)
			test.AuditLocal(t, tc.fixtureDir, tc.file, New(Config{}), strings.Split(tc.file, ".")[0], tc.expectedErrors)
		})
	}
}
