package limits

import (
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
		{"resources_limit_nil_v1beta1.yml", "", "", []string{LimitsNotSet}},
		{"resources_limit_no_cpu_v1beta1.yml", "", "", []string{LimitsCPUNotSet}},
		{"resources_limit_no_memory_v1beta1.yml", "", "", []string{LimitsMemoryNotSet}},
		{"resources_limit_v1beta1.yml", "", "", []string{}},
		{"resources_limit_v1beta1.yml", "600m", "", []string{LimitsCPUExceeded}},
		{"resources_limit_v1beta1.yml", "", "384", []string{LimitsMemoryExceeded}},
		{"resources_limit_v1beta1.yml", "600m", "384", []string{LimitsCPUExceeded, LimitsMemoryExceeded}},
		{"resources_limit_v1beta1.yml", "750m", "512Mi", []string{}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			auditor, err := New(Config{CPU: tt.maxCPU, Memory: tt.maxMemory})
			assert.Nil(t, err)
			test.Audit(t, fixtureDir, tt.file, auditor, tt.expectedErrors)
		})
	}

	t.Run("Bad arguments", func(t *testing.T) {
		_, err := New(Config{CPU: "badvalue", Memory: ""})
		assert.NotNil(t, err)

		_, err = New(Config{CPU: "", Memory: "badvalue"})
		assert.NotNil(t, err)
	})
}
