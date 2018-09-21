package cmd

import "testing"

func TestSecurityContextNil_Privileged(t *testing.T) {
	runAuditTest(t, "security_context_nil.yml", auditPrivileged, []int{ErrorPrivilegedNil})
}

func TestPrivilegedNil(t *testing.T) {
	runAuditTest(t, "privileged_nil.yml", auditPrivileged, []int{ErrorPrivilegedNil})
}

func TestPrivilegedTrue(t *testing.T) {
	runAuditTest(t, "privileged_true.yml", auditPrivileged, []int{ErrorPrivilegedTrue})
}

func TestPrivilegedTrueAllowed(t *testing.T) {
	runAuditTest(t, "privileged_true_allowed.yml", auditPrivileged, []int{ErrorPrivilegedTrueAllowed})
}

func TestPrivilegedMisconfiguredAllow(t *testing.T) {
	runAuditTest(t, "privileged_misconfigured_allow.yml", auditPrivileged, []int{ErrorMisconfiguredKubeauditAllow})
}
