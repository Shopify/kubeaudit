package cmd

import (
	"testing"
)

func TestSecurityContextNilV1(t *testing.T) {
	runAuditTest(t, "security_context_nil_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootPSCNilCSCNil})
}

func TestRunAsNonRootNilV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_nil_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootPSCNilCSCNil})
}

func TestRunAsNonRootFalseV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_false_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootPSCTrueFalseCSCFalse})
}

func TestRunAsRootFalseAllowedV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_false_allowed_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalseAllowed})
}

func TestRunAsNonRootMisconfiguredAllowContainerV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_misconfigured_allow_container_v1.yml", auditRunAsNonRoot, []int{ErrorMisconfiguredKubeauditAllow})
}
func TestRunAsNonRootMisconfiguredAllowPodV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_misconfigured_allow_pod_v1.yml", auditRunAsNonRoot, []int{ErrorMisconfiguredKubeauditAllow})
}

func TestPSCFalseCSCNilRunAsNonRootV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootPSCFalseCSCNil})
}

func TestPSCTrueCSCFalseRunAsNonRootFalseV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_true_csc_false_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootPSCTrueFalseCSCFalse})
}

func TestPSCFalseCSCFalseRunAsNonRootFalseV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_csc_false_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootPSCTrueFalseCSCFalse})
}
func TestPSCRunAsRootFalseAllowedV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_allowed_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalseAllowed})
}

func TestPSCFalseCSCTrueRunAsNonRootFalseV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_csc_true_v1.yml", auditRunAsNonRoot, []int{})
}

func TestPSCFalseCSCNilMultipleRunAsNonRootFalseV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_csc_nil_multiple_cont_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootPSCFalseCSCNil})
}

func TestPSCFalseCSCTrueMultipleRunAsNonRootFalseV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_csc_true_multiple_cont_v1.yml", auditRunAsNonRoot, []int{})
}

func TestPSCRunAsRootFalseAllowedMultiContainersV1(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_allowed_multi_containers_multi_labels_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalseAllowed, ErrorRunAsNonRootFalseAllowed})
}

func TestPSCRunAsRootFalseAllowedMultiContainersV2(t *testing.T) {
	runAuditTest(t, "run_as_non_root_psc_false_allowed_multi_containers_single_label_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootPSCTrueFalseCSCFalse, ErrorRunAsNonRootPSCTrueFalseCSCFalse})
}

func TestAllowAuditPSCRunAsRootFalseAllowedMultiContainersFromConfigV2(t *testing.T) {
	rootConfig.auditConfig = "../configs/allow_audit_from_config.yml"
	runAuditTest(t, "run_as_non_root_psc_false_allowed_multi_containers_single_label_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootPSCTrueFalseCSCFalse, ErrorRunAsNonRootPSCTrueFalseCSCFalse})
}
func TestAllowRunAsNonRootFromConfig(t *testing.T) {
	rootConfig.auditConfig = "../configs/allow_run_as_non_root_from_config.yml"
	runAuditTest(t, "security_context_nil_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalseAllowed})
	runAuditTest(t, "run_as_non_root_nil_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalseAllowed})
	runAuditTest(t, "run_as_non_root_false_v1.yml", auditRunAsNonRoot, []int{ErrorRunAsNonRootFalseAllowed})
}
