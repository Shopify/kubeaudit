package runasuser

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditRunAsUser(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"run-as-user-nil.yml", fixtureDir, []string{RunAsUserPSCNilCSCNil}},
		{"run-as-user-0.yml", fixtureDir, []string{RunAsUserCSCRoot}},
		{"run-as-user-0-allowed.yml", fixtureDir, []string{override.GetOverriddenResultName(RunAsUserCSCRoot)}},
		{"run-as-user-redundant-override-container.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"run-as-user-redundant-override-pod.yml", fixtureDir, []string{kubeaudit.RedundantAuditorOverride}},
		{"run-as-user-psc-0.yml", fixtureDir, []string{RunAsUserPSCRootCSCNil}},
		{"run-as-user-psc-1.yml", fixtureDir, []string{}},
		{"run-as-user-psc-1-csc-0.yml", fixtureDir, []string{RunAsUserCSCRoot}},
		{"run-as-user-psc-0-csc-0.yml", fixtureDir, []string{RunAsUserCSCRoot}},
		{"run-as-user-psc-0-allowed.yml", fixtureDir, []string{override.GetOverriddenResultName(RunAsUserPSCRootCSCNil)}},
		{"run-as-user-psc-0-csc-1.yml", fixtureDir, []string{}},
		{"run-as-user-psc-0-csc-nil-multiple-cont.yml", fixtureDir, []string{RunAsUserPSCRootCSCNil}},
		{"run-as-user-psc-0-csc-1-multiple-cont.yml", fixtureDir, []string{}},
		{"run-as-user-psc-0-allowed-multi-containers-multi-labels.yml", fixtureDir, []string{
			override.GetOverriddenResultName(RunAsUserPSCRootCSCNil),
			kubeaudit.RedundantAuditorOverride,
		}},
		{"run-as-user-psc-0-allowed-multi-containers-single-label.yml", fixtureDir, []string{
			kubeaudit.RedundantAuditorOverride, RunAsUserCSCRoot,
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
