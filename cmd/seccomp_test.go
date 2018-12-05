package cmd

import (
	"testing"
)

func TestSeccompEnabledPodV1(t *testing.T) {
	runAuditTest(t, "seccomp_enabled_pod_v1.yml", auditSeccomp, []int{})
}

func TestSeccompEnabledV1(t *testing.T) {
	runAuditTest(t, "seccomp_enabled_v1.yml", auditSeccomp, []int{})
}

func TestSeccompAnnotationMissingV1(t *testing.T) {
	runAuditTest(t, "seccomp_annotation_missing_v1.yml", auditSeccomp, []int{ErrorSeccompAnnotationMissing})
}

func TestSeccompBadValuePodV1(t *testing.T) {
	runAuditTest(t, "seccomp_disabled_pod_v1.yml", auditSeccomp, []int{ErrorSeccompDisabledPod})
}

func TestSeccompBadValueV1(t *testing.T) {
	runAuditTest(t, "seccomp_disabled_v1.yml", auditSeccomp, []int{ErrorSeccompDisabled})
}

func TestSeccompDeprecatedValuePodV1(t *testing.T) {
	runAuditTest(t, "seccomp_deprecated_pod_v1.yml", auditSeccomp, []int{ErrorSeccompDeprecatedPod})
}

func TestSeccompDeprecatedValueV1(t *testing.T) {
	runAuditTest(t, "seccomp_deprecated_v1.yml", auditSeccomp, []int{ErrorSeccompDeprecated})
}
