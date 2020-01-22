package privileged

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditPrivileged(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"privileged_nil_v1.yml", fixtureDir, []string{PrivilegedNil}},
		{"privileged_true_v1.yml", fixtureDir, []string{PrivilegedTrue}},
		{"privileged_true_allowed_v1.yml", fixtureDir, []string{override.GetOverriddenResultName(PrivilegedTrue)}},
		{"privileged_redundant_override_v1.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"privileged_true_allowed_multi_containers_multi_labels_v1.yml", fixtureDir, []string{override.GetOverriddenResultName(PrivilegedTrue)}},
		{"privileged_true_allowed_multi_containers_single_label_v1.yml", fixtureDir, []string{
			PrivilegedTrue,
			override.GetOverriddenResultName(PrivilegedTrue)},
		},

		// Shared fixtures
		{"security_context_nil_v1.yml", test.SharedFixturesDir, []string{PrivilegedNil}},
		{"security_context_nil_v1beta1.yml", test.SharedFixturesDir, []string{PrivilegedNil}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, tt.fixtureDir, tt.file, New(), tt.expectedErrors)
		})
	}
}
