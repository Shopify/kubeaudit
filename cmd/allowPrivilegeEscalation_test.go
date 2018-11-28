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

func TestSecurityContextNil_APEV1Beta1(t *testing.T) {
	runAuditTest(t, "security_context_nil_v1beta1.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationNil})
}

func TestAllowPrivilegeEscalationNilV1Beta1(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_nil_v1beta1.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationNil})
}

func TestAllowPrivilegeEscalationTrueV1Beta1(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_true_v1beta1.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationTrue})
}

func TestAllowPrivilegeEscalationTrueAllowedV1Beta1(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_true_allowed_v1beta1.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationTrueAllowed})
}

func TestAllowPrivilegeEscalationMisconfiguredAllowV1Beta1(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_misconfigured_allow_v1beta1.yml", auditAllowPrivilegeEscalation, []int{ErrorMisconfiguredKubeauditAllow})
}
