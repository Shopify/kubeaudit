package mountds

import (
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
		{"docker-sock-mounted.yml", fixtureDir, []string{DockerSocketMounted}},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			test.AuditManifest(t, tc.fixtureDir, tc.file, New(), tc.expectedErrors)
			test.AuditLocal(t, tc.fixtureDir, tc.file, New(), strings.Split(tc.file, ".")[0], tc.expectedErrors)
		})
	}
}
