package cmd

import (
	"github.com/Shopify/kubeaudit/fakeaudit"
	"testing"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeDeploymentRANR")
	fakeaudit.CreateFakeDeploymentRunAsNonRoot("fakeDeploymentRANR")
	fakeaudit.CreateFakeNamespace("fakeStatefulSetRANR")
	fakeaudit.CreateFakeStatefulSetRunAsNonRoot("fakeStatefulSetRANR")
	fakeaudit.CreateFakeNamespace("fakeDaemonSetRANR")
	fakeaudit.CreateFakeDaemonSetRunAsNonRoot("fakeDaemonSetRANR")
	fakeaudit.CreateFakeNamespace("fakePodRANR")
	fakeaudit.CreateFakePodRunAsNonRoot("fakePodRANR")
	fakeaudit.CreateFakeNamespace("fakeReplicationControllerRANR")
	fakeaudit.CreateFakeReplicationControllerRunAsNonRoot("fakeReplicationControllerRANR")
}

func TestDeploymentRANR(t *testing.T) {
	fakeDeployments := fakeaudit.GetDeployments("fakeDeploymentRANR")
	wg.Add(1)
	results := auditRunAsNonRoot(kubeAuditDeployments{list: fakeDeployments})

	if len(results) != 3 {
		t.Error("Test 1: Failed to caught all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeDeploymentRANR1" && result.err != 3 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeDeploymentRANR1.yml")
		}

		if result.name == "fakeDeploymentRANR2" && result.err != 1 {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakeDeploymentRANR2.yml")
		}

		if result.name == "fakeDeploymentRANR3" && result.err != 2 {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakeDeploymentRANR3.yml")
		}
	}
}

func TestStatefulSetRANR(t *testing.T) {
	fakeStatefulSets := fakeaudit.GetStatefulSets("fakeStatefulSetRANR")
	wg.Add(1)
	results := auditRunAsNonRoot(kubeAuditStatefulSets{list: fakeStatefulSets})

	if len(results) != 3 {
		t.Error("Test 1: Failed to caught all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeStatefulSetRANR1" && result.err != 3 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeStatefulSetRANR1.yml")
		}

		if result.name == "fakeStatefulSetRANR2" && result.err != 1 {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakeStatefulSetRANR2.yml")
		}

		if result.name == "fakeStatefulSetRANR3" && result.err != 2 {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakeStatefulSetRANR3.yml")
		}
	}
}

func TestDaemonSetRANR(t *testing.T) {
	fakeDaemonSets := fakeaudit.GetDaemonSets("fakeDaemonSetRANR")
	wg.Add(1)
	results := auditRunAsNonRoot(kubeAuditDaemonSets{list: fakeDaemonSets})

	if len(results) != 3 {
		t.Error("Test 1: Failed to caught all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeDaemonSetRANR1" && result.err != 3 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeDaemonSetRANR1.yml")
		}

		if result.name == "fakeDaemonSetRANR2" && result.err != 1 {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakeDaemonSetRANR2.yml")
		}

		if result.name == "fakeDaemonSetRANR3" && result.err != 2 {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakeDaemonSetRANR3.yml")
		}
	}
}

func TestPodRANR(t *testing.T) {
	fakePods := fakeaudit.GetPods("fakePodRANR")
	wg.Add(1)
	results := auditRunAsNonRoot(kubeAuditPods{list: fakePods})

	if len(results) != 3 {
		t.Error("Test 1: Failed to caught all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakePodRANR1" && result.err != 3 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakePodRANR1.yml")
		}

		if result.name == "fakePodRANR2" && result.err != 1 {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakePodRANR2.yml")
		}

		if result.name == "fakePodRANR3" && result.err != 2 {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakePodRANR3.yml")
		}
	}
}

func TestReplicationControllerRANR(t *testing.T) {
	fakeReplicationControllers := fakeaudit.GetReplicationControllers("fakeReplicationControllerRANR")
	wg.Add(1)
	results := auditRunAsNonRoot(kubeAuditReplicationControllers{list: fakeReplicationControllers})

	if len(results) != 3 {
		t.Error("Test 1: Failed to caught all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeReplicationControllerRANR1" && result.err != 3 {
			t.Error("Test 2: Failed to identify security context missing. Refer: fakeReplicationControllerRANR1.yml")
		}

		if result.name == "fakeReplicationControllerRANR2" && result.err != 1 {
			t.Error("Test 3: Failed to identify RunAsNonRoot was nil. Refer: fakeReplicationControllerRANR2.yml")
		}

		if result.name == "fakeReplicationControllerRANR3" && result.err != 2 {
			t.Error("Test 4: Failed to identify RunAsNonRoot was false. Refer: fakeReplicationControllerRANR3.yml")
		}
	}
}
