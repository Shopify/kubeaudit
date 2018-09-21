package cmd

import (
	"testing"
)

func TestDaemonSetInNamespace(t *testing.T) {
	runAuditTestInNamespace(t, "fakeDaemonSetPrivileged", "privileged_true.yml", auditPrivileged, []int{ErrorPrivilegedTrue})
}

func TestDaemonSetNotInNamespace(t *testing.T) {
	runAuditTestInNamespace(t, "otherFakeDaemonSetPrivileged", "privileged_true.yml", auditPrivileged, []int{ErrorPrivilegedTrue})
}

func TestDeploymentInNamespace(t *testing.T) {
	runAuditTestInNamespace(t, "fakeDeploymentSC", "capabilities_some_dropped.yml", auditCapabilities, []int{ErrorCapabilityNotDropped})
}

func TestDeploymentNotInNamespace(t *testing.T) {
	runAuditTestInNamespace(t, "otherFakeDeploymentSC", "capabilities_some_dropped.yml", auditCapabilities, []int{ErrorCapabilityNotDropped})
}

func TestStatefulSetInNamespace(t *testing.T) {
	runAuditTestInNamespace(t, "fakeStatefulSetRORF", "read_only_root_filesystem_nil.yml", auditReadOnlyRootFS, []int{ErrorReadOnlyRootFilesystemNil})
}

func TestStatefulSetNotInNamespace(t *testing.T) {
	runAuditTestInNamespace(t, "otherFakeStatefulSetRORF", "read_only_root_filesystem_nil.yml", auditReadOnlyRootFS, []int{ErrorReadOnlyRootFilesystemNil})
}

func TestReplicationControllerInNamespace(t *testing.T) {
	runAuditTestInNamespace(t, "fakeReplicationControllerASAT", "service_account_token_nil_and_no_name.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenNilAndNoName})
}

func TestReplicationControllerNotInNamespace(t *testing.T) {
	runAuditTestInNamespace(t, "otherFakeReplicationControllerASAT", "service_account_token_nil_and_no_name.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenNilAndNoName})
}
