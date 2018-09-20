package cmd

import (
	"testing"
)

func TestAuditAll(t *testing.T) {
	requiredErrors := []int{
		ErrorAllowPrivilegeEscalationNil, ErrorAutomountServiceAccountTokenNilAndNoName, ErrorCapabilityNotDropped, ErrorImageTagMissing,
		ErrorPrivilegedNil, ErrorReadOnlyRootFilesystemNil, ErrorResourcesLimitsNil, ErrorRunAsNonRootNil,
	}
	runAuditTest(t, "audit_all.yml", mergeAuditFunctions(allAuditFunctions), requiredErrors)
}
