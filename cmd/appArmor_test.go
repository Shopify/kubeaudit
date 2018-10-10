package cmd

import (
	"testing"
)

func TestAppArmorEnabled(t *testing.T) {
	runAuditTest(t, "apparmor_enabled.yml", auditAppArmor, []int{})
}

func TestAppArmorAnnotationMissing(t *testing.T) {
	runAuditTest(t, "apparmor_annotation_missing.yml", auditAppArmor, []int{ErrorAppArmorAnnotationMissing})
}

func TestAppArmorBadValue(t *testing.T) {
	runAuditTest(t, "apparmor_disabled.yml", auditAppArmor, []int{ErrorAppArmorDisabled})
}
