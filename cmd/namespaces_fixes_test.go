package cmd

import (
	"testing"
)

func TestFixHostNetworkTrueV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "host_network_true_v1.yml", auditNamespaces)
	switch kubeType := resource.(type) {
	case *PodV1:
		assert.False(kubeType.Spec.HostNetwork)
	}
}
func TestFixHostIPCTrueV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "host_IPC_true_v1.yml", auditNamespaces)
	switch kubeType := resource.(type) {
	case *PodV1:
		assert.False(kubeType.Spec.HostNetwork)
	}
}
func TestFixHostPIDTrueV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "host_PID_true_v1.yml", auditNamespaces)
	switch kubeType := resource.(type) {
	case *PodV1:
		assert.False(kubeType.Spec.HostNetwork)
	}
}

func TestFixHostNetworkTrueAllowedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "host_network_true_allowed_v1.yml", auditNamespaces)
	switch kubeType := resource.(type) {
	case *PodV1:
		assert.True(kubeType.Spec.HostNetwork)
	}
}
func TestFixHostIPCTrueAllowedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "host_IPC_true_allowed_v1.yml", auditNamespaces)
	switch kubeType := resource.(type) {
	case *PodV1:
		assert.True(kubeType.Spec.HostIPC)
	}
}
func TestFixHostPIDTrueAllowedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "host_PID_true_allowed_v1.yml", auditNamespaces)
	switch kubeType := resource.(type) {
	case *PodV1:
		assert.True(kubeType.Spec.HostPID)
	}
}

func TestFixNamespacesMisconfiguredAllowV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "namespaces_misconfigured_allow_v1.yml", auditNamespaces)
	switch kubeType := resource.(type) {
	case *PodV1:
		assert.False(kubeType.Spec.HostNetwork)
	}
}
func TestFixNamespacesAllTrueV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "namespaces_all_true_v1.yml", auditNamespaces)
	switch kubeType := resource.(type) {
	case *PodV1:
		assert.False(kubeType.Spec.HostNetwork)
		assert.False(kubeType.Spec.HostPID)
		assert.False(kubeType.Spec.HostIPC)
	}
}
func TestFixNamespacesAllTrueAllowedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "namespaces_all_true_allowed_v1.yml", auditNamespaces)
	switch kubeType := resource.(type) {
	case *PodV1:
		assert.True(kubeType.Spec.HostNetwork)
		assert.True(kubeType.Spec.HostPID)
		assert.True(kubeType.Spec.HostIPC)
	}
}
