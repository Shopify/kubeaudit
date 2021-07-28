package asat

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const fixtureDir = "fixtures"

func TestAuditAutomountServiceAccountToken(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
		testLocalMode  bool
	}{
		// When this yaml is applied into the cluster, both the deprecated and new service account fields are populated
		// with the service account value, so there is no error in local mode
		{"service-account-token-deprecated.yml", []string{AutomountServiceAccountTokenDeprecated}, false},
		{"service-account-token-true-and-no-name.yml", []string{AutomountServiceAccountTokenTrueAndDefaultSA}, true},
		{"service-account-token-nil-and-no-name.yml", []string{AutomountServiceAccountTokenTrueAndDefaultSA}, true},
		{"service-account-token-true-allowed.yml", []string{
			override.GetOverriddenResultName(AutomountServiceAccountTokenTrueAndDefaultSA)}, true,
		},
		{"service-account-token-true-and-default-name.yml", []string{AutomountServiceAccountTokenTrueAndDefaultSA}, true},
		{"service-account-token-false.yml", []string{}, true},
		{"service-account-token-redundant-override.yml", []string{kubeaudit.RedundantAuditorOverride}, true},
		{"service-account-token-nil-and-no-name-and-default-sa.yml", []string{}, true},
		{"service-account-token-true-and-default-sa.yml", []string{AutomountServiceAccountTokenTrueAndDefaultSA}, true},
	}

	for _, tc := range cases {
		// This line is needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			test.AuditManifest(t, fixtureDir, tc.file, New(), tc.expectedErrors)
			if tc.testLocalMode {
				test.AuditLocal(t, fixtureDir, tc.file, New(), strings.Split(tc.file, ".")[0], tc.expectedErrors)
			}
		})
	}
}
