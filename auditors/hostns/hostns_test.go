package hostns

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/internal/test"
)

const fixtureDir = "fixtures"

func TestAuditHostNamespaces(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
	}{
		{"host_network_true_v1.yml", []string{NamespaceHostNetworkTrue}},
		{"host_IPC_true_v1.yml", []string{NamespaceHostIPCTrue}},
		{"host_PID_true_v1.yml", []string{NamespaceHostPIDTrue}},
		{"host_network_true_allowed_v1.yml", []string{override.GetOverriddenResultName(NamespaceHostNetworkTrue)}},
		{"host_IPC_true_allowed_v1.yml", []string{override.GetOverriddenResultName(NamespaceHostIPCTrue)}},
		{"host_PID_true_allowed_v1.yml", []string{override.GetOverriddenResultName(NamespaceHostPIDTrue)}},
		{"namespaces_redundant_override_v1.yml", []string{kubeaudit.RedundantAuditorOverride}},
		{"namespaces_all_true_v1.yml", []string{NamespaceHostNetworkTrue, NamespaceHostIPCTrue, NamespaceHostPIDTrue}},
		{"namespaces_all_true_allowed_v1.yml", []string{
			override.GetOverriddenResultName(NamespaceHostNetworkTrue),
			override.GetOverriddenResultName(NamespaceHostIPCTrue),
			override.GetOverriddenResultName(NamespaceHostPIDTrue),
		}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, fixtureDir, tt.file, New(), tt.expectedErrors)
		})
	}
}
