package cmd

import (
	"testing"
)

func TestSeccompEnabledPod(t *testing.T) {
	runAuditTest(t, "seccomp_enabled_pod.yml", auditSeccomp, []int{})
}

func TestSeccompEnabled(t *testing.T) {
	runAuditTest(t, "seccomp_enabled.yml", auditSeccomp, []int{})
}

func TestSeccompAnnotationMissing(t *testing.T) {
	runAuditTest(t, "seccomp_annotation_missing.yml", auditSeccomp, []int{ErrorSeccompAnnotationMissing})
}

func TestSeccompBadValuePod(t *testing.T) {
	runAuditTest(t, "seccomp_disabled_pod.yml", auditSeccomp, []int{ErrorSeccompDisabledPod})
}

func TestSeccompBadValue(t *testing.T) {
	runAuditTest(t, "seccomp_disabled.yml", auditSeccomp, []int{ErrorSeccompDisabled})
}

func TestSeccompDeprecatedValuePod(t *testing.T) {
	runAuditTest(t, "seccomp_deprecated_pod.yml", auditSeccomp, []int{ErrorSeccompDeprecatedPod})
}

func TestSeccompDeprecatedValue(t *testing.T) {
	runAuditTest(t, "seccomp_deprecated.yml", auditSeccomp, []int{ErrorSeccompDeprecated})
}
