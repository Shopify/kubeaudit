package cmd

import (
	"testing"
)

func TestSecurityContextNIL_APE(t *testing.T) {
	runTest(t, "security_context_nil.yml", auditAllowPrivilegeEscalation, ErrorAllowPrivilegeEscalationNIL)
}

func TestAllowPrivilegeEscalationNil(t *testing.T) {
	runTest(t, "allow_privilege_escalation_nil.yml", auditAllowPrivilegeEscalation, ErrorAllowPrivilegeEscalationNIL)
}

func TestAllowPrivilegeEscalationTrue(t *testing.T) {
	runTest(t, "allow_privilege_escalation_true.yml", auditAllowPrivilegeEscalation, ErrorAllowPrivilegeEscalationTrue)
}

func TestAllowPrivilegeEscalationTrueAllowed(t *testing.T) {
	runTest(t, "allow_privilege_escalation_true_allowed.yml", auditAllowPrivilegeEscalation, ErrorAllowPrivilegeEscalationTrueAllowed)
}

func TestAllowPrivilegeEscalationMisconfiguredAllow(t *testing.T) {
	runTest(t, "allow_privilege_escalation_misconfigured_allow.yml", auditAllowPrivilegeEscalation, ErrorMisconfiguredKubeauditAllow)
}
