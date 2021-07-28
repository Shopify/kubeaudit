package capabilities

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const fixtureDir = "fixtures"

func TestAuditCapabilities(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"capabilities-nil.yml", fixtureDir, []string{CapabilityOrSecurityContextMissing}},
		{"capabilities-added.yml", fixtureDir, []string{CapabilityAdded}},
		{"capabilities-added-not-dropped.yml", fixtureDir, []string{CapabilityAdded, CapabilityShouldDropAll}},
		{"capabilities-some-allowed.yml", fixtureDir, []string{
			override.GetOverriddenResultName(CapabilityAdded),
			CapabilityAdded,
		}},
		{"capabilities-some-dropped.yml", fixtureDir, []string{CapabilityShouldDropAll}},
		{"capabilities-dropped-all.yml", fixtureDir, []string{}},
		{"capabilities-some-allowed-multi-containers-all-labels.yml", fixtureDir, []string{
			CapabilityAdded,
			CapabilityShouldDropAll,
			override.GetOverriddenResultName(CapabilityAdded),
		}},
		{"capabilities-some-allowed-multi-containers-some-labels.yml", fixtureDir, []string{
			CapabilityAdded,
			CapabilityShouldDropAll,
			override.GetOverriddenResultName(CapabilityAdded),
		}},
		{"capabilities-some-allowed-multi-containers-mix-labels.yml", fixtureDir, []string{
			CapabilityAdded,
			CapabilityShouldDropAll,
			override.GetOverriddenResultName(CapabilityAdded),
		}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			test.AuditManifest(t, tc.fixtureDir, tc.file, New(Config{}), tc.expectedErrors)
			test.AuditLocal(t, tc.fixtureDir, tc.file, New(Config{}), strings.Split(tc.file, ".")[0], tc.expectedErrors)
		})
	}
}
