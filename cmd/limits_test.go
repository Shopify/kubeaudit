package cmd

import "testing"

func TestResourcesLimitsNilV1Beta1(t *testing.T) {
	runAuditTest(t, "resources_limit_nil_v1beta1.yml", auditLimits, []int{ErrorResourcesLimitsNil})
}

func TestResourcesNoCPULimitV1Beta1(t *testing.T) {
	runAuditTest(t, "resources_limit_no_cpu_v1beta1.yml", auditLimits, []int{ErrorResourcesLimitsCPUNil})
}

func TestResourcesNoMemoryLimitV1Beta1(t *testing.T) {
	runAuditTest(t, "resources_limit_no_memory_v1beta1.yml", auditLimits, []int{ErrorResourcesLimitsMemoryNil})
}
func TestResourcesCPULimitExceededV1Beta1(t *testing.T) {
	runAuditTest(t, "resources_limit_v1beta1.yml", auditLimits, []int{ErrorResourcesLimitsCPUExceeded}, "600m", "")
}

func TestResourcesMemoryLimitExceededV1Beta1(t *testing.T) {
	runAuditTest(t, "resources_limit_v1beta1.yml", auditLimits, []int{ErrorResourcesLimitsMemoryExceeded}, "", "384")
}
