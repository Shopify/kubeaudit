package cmd

import (
	"testing"
)

func TestSecurityContextNil_APEV1(t *testing.T) {
	runAuditTest(t, "security_context_nil_v1.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationNil})
}

func TestAllowPrivilegeEscalationNilV1(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_nil_v1.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationNil})
}

func TestAllowPrivilegeEscalationTrueV1(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_true_v1.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationTrue})
}

func TestAllowPrivilegeEscalationTrueAllowedV1(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_true_allowed_v1.yml", auditAllowPrivilegeEscalation, []int{ErrorAllowPrivilegeEscalationTrueAllowed})
}

func TestAllowPrivilegeEscalationMisconfiguredAllowV1(t *testing.T) {
	runAuditTest(t, "allow_privilege_escalation_misconfigured_allow_v1.yml", auditAllowPrivilegeEscalation, []int{ErrorMisconfiguredKubeauditAllow})
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
