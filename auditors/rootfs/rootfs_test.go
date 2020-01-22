package rootfs

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditReadOnlyRootFilesystem(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"read_only_root_filesystem_nil_v1.yml", fixtureDir, []string{ReadOnlyRootFilesystemNil}},
		{"read_only_root_filesystem_false_v1.yml", fixtureDir, []string{ReadOnlyRootFilesystemFalse}},
		{"read_only_root_filesystem_false_allowed_v1.yml", fixtureDir, []string{override.GetOverriddenResultName(ReadOnlyRootFilesystemFalse)}},
		{"read_only_root_filesystem_redundant_override_v1.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"read_only_root_filesystem_false_allowed_multi_container_multi_labels_v1.yml", fixtureDir, []string{override.GetOverriddenResultName(ReadOnlyRootFilesystemFalse)}},
		{"read_only_root_filesystem_false_allowed_multi_container_single_label_v1.yml", fixtureDir, []string{
			override.GetOverriddenResultName(ReadOnlyRootFilesystemFalse), ReadOnlyRootFilesystemFalse,
		}},

		// Shared fixtures
		{"security_context_nil_v1.yml", test.SharedFixturesDir, []string{ReadOnlyRootFilesystemNil}},
		{"security_context_nil_v1beta1.yml", test.SharedFixturesDir, []string{ReadOnlyRootFilesystemNil}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, tt.fixtureDir, tt.file, New(), tt.expectedErrors)
		})
	}
}
