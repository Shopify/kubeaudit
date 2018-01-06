package cmd

import (
	"testing"
)

func TestSecurityContextNIL_APE(t *testing.T) {
	runAuditTest(t, "security_context_nil.yml", auditAllowPrivilegeEscalation, []int{ErrorSecurityContextNIL})
}

func TestAllowPrivilegeEscalationNil(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_nil.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationNIL})
}

func TestAllowPrivilegeEscalationTrue(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_true.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationTrue})
}

func TestAllowPrivilegeEscalationTrueAllowed(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_true_allowed.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationTrueAllowed})
}

func TestAllowPrivilegeEscalationMisconfiguredAllow(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_misconfigured_allow.yml", auditAllowPrivilegeEscalation, []int{ErrorMisconfiguredKubeauditAllow})
}
