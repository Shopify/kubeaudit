package hostns

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/override"
)

const fixtureDir = "fixtures"

func TestAuditHostNamespaces(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
	}{
		{"host-network-true.yml", []string{NamespaceHostNetworkTrue}},
		{"host-ipc-true.yml", []string{NamespaceHostIPCTrue}},
		{"host-pid-true.yml", []string{NamespaceHostPIDTrue}},
		{"host-network-true-allowed.yml", []string{override.GetOverriddenResultName(NamespaceHostNetworkTrue)}},
		{"host-ipc-true-allowed.yml", []string{override.GetOverriddenResultName(NamespaceHostIPCTrue)}},
		{"host-pid-true-allowed.yml", []string{override.GetOverriddenResultName(NamespaceHostPIDTrue)}},
		{"namespaces-redundant-override.yml", []string{kubeaudit.RedundantAuditorOverride}},
		{"namespaces-all-true.yml", []string{NamespaceHostNetworkTrue, NamespaceHostIPCTrue, NamespaceHostPIDTrue}},
		{"namespaces-all-true-allowed.yml", []string{
			override.GetOverriddenResultName(NamespaceHostNetworkTrue),
			override.GetOverriddenResultName(NamespaceHostIPCTrue),
			override.GetOverriddenResultName(NamespaceHostPIDTrue),
		}},
	}

	for _, tc := range cases {
		// This line is needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			test.AuditManifest(t, fixtureDir, tc.file, New(), tc.expectedErrors)
			test.AuditLocal(t, fixtureDir, tc.file, New(), strings.Split(tc.file, ".")[0], tc.expectedErrors)
		})
	}
}
