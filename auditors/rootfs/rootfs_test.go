package rootfs

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const fixtureDir = "fixtures"

func TestAuditReadOnlyRootFilesystem(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"read-only-root-filesystem-nil.yml", fixtureDir, []string{ReadOnlyRootFilesystemNil}},
		{"read-only-root-filesystem-false.yml", fixtureDir, []string{ReadOnlyRootFilesystemFalse}},
		{"read-only-root-filesystem-false-allowed.yml", fixtureDir, []string{override.GetOverriddenResultName(ReadOnlyRootFilesystemFalse)}},
		{"read-only-root-filesystem-redundant-override.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"read-only-root-filesystem-false-allowed-multi-labels.yml", fixtureDir, []string{override.GetOverriddenResultName(ReadOnlyRootFilesystemFalse)}},
		{"read-only-root-filesystem-false-allowed-single-label.yml", fixtureDir, []string{
			override.GetOverriddenResultName(ReadOnlyRootFilesystemFalse), ReadOnlyRootFilesystemFalse,
		}},
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
