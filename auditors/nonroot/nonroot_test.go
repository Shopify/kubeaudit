package nonroot

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditRunAsNonRoot(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"run_as_non_root_nil_v1.yml", fixtureDir, []string{RunAsNonRootPSCNilCSCNil}},
		{"run_as_non_root_false_v1.yml", fixtureDir, []string{RunAsNonRootCSCFalse}},
		{"run_as_non_root_false_allowed_v1.yml", fixtureDir, []string{override.GetOverriddenResultName(RunAsNonRootCSCFalse)}},
		{"run_as_non_root_redundant_override_container_v1.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"run_as_non_root_redundant_override_pod_v1.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"run_as_non_root_psc_false_v1.yml", fixtureDir, []string{RunAsNonRootPSCFalseCSCNil}},
		{"run_as_non_root_psc_true_csc_false_v1.yml", fixtureDir, []string{RunAsNonRootCSCFalse}},
		{"run_as_non_root_psc_false_csc_false_v1.yml", fixtureDir, []string{RunAsNonRootCSCFalse}},
		{"run_as_non_root_psc_false_allowed_v1.yml", fixtureDir, []string{override.GetOverriddenResultName(RunAsNonRootPSCFalseCSCNil)}},
		{"run_as_non_root_psc_false_csc_true_v1.yml", fixtureDir, []string{}},
		{"run_as_non_root_psc_false_csc_nil_multiple_cont_v1.yml", fixtureDir, []string{RunAsNonRootPSCFalseCSCNil}},
		{"run_as_non_root_psc_false_csc_true_multiple_cont_v1.yml", fixtureDir, []string{}},
		{"run_as_non_root_psc_false_allowed_multi_containers_multi_labels_v1.yml", fixtureDir, []string{
			override.GetOverriddenResultName(RunAsNonRootPSCFalseCSCNil),
			kubeaudit.RedundantAuditorOverride,
		}},
		{"run_as_non_root_psc_false_allowed_multi_containers_single_label_v1.yml", fixtureDir, []string{
			kubeaudit.RedundantAuditorOverride, RunAsNonRootCSCFalse,
		}},

		// Shared fixtures
		{"security_context_nil_v1.yml", test.SharedFixturesDir, []string{RunAsNonRootPSCNilCSCNil}},
		{"security_context_nil_v1beta1.yml", test.SharedFixturesDir, []string{RunAsNonRootPSCNilCSCNil}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, tt.fixtureDir, tt.file, New(), tt.expectedErrors)
		})
	}
}
