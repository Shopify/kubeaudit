package cmd

import (
	"testing"

	"github.com/Shopify/kubeaudit/fakeaudit"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeStatefulSetRORF")
	fakeaudit.CreateFakeStatefulSetReadOnlyRootFilesystem("fakeStatefulSetRORF")
}

func TestStatefulSetRORF(t *testing.T) {
	fakeStatefulSets := fakeaudit.GetStatefulSets("fakeStatefulSetRORF")
	results := auditReadOnlyRootFS(kubeAuditStatefulSets{list: fakeStatefulSets})

	if len(results) != 3 {
		t.Error("Test 1: Failed to catch all bad configurations")
	}

	for _, result := range results {
		if result.Name == "fakeStatefulSetRORF1" && result.Occurrences[0].id != ErrorSecurityContextNIL {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeStatefulSetRORF1.yml")
		}

		if result.Name == "fakeStatefulSetRORF2" && result.Occurrences[0].id != ErrorReadOnlyRootFilesystemNIL {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakeStatefulSetRORF2.yml")
		}

		if result.Name == "fakeStatefulSetRORF3" && result.Occurrences[0].id != ErrorReadOnlyRootFilesystemFalse {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakeStatefulSetRORF3.yml")
		}
	}
}
