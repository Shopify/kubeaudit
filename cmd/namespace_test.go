package cmd

import (
	"testing"
)

func TestDaemonSetInNamespaceV1(t *testing.T) {
	runAuditTestInNamespace(t, "fakeDaemonSetPrivileged", "privileged_true_v1.yml", auditPrivileged, []int{ErrorPrivilegedTrue})
}

func TestDaemonSetNotInNamespaceV1(t *testing.T) {
	runAuditTestInNamespace(t, "otherFakeDaemonSetPrivileged", "privileged_true_v1.yml", auditPrivileged, []int{ErrorPrivilegedTrue})
}

func TestDeploymentInNamespaceV1(t *testing.T) {
	runAuditTestInNamespace(t, "fakeDeploymentSC", "capabilities_some_dropped_v1beta2.yml", auditCapabilities, []int{ErrorCapabilityNotDropped})
}

func TestDeploymentNotInNamespaceV1(t *testing.T) {
	runAuditTestInNamespace(t, "otherFakeDeploymentSC", "capabilities_some_dropped_v1beta2.yml", auditCapabilities, []int{ErrorCapabilityNotDropped})
}

func TestStatefulSetInNamespaceV1(t *testing.T) {
	runAuditTestInNamespace(t, "fakeStatefulSetRORF", "read_only_root_filesystem_nil_v1.yml", auditReadOnlyRootFS, []int{ErrorReadOnlyRootFilesystemNil})
}

func TestStatefulSetNotInNamespaceV1(t *testing.T) {
	runAuditTestInNamespace(t, "otherFakeStatefulSetRORF", "read_only_root_filesystem_nil_v1.yml", auditReadOnlyRootFS, []int{ErrorReadOnlyRootFilesystemNil})
}

func TestReplicationControllerInNamespaceV1(t *testing.T) {
	runAuditTestInNamespace(t, "fakeReplicationControllerASAT", "service_account_token_nil_and_no_name_v1.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenNilAndNoName})
}

func TestReplicationControllerNotInNamespaceV1(t *testing.T) {
	runAuditTestInNamespace(t, "otherFakeReplicationControllerASAT", "service_account_token_nil_and_no_name_v1.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenNilAndNoName})
}
