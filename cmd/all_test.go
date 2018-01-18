package cmd

import (
	"testing"
)

func TestAuditAll(t *testing.T) {
	requiredErrors := []int{
		ErrorAllowPrivilegeEscalationNIL, ErrorAutomountServiceAccountTokenNILAndNoName, ErrorCapabilityNotDropped, ErrorImageTagMissing,
		ErrorPrivilegedNIL, ErrorReadOnlyRootFilesystemNIL, ErrorResourcesLimitsNIL, ErrorRunAsNonRootNIL,
	}
	runAuditTest(t, "audit_all.yml", mergeAuditFunctions(allAuditFunctions), requiredErrors)
}
