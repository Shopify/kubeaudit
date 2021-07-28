package hostns

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
)

func TestFixHostNamespaces(t *testing.T) {
	cases := []struct {
		file                string
		expectedHostNetwork bool
		expectedHostIPC     bool
		expectedHostPID     bool
	}{
		{"host-network-true.yml", false, false, false},
		{"host-ipc-true.yml", false, false, false},
		{"host-pid-true.yml", false, false, false},
		{"host-network-true-allowed.yml", true, false, false},
		{"host-ipc-true-allowed.yml", false, true, false},
		{"host-pid-true-allowed.yml", false, false, true},
		{"namespaces-redundant-override.yml", false, false, false},
		{"namespaces-all-true.yml", false, false, false},
		{"namespaces-all-true-allowed.yml", true, true, true},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, fixtureDir, tc.file, New())
			for _, resource := range resources {
				podSpec := k8s.GetPodSpec(resource)
				assert.Equal(t, tc.expectedHostNetwork, podSpec.HostNetwork)
				assert.Equal(t, tc.expectedHostIPC, podSpec.HostIPC)
				assert.Equal(t, tc.expectedHostPID, podSpec.HostPID)
			}
		})
	}
}
