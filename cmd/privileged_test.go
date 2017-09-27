package cmd

import (
	"github.com/Shopify/kubeaudit/fakeaudit"
	"testing"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeDeploymentPrivileged")
	fakeaudit.CreateFakeDeploymentPrivileged("fakeDeploymentPrivileged")
	fakeaudit.CreateFakeNamespace("fakeStatefulSetPrivileged")
	fakeaudit.CreateFakeStatefulSetPrivileged("fakeStatefulSetPrivileged")
	fakeaudit.CreateFakeNamespace("fakeDaemonSetPrivileged")
	fakeaudit.CreateFakeDaemonSetPrivileged("fakeDaemonSetPrivileged")
	fakeaudit.CreateFakeNamespace("fakePodPrivileged")
	fakeaudit.CreateFakePodPrivileged("fakePodPrivileged")
	fakeaudit.CreateFakeNamespace("fakeReplicationControllerPrivileged")
	fakeaudit.CreateFakeReplicationControllerPrivileged("fakeReplicationControllerPrivileged")
}

func TestDeploymentPrivileged(t *testing.T) {
	fakeDeployments := fakeaudit.GetDeployments("fakeDeploymentPrivileged")
	wg.Add(1)
	results := auditPrivileged(kubeAuditDeployments{list: fakeDeployments})

	if len(results) != 2 {
		t.Error("Test 1: Failed to caught all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeDeploymentPrivileged1" && result.err != 2 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeDeploymentPrivileged1.yml")
		}

		if result.name == "fakeDeploymentPrivileged2" && result.err != 1 {
			t.Error("Test 3: Failed to identify Privileged was set to true. Refer: fakeDeploymentPrivileged2.yml")
		}
	}
}

func TestStatefulSetPrivileged(t *testing.T) {
	fakeStatefulSets := fakeaudit.GetStatefulSets("fakeStatefulSetPrivileged")
	wg.Add(1)
	results := auditPrivileged(kubeAuditStatefulSets{list: fakeStatefulSets})

	if len(results) != 2 {
		t.Error("Test 1: Failed to caught all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeStatefulSetPrivileged1" && result.err != 2 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeStatefulSetPrivileged1.yml")
		}

		if result.name == "fakeStatefulSetPrivileged2" && result.err != 1 {
			t.Error("Test 3: Failed to identify Privileged was set to true. Refer: fakeStatefulSetPrivileged2.yml")
		}
	}
}

func TestDaemonSetPrivileged(t *testing.T) {
	fakeDaemonSets := fakeaudit.GetDaemonSets("fakeDaemonSetPrivileged")
	wg.Add(1)
	results := auditPrivileged(kubeAuditDaemonSets{list: fakeDaemonSets})

	if len(results) != 2 {
		t.Error("Test 1: Failed to caught all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeDaemonSetPrivileged1" && result.err != 2 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeDaemonSetPrivileged1.yml")
		}

		if result.name == "fakeDaemonSetPrivileged2" && result.err != 1 {
			t.Error("Test 3: Failed to identify Privileged was set to true. Refer: fakeDaemonSetPrivileged2.yml")
		}
	}
}

func TestPodPrivileged(t *testing.T) {
	fakePods := fakeaudit.GetPods("fakePodPrivileged")
	wg.Add(1)
	results := auditPrivileged(kubeAuditPods{list: fakePods})

	if len(results) != 2 {
		t.Error("Test 1: Failed to caught all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakePodPrivileged1" && result.err != 2 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakePodPrivileged1.yml")
		}

		if result.name == "fakePodPrivileged2" && result.err != 1 {
			t.Error("Test 3: Failed to identify Privileged was set to true. Refer: fakePodPrivileged2.yml")
		}
	}
}

func TestReplicationControllerPrivileged(t *testing.T) {
	fakeReplicationControllers := fakeaudit.GetReplicationControllers("fakeReplicationControllerPrivileged")
	wg.Add(1)
	results := auditPrivileged(kubeAuditReplicationControllers{list: fakeReplicationControllers})

	if len(results) != 2 {
		t.Error("Test 1: Failed to caught all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeReplicationControllerPrivileged1" && result.err != 2 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeReplicationControllerPrivileged1.yml")
		}

		if result.name == "fakeReplicationControllerPrivileged2" && result.err != 1 {
			t.Error("Test 3: Failed to identify Privileged was set to true. Refer: fakeReplicationControllerPrivileged2.yml")
		}
	}
}
