package cmd

import "testing"

func TestResourcesLimitsNil(t *testing.T) {
	runAuditTest(t, "resources_limit_nil.yml", auditLimits, []int{ErrorResourcesLimitsNil})
}

func TestResourcesNoCPULimit(t *testing.T) {
	runAuditTest(t, "resources_limit_no_cpu.yml", auditLimits, []int{ErrorResourcesLimitsCPUNil})
}

func TestResourcesNoMemoryLimit(t *testing.T) {
	runAuditTest(t, "resources_limit_no_memory.yml", auditLimits, []int{ErrorResourcesLimitsMemoryNil})
}
func TestResourcesCPULimitExceeded(t *testing.T) {
	runAuditTest(t, "resources_limit.yml", auditLimits, []int{ErrorResourcesLimitsCPUExceeded}, "600m", "")
}

func TestResourcesMemoryLimitExceeded(t *testing.T) {
	runAuditTest(t, "resources_limit.yml", auditLimits, []int{ErrorResourcesLimitsMemoryExceeded}, "", "384")
}
