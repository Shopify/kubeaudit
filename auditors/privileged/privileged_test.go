package privileged

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const fixtureDir = "fixtures"

func TestAuditPrivileged(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"privileged-nil.yml", fixtureDir, []string{PrivilegedNil}},
		{"privileged-true.yml", fixtureDir, []string{PrivilegedTrue}},
		{"privileged-true-allowed.yml", fixtureDir, []string{override.GetOverriddenResultName(PrivilegedTrue)}},
		{"privileged-redundant-override.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"privileged-true-allowed-multi-containers-multi-labels.yml", fixtureDir, []string{override.GetOverriddenResultName(PrivilegedTrue)}},
		{"privileged-true-allowed-multi-containers-single-label.yml", fixtureDir, []string{
			PrivilegedTrue,
			override.GetOverriddenResultName(PrivilegedTrue)},
		},
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
