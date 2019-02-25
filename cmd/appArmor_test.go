package cmd

import (
	"testing"
)

func TestAppArmorEnabledV1(t *testing.T) {
	runAuditTest(t, "apparmor_enabled_v1.yml", auditAppArmor, []int{})
}

func TestAppArmorAnnotationMissingV1(t *testing.T) {
	runAuditTest(t, "apparmor_annotation_missing_v1.yml", auditAppArmor, []int{ErrorAppArmorAnnotationMissing})
}

func TestAppArmorBadValueV1(t *testing.T) {
	runAuditTest(t, "apparmor_disabled_v1.yml", auditAppArmor, []int{ErrorAppArmorDisabled})
}
