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
