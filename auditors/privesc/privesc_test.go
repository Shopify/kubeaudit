package privesc

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditPrivilegeEscalation(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"allow_privilege_escalation_nil_v1.yml", fixtureDir, []string{AllowPrivilegeEscalationNil}},
		{"allow_privilege_escalation_true_v1.yml", fixtureDir, []string{AllowPrivilegeEscalationTrue}},
		{"allow_privilege_escalation_true_allowed_v1.yml", fixtureDir, []string{override.GetOverriddenResultName(AllowPrivilegeEscalationTrue)}},
		{"allow_privilege_escalation_redundant_override_v1.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"allow_privilege_escalation_nil_v1beta1.yml", fixtureDir, []string{AllowPrivilegeEscalationNil}},
		{"allow_privilege_escalation_true_v1beta1.yml", fixtureDir, []string{AllowPrivilegeEscalationTrue}},
		{"allow_privilege_escalation_true_allowed_v1beta1.yml", fixtureDir, []string{override.GetOverriddenResultName(AllowPrivilegeEscalationTrue)}},
		{"allow_privilege_escalation_redundant_override_v1beta1.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"allow_privilege_escalation_true_multiple_allowed_multiple_containers_v1beta.yml", fixtureDir, []string{override.GetOverriddenResultName(AllowPrivilegeEscalationTrue)}},
		{"allow_privilege_escalation_true_single_allowed_multiple_containers_v1beta.yml", fixtureDir, []string{AllowPrivilegeEscalationTrue, override.GetOverriddenResultName(AllowPrivilegeEscalationTrue)}},

		// Shared fixtures
		{"security_context_nil_v1.yml", test.SharedFixturesDir, []string{AllowPrivilegeEscalationNil}},
		{"security_context_nil_v1beta1.yml", test.SharedFixturesDir, []string{AllowPrivilegeEscalationNil}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, tt.fixtureDir, tt.file, New(), tt.expectedErrors)
		})
	}
}
