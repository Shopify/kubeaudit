package cmd

import "testing"

func TestSecurityContextNIL_Privileged(t *testing.T) {
	runAuditTest(t, "security_context_nil.yml", auditPrivileged, []int{ErrorPrivilegedNIL})
}

func TestPrivilegedNIL(t *testing.T) {
	runAuditTest(t, "privileged_nil.yml", auditPrivileged, []int{ErrorPrivilegedNIL})
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
