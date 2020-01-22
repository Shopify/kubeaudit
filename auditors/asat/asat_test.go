package asat

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditAutomountServiceAccountToken(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
	}{
		{"service_account_token_deprecated_v1.yml", []string{AutomountServiceAccountTokenDeprecated}},
		{"service_account_token_true_and_no_name_v1.yml", []string{AutomountServiceAccountTokenTrueAndDefaultSA}},
		{"service_account_token_nil_and_no_name_v1.yml", []string{AutomountServiceAccountTokenTrueAndDefaultSA}},
		{"service_account_token_true_allowed_v1.yml", []string{
			override.GetOverriddenResultName(AutomountServiceAccountTokenTrueAndDefaultSA)},
		},
		{"service_account_token_true_and_default_name_v1.yml", []string{AutomountServiceAccountTokenTrueAndDefaultSA}},
		{"service_account_token_true_and_no_name_multiple_resources_v1.yml", []string{AutomountServiceAccountTokenTrueAndDefaultSA}},
		{"service_account_token_deprecated_multiple_resources_v1.yml", []string{AutomountServiceAccountTokenDeprecated}},
		{"service_account_token_redundant_override_v1.yml", []string{kubeaudit.RedundantAuditorOverride}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, fixtureDir, tt.file, New(), tt.expectedErrors)
		})
	}
}
