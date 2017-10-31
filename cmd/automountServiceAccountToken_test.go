package cmd

import (
	"testing"

	"github.com/Shopify/kubeaudit/fakeaudit"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeReplicationControllerASAT")
	fakeaudit.CreateFakeReplicationControllerAutomountServiceAccountToken("fakeReplicationControllerASAT")
}

func TestReplicationControllerASAT(t *testing.T) {
	fakeReplicationControllers := fakeaudit.GetReplicationControllers("fakeReplicationControllerASAT")
	results := auditAutomountServiceAccountToken(kubeAuditReplicationControllers{list: fakeReplicationControllers})

	if len(results) != 2 {
		t.Error("Test 1: Failed to detect all bad configuarations")
	}

	for _, result := range results {
		if result.Name == "fakeReplicationControllerASAT1" && result.Occurrences[0].id != ErrorServiceAccountTokenDeprecated {
			t.Error("Test 2: Failed to identify deprecated service account name. Refer: fakeReplicationControllerASAT1.yml")
		}

		if result.Name == "fakeReplicationControllerASAT2" && result.Occurrences[0].id != ErrorServiceAccountTokenTrueAndNoName {
			t.Error("Test 3: Failed to identify automountServiceAccountToken set to true. Refer: fakeReplicationControllerASAT2.yml")
		}
	}
}
