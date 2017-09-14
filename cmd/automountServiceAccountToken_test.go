package cmd

import (
	"github.com/Shopify/kubeaudit/fakeaudit"
	"testing"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeDeploymentASAT")
	fakeaudit.CreateFakeDeploymentAutomountServiceAccountToken("fakeDeploymentASAT")
	fakeaudit.CreateFakeNamespace("fakeStatefulSetASAT")
	fakeaudit.CreateFakeStatefulSetAutomountServiceAccountToken("fakeStatefulSetASAT")
	fakeaudit.CreateFakeNamespace("fakeDaemonSetASAT")
	fakeaudit.CreateFakeDaemonSetAutomountServiceAccountToken("fakeDaemonSetASAT")
	fakeaudit.CreateFakeNamespace("fakePodASAT")
	fakeaudit.CreateFakePodAutomountServiceAccountToken("fakePodASAT")
	fakeaudit.CreateFakeNamespace("fakeReplicationControllerASAT")
	fakeaudit.CreateFakeReplicationControllerAutomountServiceAccountToken("fakeReplicationControllerASAT")
}

func TestDeploymentASAT(t *testing.T) {
	fakeDeployments := fakeaudit.GetDeployments("fakeDeploymentASAT")
	wg.Add(1)
	results := auditAutomountServiceAccountToken(kubeAuditDeployments{list: fakeDeployments})

	if len(results) != 2 {
		t.Error("Test 1: Failed to detect all bad configuarations")
	}

	for _, result := range results {
		if result.name == "fakeDeploymentASAT1" && (result.err != 3 || result.dsa == "") {
			t.Error("Test 2: Failed to identify deprecated service account name. Refer: fakeDeploymentASAT1.yml")
		}

		if result.name == "fakeDeploymentASAT2" && result.err != 2 {
			t.Error("Test 3: Failed to identify automountServiceAccountToken set to true. Refer: fakeDeploymentASAT2.yml")
		}
	}
}

func TestStatefulSetASAT(t *testing.T) {
	fakeStatefulSets := fakeaudit.GetStatefulSets("fakeStatefulSetASAT")
	wg.Add(1)
	results := auditAutomountServiceAccountToken(kubeAuditStatefulSets{list: fakeStatefulSets})

	if len(results) != 2 {
		t.Error("Test 1: Failed to detect all bad configuarations")
	}

	for _, result := range results {
		if result.name == "fakeStatefulSetASAT1" && (result.err != 3 || result.dsa == "") {
			t.Error("Test 2: Failed to identify deprecated service account name. Refer: fakeStatefulSetASAT1.yml")
		}

		if result.name == "fakeStatefulSetASAT2" && result.err != 2 {
			t.Error("Test 3: Failed to identify automountServiceAccountToken set to true. Refer: fakeStatefulSetASAT2.yml")
		}
	}
}

func TestDaemonSetASAT(t *testing.T) {
	fakeDaemonSets := fakeaudit.GetDaemonSets("fakeDaemonSetASAT")
	wg.Add(1)
	results := auditAutomountServiceAccountToken(kubeAuditDaemonSets{list: fakeDaemonSets})

	if len(results) != 2 {
		t.Error("Test 1: Failed to detect all bad configuarations")
	}

	for _, result := range results {
		if result.name == "fakeDaemonSetASAT1" && (result.err != 3 || result.dsa == "") {
			t.Error("Test 2: Failed to identify deprecated service account name. Refer: fakeDaemonSetASAT1.yml")
		}

		if result.name == "fakeDaemonSetASAT2" && result.err != 2 {
			t.Error("Test 3: Failed to identify automountServiceAccountToken set to true. Refer: fakeDaemonSetASAT2.yml")
		}
	}
}

func TestPodASAT(t *testing.T) {
	fakePods := fakeaudit.GetPods("fakePodASAT")
	wg.Add(1)
	results := auditAutomountServiceAccountToken(kubeAuditPods{list: fakePods})

	if len(results) != 2 {
		t.Error("Test 1: Failed to detect all bad configuarations")
	}

	for _, result := range results {
		if result.name == "fakePodASAT1" && (result.err != 3 || result.dsa == "") {
			t.Error("Test 2: Failed to identify deprecated service account name. Refer: fakePodASAT1.yml")
		}

		if result.name == "fakePodASAT2" && result.err != 2 {
			t.Error("Test 3: Failed to identify automountServiceAccountToken set to true. Refer: fakePodASAT2.yml")
		}
	}
}

func TestReplicationControllerASAT(t *testing.T) {
	fakeReplicationControllers := fakeaudit.GetReplicationControllers("fakeReplicationControllerASAT")
	wg.Add(1)
	results := auditAutomountServiceAccountToken(kubeAuditReplicationControllers{list: fakeReplicationControllers})

	if len(results) != 2 {
		t.Error("Test 1: Failed to detect all bad configuarations")
	}

	for _, result := range results {
		if result.name == "fakeReplicationControllerASAT1" && (result.err != 3 || result.dsa == "") {
			t.Error("Test 2: Failed to identify deprecated service account name. Refer: fakeReplicationControllerASAT1.yml")
		}

		if result.name == "fakeReplicationControllerASAT2" && result.err != 2 {
			t.Error("Test 3: Failed to identify automountServiceAccountToken set to true. Refer: fakeReplicationControllerASAT2.yml")
		}
	}
}
