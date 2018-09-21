package cmd

import (
	"testing"
)

func TestSecurityContextNil_APE(t *testing.T) {
	runAuditTest(t, "security_context_nil.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationNil})
}

func TestAllowPrivilegeEscalationNil(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_nil.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationNil})
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
