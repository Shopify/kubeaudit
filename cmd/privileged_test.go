package cmd

import (
	"testing"

	"github.com/Shopify/kubeaudit/fakeaudit"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeDaemonSetPrivileged")
	fakeaudit.CreateFakeDaemonSetPrivileged("fakeDaemonSetPrivileged")
}

func TestDaemonSetPrivileged(t *testing.T) {
	fakeDaemonSets := fakeaudit.GetDaemonSets("fakeDaemonSetPrivileged")
	wg.Add(1)
	results := auditPrivileged(kubeAuditDaemonSets{list: fakeDaemonSets})

	if len(results) != 3 {
		t.Error("Test 1: Failed to catch all bad configurations")
	}

	for _, result := range results {
		if result.Name == "fakeDaemonSetPrivileged1" && result.Occurrences[0].id != ErrorSecurityContextNIL {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeDaemonSetPrivileged1.yml")
		}

		if result.Name == "fakeDaemonSetPrivileged2" && result.Occurrences[0].id != ErrorPrivilegedTrue {
			t.Error("Test 3: Failed to identify Privileged was set to true. Refer: fakeDaemonSetPrivileged2.yml")
		}
	}
}
