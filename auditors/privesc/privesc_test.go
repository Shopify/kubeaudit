package privesc

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const fixtureDir = "fixtures"

func TestAuditPrivilegeEscalation(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"allow-privilege-escalation-nil.yml", fixtureDir, []string{AllowPrivilegeEscalationNil}},
		{"allow-privilege-escalation-redundant-override.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"allow-privilege-escalation-true-allowed.yml", fixtureDir, []string{override.GetOverriddenResultName(AllowPrivilegeEscalationTrue)}},
		{"allow-privilege-escalation-true-multi-allowed-multi-containers.yml", fixtureDir, []string{override.GetOverriddenResultName(AllowPrivilegeEscalationTrue)}},
		{"allow-privilege-escalation-true-single-allowed-multi-containers.yml", fixtureDir, []string{AllowPrivilegeEscalationTrue, override.GetOverriddenResultName(AllowPrivilegeEscalationTrue)}},
		{"allow-privilege-escalation-true.yml", fixtureDir, []string{AllowPrivilegeEscalationTrue}},
	}

	for _, tc := range cases {
		// This line is needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			test.AuditManifest(t, tc.fixtureDir, tc.file, New(), tc.expectedErrors)
			test.AuditLocal(t, tc.fixtureDir, tc.file, New(), strings.Split(tc.file, ".")[0], tc.expectedErrors)
		})
	}
}
