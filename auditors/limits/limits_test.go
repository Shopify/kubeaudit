package limits

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

const fixtureDir = "fixtures"

func TestAuditLimits(t *testing.T) {
	cases := []struct {
		file           string
		maxCPU         string
		maxMemory      string
		expectedErrors []string
	}{
		{"resources-limit-nil.yml", "", "", []string{LimitsNotSet}},
		{"resources-limit-no-cpu.yml", "", "", []string{LimitsCPUNotSet}},
		{"resources-limit-no-memory.yml", "", "", []string{LimitsMemoryNotSet}},
		{"resources-limit.yml", "", "", []string{}},
		{"resources-limit.yml", "600m", "", []string{LimitsCPUExceeded}},
		{"resources-limit.yml", "", "384", []string{LimitsMemoryExceeded}},
		{"resources-limit.yml", "600m", "384", []string{LimitsCPUExceeded, LimitsMemoryExceeded}},
		{"resources-limit.yml", "750m", "512Mi", []string{}},
	}

	for i, tc := range cases {
		// These lines are needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		i := i
		t.Run(fmt.Sprintf("%s %s %s", tc.file, tc.maxCPU, tc.maxMemory), func(t *testing.T) {
			t.Parallel()
			auditor, err := New(Config{CPU: tc.maxCPU, Memory: tc.maxMemory})
			assert.Nil(t, err)
			test.AuditManifest(t, fixtureDir, tc.file, auditor, tc.expectedErrors)
			test.AuditLocal(t, fixtureDir, tc.file, auditor, fmt.Sprintf("%s%d", strings.Split(tc.file, ".")[0], i), tc.expectedErrors)
		})
	}

	t.Run("Bad arguments", func(t *testing.T) {
		_, err := New(Config{CPU: "badvalue", Memory: ""})
		assert.NotNil(t, err)

		_, err = New(Config{CPU: "", Memory: "badvalue"})
		assert.NotNil(t, err)
	})
}
