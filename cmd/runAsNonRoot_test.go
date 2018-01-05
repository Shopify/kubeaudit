package cmd

import (
	"testing"
)

func TestSecurityContextNIL(t *testing.T) {
	runAuditTest(t, "security_context_nil.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootNIL})
}

func TestRunAsNonRootNil(t *testing.T) {
	runAuditTest(t, "run_as_non_root_nil.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootNIL})
}

func TestRunAsNonRootFalse(t *testing.T) {
	runAuditTest(t, "run_as_non_root_false.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalse})
}

func TestAllowRunAsRootFalseAllowed(t *testing.T) {
	runAuditTest(t, "run_as_non_root_false_allowed.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalseAllowed})
}

func TestAllowRunAsNonRootMisconfiguredAllow(t *testing.T) {
	runAuditTest(t, "run_as_non_root_misconfigured_allow.yml", auditRunAsNonRoot, []int{ErrorMisconfiguredKubeauditAllow})
}
