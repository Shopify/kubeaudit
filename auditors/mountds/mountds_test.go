package mountds

import (
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
		{"docker_sock_mounted.yml", fixtureDir, []string{DockerSocketMounted}},

		// Shared fixtures
		{"security_context_nil_v1.yml", test.SharedFixturesDir, []string{}},
		{"security_context_nil_v1beta1.yml", test.SharedFixturesDir, []string{}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, tt.fixtureDir, tt.file, New(), tt.expectedErrors)
		})
	}
}
