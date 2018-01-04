package cmd

import (
	"testing"
)

func TestDaemonSetInNamespace(t *testing.T) {
	runTestInNamespace(t, "fakeDaemonSetPrivileged", "privileged_true.yml", auditPrivileged, ErrorPrivilegedTrue)
}

func TestDaemonSetNotInNamespace(t *testing.T) {
	runTestInNamespace(t, "otherFakeDaemonSetPrivileged", "privileged_true.yml", auditPrivileged, ErrorPrivilegedTrue)
}

func TestDeploymentInNamespace(t *testing.T) {
	runTestInNamespace(t, "fakeDeploymentSC", "capabilities_some_dropped.yml", auditCapabilities, ErrorCapabilityNotDropped)
}

func TestDeploymentNotInNamespace(t *testing.T) {
	runTestInNamespace(t, "otherFakeDeploymentSC", "capabilities_some_dropped.yml", auditCapabilities, ErrorCapabilityNotDropped)
}

func TestStatefulSetInNamespace(t *testing.T) {
	runTestInNamespace(t, "fakeStatefulSetRORF", "read_only_root_filesystem_nil.yml", auditReadOnlyRootFS, ErrorReadOnlyRootFilesystemNIL)
}

func TestStatefulSetNotInNamespace(t *testing.T) {
	runTestInNamespace(t, "otherFakeStatefulSetRORF", "read_only_root_filesystem_nil.yml", auditReadOnlyRootFS, ErrorReadOnlyRootFilesystemNIL)
}

func TestReplicationControllerInNamespace(t *testing.T) {
	runTestInNamespace(t, "fakeReplicationControllerASAT", "service_account_token_nil_and_no_name.yml", auditAutomountServiceAccountToken, ErrorAutomountServiceAccountTokenNILAndNoName)
}

func TestReplicationControllerNotInNamespace(t *testing.T) {
	runTestInNamespace(t, "otherFakeReplicationControllerASAT", "service_account_token_nil_and_no_name.yml", auditAutomountServiceAccountToken, ErrorAutomountServiceAccountTokenNILAndNoName)
}
