package cmd

import "testing"

func TestHostNetworkTrueV1(t *testing.T) {
	runAuditTest(t, "host_network_true_v1.yml", auditNamespaces, []int{ErrorNamespaceHostNetworkTrue})
}

func TestHostIPCTrueV1(t *testing.T) {
	runAuditTest(t, "host_IPC_true_v1.yml", auditNamespaces, []int{ErrorNamespaceHostIPCTrue})
}

func TestHostPIDTrueV1(t *testing.T) {
	runAuditTest(t, "host_PID_true_v1.yml", auditNamespaces, []int{ErrorNamespaceHostPIDTrue})
}
func TestHostNetworkTrueAllowedV1(t *testing.T) {
	runAuditTest(t, "host_network_true_allowed_v1.yml", auditNamespaces, []int{ErrorNamespaceHostNetworkTrueAllowed})
}

func TestHostIPCTrueAllowedV1(t *testing.T) {
	runAuditTest(t, "host_IPC_true_allowed_v1.yml", auditNamespaces, []int{ErrorNamespaceHostIPCTrueAllowed})
}

func TestHostPIDTrueAllowedV1(t *testing.T) {
	runAuditTest(t, "host_PID_true_allowed_v1.yml", auditNamespaces, []int{ErrorNamespaceHostPIDTrueAllowed})
}
func TestNamespacesMisconfiguredAllowV1(t *testing.T) {
	runAuditTest(t, "namespaces_misconfigured_allow_v1.yml", auditNamespaces, []int{ErrorMisconfiguredKubeauditAllow})
}

func TestNamespacesAllTrueV1(t *testing.T) {
	runAuditTest(t, "namespaces_all_true_v1.yml", auditNamespaces, []int{ErrorNamespaceHostPIDTrue, ErrorNamespaceHostIPCTrue, ErrorNamespaceHostNetworkTrue})
}

func TestNamespacesAllTrueAllowedV1(t *testing.T) {
	runAuditTest(t, "namespaces_all_true_allowed_v1.yml", auditNamespaces, []int{ErrorNamespaceHostPIDTrueAllowed, ErrorNamespaceHostIPCTrueAllowed, ErrorNamespaceHostNetworkTrueAllowed})
}

func TestAllowNamespacesFromConfig(t *testing.T) {
	rootConfig.auditConfig = "../configs/allow_namespaces_from_config.yml"
	runAuditTest(t, "host_network_true_v1.yml", auditNamespaces, []int{ErrorNamespaceHostNetworkTrueAllowed, ErrorMisconfiguredKubeauditAllow})
	runAuditTest(t, "host_IPC_true_v1.yml", auditNamespaces, []int{ErrorNamespaceHostIPCTrueAllowed, ErrorMisconfiguredKubeauditAllow})
	runAuditTest(t, "host_PID_true_v1.yml", auditNamespaces, []int{ErrorNamespaceHostPIDTrueAllowed, ErrorMisconfiguredKubeauditAllow})
	rootConfig.auditConfig = ""
}
