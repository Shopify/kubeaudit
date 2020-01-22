package hostns

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestFixHostNamespaces(t *testing.T) {
	cases := []struct {
		file                string
		expectedHostNetwork bool
		expectedHostIPC     bool
		expectedHostPID     bool
	}{
		{"host_network_true_v1.yml", false, false, false},
		{"host_IPC_true_v1.yml", false, false, false},
		{"host_PID_true_v1.yml", false, false, false},
		{"host_network_true_allowed_v1.yml", true, false, false},
		{"host_IPC_true_allowed_v1.yml", false, true, false},
		{"host_PID_true_allowed_v1.yml", false, false, true},
		{"namespaces_redundant_override_v1.yml", false, false, false},
		{"namespaces_all_true_v1.yml", false, false, false},
		{"namespaces_all_true_allowed_v1.yml", true, true, true},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, fixtureDir, tt.file, New())
			for _, resource := range resources {
				podSpec := k8s.GetPodSpec(resource)
				assert.Equal(t, tt.expectedHostNetwork, podSpec.HostNetwork)
				assert.Equal(t, tt.expectedHostIPC, podSpec.HostIPC)
				assert.Equal(t, tt.expectedHostPID, podSpec.HostPID)
			}
		})
	}
}
