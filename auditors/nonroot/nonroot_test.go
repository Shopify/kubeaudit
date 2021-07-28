package nonroot

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const fixtureDir = "fixtures"

func TestAuditRunAsNonRoot(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"run-as-non-root-nil.yml", fixtureDir, []string{RunAsNonRootPSCNilCSCNil}},
		{"run-as-non-root-false.yml", fixtureDir, []string{RunAsNonRootCSCFalse}},
		{"run-as-non-root-false-allowed.yml", fixtureDir, []string{override.GetOverriddenResultName(RunAsNonRootCSCFalse)}},
		{"run-as-non-root-redundant-override-container.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"run-as-non-root-redundant-override-pod.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"run-as-non-root-psc-false.yml", fixtureDir, []string{RunAsNonRootPSCFalseCSCNil}},
		{"run-as-non-root-psc-true.yml", fixtureDir, []string{}},
		{"run-as-non-root-psc-true-csc-false.yml", fixtureDir, []string{RunAsNonRootCSCFalse}},
		{"run-as-non-root-psc-false-csc-false.yml", fixtureDir, []string{RunAsNonRootCSCFalse}},
		{"run-as-non-root-psc-false-allowed.yml", fixtureDir, []string{override.GetOverriddenResultName(RunAsNonRootPSCFalseCSCNil)}},
		{"run-as-non-root-psc-false-csc-true.yml", fixtureDir, []string{}},
		{"run-as-non-root-psc-false-csc-nil-multiple-cont.yml", fixtureDir, []string{RunAsNonRootPSCFalseCSCNil}},
		{"run-as-non-root-psc-false-csc-true-multiple-cont.yml", fixtureDir, []string{}},
		{"run-as-non-root-psc-false-allowed-multi-containers-multi-labels.yml", fixtureDir, []string{
			override.GetOverriddenResultName(RunAsNonRootPSCFalseCSCNil),
			kubeaudit.RedundantAuditorOverride,
		}},
		{"run-as-non-root-psc-false-allowed-multi-containers-single-label.yml", fixtureDir, []string{
			kubeaudit.RedundantAuditorOverride, RunAsNonRootCSCFalse,
		}},
		{"run-as-user-0.yml", fixtureDir, []string{RunAsUserCSCRoot}},
		{"run-as-user-0-allowed.yml", fixtureDir, []string{override.GetOverriddenResultName(RunAsUserCSCRoot)}},
		{"run-as-user-psc-0.yml", fixtureDir, []string{RunAsUserPSCRoot}},
		{"run-as-user-psc-0-allowed.yml", fixtureDir, []string{override.GetOverriddenResultName(RunAsUserPSCRoot)}},
		{"run-as-user-psc-1.yml", fixtureDir, []string{}},
		{"run-as-user-psc-1-csc-0.yml", fixtureDir, []string{RunAsUserCSCRoot}},
		{"run-as-user-psc-0-csc-0.yml", fixtureDir, []string{RunAsUserCSCRoot}},
		{"run-as-user-psc-0-csc-1.yml", fixtureDir, []string{RunAsUserPSCRoot}},
		{"run-as-user-redundant-override-container.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"run-as-user-redundant-override-pod.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"run-as-user-psc-0-csc-nil-multiple-cont.yml", fixtureDir, []string{RunAsUserPSCRoot}},
		{"run-as-user-psc-0-csc-1-multiple-cont.yml", fixtureDir, []string{RunAsUserPSCRoot}},
		{"run-as-user-psc-0-allowed-multi-containers-multi-labels.yml", fixtureDir, []string{
			override.GetOverriddenResultName(RunAsUserPSCRoot),
		}},
		{"run-as-user-psc-0-allowed-multi-containers-single-label.yml", fixtureDir, []string{
			override.GetOverriddenResultName(RunAsUserCSCRoot), RunAsUserPSCRoot,
		}},
		{"run-as-user-0-run-as-non-root-true.yml", fixtureDir, []string{RunAsUserCSCRoot}},
		{"run-as-user-0-run-as-non-root-false.yml", fixtureDir, []string{RunAsUserCSCRoot}},
		{"run-as-user-psc-0-run-as-non-root-psc-true.yml", fixtureDir, []string{RunAsUserPSCRoot}},
		{"run-as-user-psc-0-run-as-non-root-psc-false.yml", fixtureDir, []string{RunAsUserPSCRoot}},
		{"run-as-user-1-run-as-non-root-true.yml", fixtureDir, []string{}},
		{"run-as-user-1-run-as-non-root-false.yml", fixtureDir, []string{}},
		{"run-as-user-psc-1-run-as-non-root-psc-true.yml", fixtureDir, []string{}},
		{"run-as-user-psc-1-run-as-non-root-psc-false.yml", fixtureDir, []string{}},
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
