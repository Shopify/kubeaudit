package seccomp

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
)

func TestAuditSeccomp(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
	}{
		{"seccomp-profile-missing.yml", []string{SeccompProfileMissing}},
		{"seccomp-profile-missing-disabled-container.yml", []string{SeccompProfileMissing, SeccompDisabledContainer}},
		{"seccomp-disabled-pod.yml", []string{SeccompDisabledPod}},
		{"seccomp-disabled.yml", []string{SeccompDisabledContainer}},
		{"seccomp-enabled-pod.yml", nil},
		{"seccomp-enabled.yml", nil},
	}

	for _, tc := range cases {
		// This line is needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			test.AuditManifest(t, fixtureDir, tc.file, New(), tc.expectedErrors)
			test.AuditLocal(t, fixtureDir, tc.file, New(), strings.Split(tc.file, ".")[0], tc.expectedErrors)
		})
	}
}
