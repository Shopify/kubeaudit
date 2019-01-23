package cmd

import (
	"testing"
)

func TestSecurityContextNilV1(t *testing.T) {
	runAuditTest(t, "security_context_nil_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootNil})
}

func TestRunAsNonRootNilV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_nil_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootNil})
}

func TestRunAsNonRootFalseV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_false_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalse})
}

func TestRunAsRootFalseAllowedV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_false_allowed_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalseAllowed})
}

func TestRunAsNonRootMisconfiguredAllowV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_misconfigured_allow_v1.yml", auditRunAsNonRoot, []int{ErrorMisconfiguredKubeauditAllow})
}

func TestPSCRunAsNonRootFalseV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalse})
}

func TestPSCTrueCSCFalseRunAsNonRootFalseV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_true_csc_false_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalse})
}

func TestPSCFalseCSCFalseRunAsNonRootFalseV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_csc_false_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalse})
}
func TestPSCRunAsRootFalseAllowedV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_allowed_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalseAllowed})
}
