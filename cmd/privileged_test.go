package cmd

import "testing"

func TestSecurityContextNIL_Privileged(t *testing.T) {
	runTest(t, "security_context_nil.yml", auditPrivileged, ErrorSecurityContextNIL)
}

func TestPrivilegedNIL(t *testing.T) {
	runTest(t, "privileged_nil.yml", auditPrivileged, ErrorPrivilegedNIL)
}

func TestPrivilegedTrue(t *testing.T) {
	runTest(t, "privileged_true.yml", auditPrivileged, ErrorPrivilegedTrue)
}

func TestPrivilegedTrueAllowed(t *testing.T) {
	runTest(t, "privileged_true_allowed.yml", auditPrivileged, ErrorPrivilegedTrueAllowed)
}

func TestPrivilegedMisconfiguredAllow(t *testing.T) {
	runTest(t, "privileged_misconfigured_allow.yml", auditPrivileged, ErrorMisconfiguredKubeauditAllow)
}
