package cmd

import (
	"testing"
)

func TestSecurityContextNIL(t *testing.T) {
	runTest(t, "security_context_nil.yml", auditRunAsNonRoot, ErrorSecurityContextNIL)
}

func TestRunAsNonRootNil(t *testing.T) {
	runTest(t, "run_as_non_root_nil.yml", auditRunAsNonRoot, ErrorRunAsNonRootNIL)
}

func TestRunAsNonRootFalse(t *testing.T) {
	runTest(t, "run_as_non_root_false.yml", auditRunAsNonRoot, ErrorRunAsNonRootFalse)
}
