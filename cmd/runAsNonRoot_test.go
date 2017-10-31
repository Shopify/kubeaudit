package cmd

import (
	"testing"

	"github.com/Shopify/kubeaudit/fakeaudit"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeDeploymentRANR")
	fakeaudit.CreateFakeDeploymentRunAsNonRoot("fakeDeploymentRANR")
}

func TestDeploymentRANR(t *testing.T) {
	fakeDeployments := fakeaudit.GetDeployments("fakeDeploymentRANR")
	results := auditRunAsNonRoot(kubeAuditDeployments{list: fakeDeployments})

	if len(results) != 3 {
		t.Error("Test 1: Failed to catch all bad configurations")
	}

	for _, result := range results {
		if result.Name == "fakeDeploymentRANR1" && result.Occurrences[0].id != ErrorSecurityContextNIL {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeDeploymentRANR1.yml")
		}

		if result.Name == "fakeDeploymentRANR2" && result.Occurrences[0].id != ErrorRunAsNonRootNIL {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakeDeploymentRANR2.yml")
		}

		if result.Name == "fakeDeploymentRANR3" && result.Occurrences[0].id != ErrorRunAsNonRootFalse {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakeDeploymentRANR3.yml")
		}
	}
}
