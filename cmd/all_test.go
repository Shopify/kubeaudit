package cmd

import (
	"testing"
)

func TestAuditAll(t *testing.T) {
	requiredErrors := []int{
		ErrorAllowPrivilegeEscalationNil, ErrorAutomountServiceAccountTokenNilAndNoName, ErrorCapabilityNotDropped,
		ErrorImageTagMissing, ErrorPrivilegedNil, ErrorReadOnlyRootFilesystemNil, ErrorResourcesLimitsNil,
		ErrorRunAsNonRootNil, ErrorAppArmorAnnotationMissing, ErrorSeccompAnnotationMissing,
	}
	runAuditTest(t, "audit_all.yml", mergeAuditFunctions(allAuditFunctions), requiredErrors)
}

func TestAuditAllV1beta1(t *testing.T) {
	requiredErrors := []int{
		ErrorAllowPrivilegeEscalationNil, ErrorAutomountServiceAccountTokenNilAndNoName, ErrorCapabilityNotDropped,
		ErrorImageTagMissing, ErrorPrivilegedNil, ErrorReadOnlyRootFilesystemNil, ErrorResourcesLimitsNil,
		ErrorRunAsNonRootNil, ErrorAppArmorAnnotationMissing, ErrorSeccompAnnotationMissing,
	}
	runAuditTest(t, "audit_all_v1beta1.yml", mergeAuditFunctions(allAuditFunctions), requiredErrors)
}
