package cmd

import "testing"

func TestSecurityContextNil_PrivilegedV1(t *testing.T) {
	runAuditTest(t, "security_context_nil_v1.yml", auditPrivileged, []int{ErrorPrivilegedNil})
}

func TestPrivilegedNilV1(t *testing.T) {
	runAuditTest(t, "privileged_nil_v1.yml", auditPrivileged, []int{ErrorPrivilegedNil})
}

func TestPrivilegedTrueV1(t *testing.T) {
	runAuditTest(t, "privileged_true_v1.yml", auditPrivileged, []int{ErrorPrivilegedTrue})
}

func TestPrivilegedTrueAllowedV1(t *testing.T) {
	runAuditTest(t, "privileged_true_allowed_v1.yml", auditPrivileged, []int{ErrorPrivilegedTrueAllowed})
}

func TestPrivilegedMisconfiguredAllowV1(t *testing.T) {
	runAuditTest(t, "privileged_misconfigured_allow_v1.yml", auditPrivileged, []int{ErrorMisconfiguredKubeauditAllow})
}

func TestPrivilegedTrueAllowedMultiContainerMultiLabelsV1(t *testing.T) {
	runAuditTest(t, "privileged_true_allowed_multi_containers_multi_labels_v1.yml", auditPrivileged, []int{ErrorPrivilegedTrueAllowed})
}

func TestPrivilegedTrueAllowedMultiContainerSingleLabelV1(t *testing.T) {
	runAuditTest(t, "privileged_true_allowed_multi_containers_single_label_v1.yml", auditPrivileged, []int{ErrorPrivilegedTrueAllowed, ErrorPrivilegedTrue})
}
