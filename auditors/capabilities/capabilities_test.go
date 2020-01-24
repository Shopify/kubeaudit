package capabilities

import (
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
		{"capabilities_nil_v1beta2.yml", fixtureDir, []string{CapabilityNotDropped}},
		{"capabilities_added_v1beta2.yml", fixtureDir, []string{CapabilityAdded}},
		{"capabilities_added_not_dropped_v1beta2.yml", fixtureDir, []string{CapabilityAdded, CapabilityNotDropped}},
		{"capabilities_some_allowed_v1beta2.yml", fixtureDir, []string{
			CapabilityAdded,
			override.GetOverriddenResultName(CapabilityAdded),
			override.GetOverriddenResultName(CapabilityNotDropped),
		}},
		{"capabilities_some_dropped_v1beta2.yml", fixtureDir, []string{CapabilityNotDropped}},
		{"capabilities_dropped_all_v1beta2.yml", fixtureDir, []string{}},
		{"capabilities_redundant_override_v1beta2.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"capabilities_some_allowed_multi_containers_all_container_labels_v1beta2.yml", fixtureDir, []string{
			CapabilityAdded,
			override.GetOverriddenResultName(CapabilityAdded),
			override.GetOverriddenResultName(CapabilityNotDropped),
		}},
		{"capabilities_some_allowed_multi_containers_some_container_labels_v1beta2.yml", fixtureDir, []string{
			CapabilityAdded,
			CapabilityNotDropped,
			override.GetOverriddenResultName(CapabilityAdded),
			override.GetOverriddenResultName(CapabilityNotDropped),
		}},
		{"capabilities_some_allowed_multi_containers_mix_labels_v1beta2.yml", fixtureDir, []string{
			CapabilityAdded,
			override.GetOverriddenResultName(CapabilityAdded),
			override.GetOverriddenResultName(CapabilityNotDropped),
		}},

		// Shared fixtures
		{"security_context_nil_v1.yml", test.SharedFixturesDir, []string{CapabilityNotDropped}},
		{"security_context_nil_v1beta1.yml", test.SharedFixturesDir, []string{CapabilityNotDropped}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, tt.fixtureDir, tt.file, New(nil), tt.expectedErrors)
		})
	}
}
