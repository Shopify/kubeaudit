package cmd

import (
	"testing"
)

func TestSecurityContextNil(t *testing.T) {
	runAuditTest(t, "security_context_nil.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootNil})
}

func TestRunAsNonRootNil(t *testing.T) {
	runAuditTest(t, "run_as_non_root_nil.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootNil})
}

func TestRunAsNonRootFalse(t *testing.T) {
	runAuditTest(t, "run_as_non_root_false.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalse})
}

func TestRunAsRootFalseAllowed(t *testing.T) {
	runAuditTest(t, "run_as_non_root_false_allowed.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalseAllowed})
}

func TestRunAsNonRootMisconfiguredAllow(t *testing.T) {
	runAuditTest(t, "run_as_non_root_misconfigured_allow.yml", auditRunAsNonRoot, []int{ErrorMisconfiguredKubeauditAllow})
}
