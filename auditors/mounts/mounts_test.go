package mounts

import (
	"github.com/Shopify/kubeaudit/internal/override"
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditDockerSockMounted(t *testing.T) {
	cases := []struct {
		file           string
		fixtureDir     string
		expectedErrors []string
	}{
		{"proc-mounted.yml", fixtureDir, []string{SensitivePathsMounted}},
		{"proc-mounted-allowed.yml", fixtureDir, []string{override.GetOverriddenResultName(SensitivePathsMounted)}},
		{"proc-mounted-allowed-multi-containers-multi-labels.yml", fixtureDir, []string{override.GetOverriddenResultName(SensitivePathsMounted)}},
		{"proc-mounted-allowed-multi-containers-single-label.yml", fixtureDir, []string{SensitivePathsMounted, override.GetOverriddenResultName(SensitivePathsMounted)}},
	}

	config := Config{SensitivePaths: DefaultSensitivePaths}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			test.AuditManifest(t, tc.fixtureDir, tc.file, New(config), tc.expectedErrors)
			test.AuditLocal(t, tc.fixtureDir, tc.file, New(config), strings.Split(tc.file, ".")[0], tc.expectedErrors)
		})
	}
}
