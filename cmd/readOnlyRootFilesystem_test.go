package cmd

import (
	"testing"

	"github.com/Shopify/kubeaudit/fakeaudit"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeDeploymentRORF")
	fakeaudit.CreateFakeDeploymentReadOnlyRootFilesystem("fakeDeploymentRORF")
	fakeaudit.CreateFakeNamespace("fakeStatefulSetRORF")
	fakeaudit.CreateFakeStatefulSetReadOnlyRootFilesystem("fakeStatefulSetRORF")
	fakeaudit.CreateFakeNamespace("fakeDaemonSetRORF")
	fakeaudit.CreateFakeDaemonSetReadOnlyRootFilesystem("fakeDaemonSetRORF")
	fakeaudit.CreateFakeNamespace("fakePodRORF")
	fakeaudit.CreateFakePodReadOnlyRootFilesystem("fakePodRORF")
	fakeaudit.CreateFakeNamespace("fakeReplicationControllerRORF")
	fakeaudit.CreateFakeReplicationControllerReadOnlyRootFilesystem("fakeReplicationControllerRORF")
}

func TestDeploymentRORF(t *testing.T) {
	fakeDeployments := fakeaudit.GetDeployments("fakeDeploymentRORF")
	wg.Add(1)
	results := auditReadOnlyRootFS(kubeAuditDeployments{list: fakeDeployments})

	if len(results) != 3 {
		t.Error("Test 1: Failed to catch all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeDeploymentRORF1" && result.err != 3 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeDeploymentRORF1.yml")
		}

		if result.name == "fakeDeploymentRORF2" && result.err != 1 {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakeDeploymentRORF2.yml")
		}

		if result.name == "fakeDeploymentRORF3" && result.err != 2 {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakeDeploymentRORF3.yml")
		}
	}
}

func TestStatefulSetRORF(t *testing.T) {
	fakeStatefulSets := fakeaudit.GetStatefulSets("fakeStatefulSetRORF")
	wg.Add(1)
	results := auditReadOnlyRootFS(kubeAuditStatefulSets{list: fakeStatefulSets})

	if len(results) != 3 {
		t.Error("Test 1: Failed to catch all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeStatefulSetRORF1" && result.err != 3 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeStatefulSetRORF1.yml")
		}

		if result.name == "fakeStatefulSetRORF2" && result.err != 1 {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakeStatefulSetRORF2.yml")
		}

		if result.name == "fakeStatefulSetRORF3" && result.err != 2 {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakeStatefulSetRORF3.yml")
		}
	}
}

func TestDaemonSetRORF(t *testing.T) {
	fakeDaemonSets := fakeaudit.GetDaemonSets("fakeDaemonSetRORF")
	wg.Add(1)
	results := auditReadOnlyRootFS(kubeAuditDaemonSets{list: fakeDaemonSets})

	if len(results) != 3 {
		t.Error("Test 1: Failed to catch all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeDaemonSetRORF1" && result.err != 3 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeDaemonSetRORF1.yml")
		}

		if result.name == "fakeDaemonSetRORF2" && result.err != 1 {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakeDaemonSetRORF2.yml")
		}

		if result.name == "fakeDaemonSetRORF3" && result.err != 2 {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakeDaemonSetRORF3.yml")
		}
	}
}

func TestPodRORF(t *testing.T) {
	fakePods := fakeaudit.GetPods("fakePodRORF")
	wg.Add(1)
	results := auditReadOnlyRootFS(kubeAuditPods{list: fakePods})

	if len(results) != 3 {
		t.Error("Test 1: Failed to catch all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakePodRORF1" && result.err != 3 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakePodRORF1.yml")
		}

		if result.name == "fakePodRORF2" && result.err != 1 {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakePodRORF2.yml")
		}

		if result.name == "fakePodRORF3" && result.err != 2 {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakePodRORF3.yml")
		}
	}
}

func TestReplicationControllerRORF(t *testing.T) {
	fakeReplicationControllers := fakeaudit.GetReplicationControllers("fakeReplicationControllerRORF")
	wg.Add(1)
	results := auditReadOnlyRootFS(kubeAuditReplicationControllers{list: fakeReplicationControllers})

	if len(results) != 3 {
		t.Error("Test 1: Failed to catch all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeReplicationControllerRORF1" && result.err != 3 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeReplicationControllerRORF1.yml")
		}

		if result.name == "fakeReplicationControllerRORF2" && result.err != 1 {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakeReplicationControllerRORF2.yml")
		}

		if result.name == "fakeReplicationControllerRORF3" && result.err != 2 {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakeReplicationControllerRORF3.yml")
		}
	}
}
