package cmd

import "testing"

func TestResourcesLimitsNil(t *testing.T) {
	runTest(t, "resources_limit_nil.yml", auditLimits, ErrorResourcesLimitsNIL)
}

func TestResourcesNoCPULimit(t *testing.T) {
	runTest(t, "resources_limit_no_cpu.yml", auditLimits, ErrorResourcesLimitsCpuNIL)
}

func TestResourcesNoMemoryLimit(t *testing.T) {
	runTest(t, "resources_limit_no_memory.yml", auditLimits, ErrorResourcesLimitsMemoryNIL)
}
func TestResourcesCPULimitExceeded(t *testing.T) {
	runTest(t, "resources_limit.yml", auditLimits, ErrorResourcesLimitsCpuExceeded, "600m", "")
}

func TestResourcesMemoryLimitExceeded(t *testing.T) {
	runTest(t, "resources_limit.yml", auditLimits, ErrorResourcesLimitsMemoryExceeded, "", "384")
}
