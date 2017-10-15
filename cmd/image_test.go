package cmd

import (
	"github.com/Shopify/kubeaudit/fakeaudit"
	"testing"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeDeploymentImg")
	fakeaudit.CreateFakeDeploymentImg("fakeDeploymentImg")
	fakeaudit.CreateFakeNamespace("fakeStatefulSetImg")
	fakeaudit.CreateFakeStatefulSetImg("fakeStatefulSetImg")
	fakeaudit.CreateFakeNamespace("fakeDaemonSetImg")
	fakeaudit.CreateFakeDaemonSetImg("fakeDaemonSetImg")
	fakeaudit.CreateFakeNamespace("fakePodImg")
	fakeaudit.CreateFakePodImg("fakePodImg")
	fakeaudit.CreateFakeNamespace("fakeReplicationControllerImg")
	fakeaudit.CreateFakeReplicationControllerImg("fakeReplicationControllerImg")
}

func TestDeploymentImg(t *testing.T) {
	fakeDeployments := fakeaudit.GetDeployments("fakeDeploymentImg")
	wg.Add(2)
	results := auditImages("fakeContainerImg:1.5", kubeAuditDeployments{list: fakeDeployments})

	if len(results) != 2 {
		t.Error("Test 1: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.img_name == "fakeDeploymentImg1" && result.err != 1 {
			t.Error("Test 2: Failed to identify that image tag is missing. Refer: fakeDeploymentImg1.yml")
		}

		if result.img_name == "fakeDeploymentImg2" && result.img_tag != "1.5" {
			t.Error("Test 3: Failed to identify the correct image tag which is present. Refer: fakeDeploymentImg2.yml")
		}
	}

	results = auditImages("fakeContainerImg:1.6", kubeAuditDeployments{list: fakeDeployments})

	if len(results) != 2 {
		t.Error("Test 4: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.img_name == "fakeDeploymentImg1" && result.err != 1 {
			t.Error("Test 5: Failed to identify that image tag is missing. Refer: fakeDeploymentImg1.yml")
		}

		if result.img_name == "fakeDeploymentImg2" && result.err != 2 {
			t.Error("Test 6: Failed to identify wrong image tag. Refer: fakeDeploymentImg2.yml")
		}
	}

}

func TestStatefulSetImg(t *testing.T) {
	fakeStatefulSets := fakeaudit.GetStatefulSets("fakeStatefulSetImg")
	wg.Add(2)
	results := auditImages("fakeContainerImg:1.5", kubeAuditStatefulSets{list: fakeStatefulSets})

	if len(results) != 2 {
		t.Error("Test 1: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeStatefulSetImg1" && result.err != 1 {
			t.Error("Test 2: Failed to identify that image tag is missing. Refer: fakeStatefulSetImg1.yml")
		}

		if result.name == "fakeStatefulSetImg2" && result.img_tag != "1.5" {
			t.Error("Test 3: Failed to identify the correct image tag which is present. Refer: fakeStatefulSetImg2.yml")
		}
	}

	results = auditImages("fakeContainerImg:1.6", kubeAuditStatefulSets{list: fakeStatefulSets})

	if len(results) != 2 {
		t.Error("Test 4: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeStatefulSetImg1" && result.err == 2 {
			t.Error("Test 5: Failed to identify that image tag is missing. Refer: fakeStatefulSetImg1")
		}

		if result.name == "fakeStatefulSetImg2" && result.err == 1 {
			t.Error("Test 6: Failed to identify wrong image tag. Refer: fakeStatefulSetImg2")
		}
	}

}

func TestDaemonSetImg(t *testing.T) {
	fakeDaemonSets := fakeaudit.GetDaemonSets("fakeDaemonSetImg")
	wg.Add(2)
	results := auditImages("fakeContainerImg:1.5", kubeAuditDaemonSets{list: fakeDaemonSets})

	if len(results) != 2 {
		t.Error("Test 1: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeDaemonSetImg1" && result.err != 1 {
			t.Error("Test 2: Failed to identify that image tag is missing. Refer: fakeDaemonSetImg1.yml")
		}

		if result.name == "fakeDaemonSetImg2" && result.img_tag != "1.5" {
			t.Error("Test 3: Failed to identify the correct image tag which is present. Refer: fakeDaemonSetImg2.yml")
		}
	}

	results = auditImages("fakeContainerImg:1.6", kubeAuditDaemonSets{list: fakeDaemonSets})

	if len(results) != 2 {
		t.Error("Test 4: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeDaemonSetImg1" && result.err == 2 {
			t.Error("Test 5: Failed to identify that image tag is missing. Refer: fakeDaemonSetImg1")
		}

		if result.name == "fakeDaemonSetImg2" && result.err == 1 {
			t.Error("Test 6: Failed to identify wrong image tag. Refer: fakeDaemonSetImg2")
		}
	}

}

func TestPodImg(t *testing.T) {
	fakePods := fakeaudit.GetPods("fakePodImg")
	wg.Add(2)
	results := auditImages("fakeContainerImg:1.5", kubeAuditPods{list: fakePods})

	if len(results) != 2 {
		t.Error("Test 1: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakePodImg1" && result.err != 1 {
			t.Error("Test 2: Failed to identify that image tag is missing. Refer: fakePodImg1.yml")
		}

		if result.name == "fakePodImg2" && result.img_tag != "1.5" {
			t.Error("Test 3: Failed to identify the correct image tag which is present. Refer: akePodImg2.yml")
		}
	}

	results = auditImages("fakeContainerImg:1.6", kubeAuditPods{list: fakePods})

	if len(results) != 2 {
		t.Error("Test 4: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakePodImg1" && result.err == 2 {
			t.Error("Test 5: Failed to identify that image tag is missing. Refer: fakePodImg1")
		}

		if result.name == "fakePodImg2" && result.err == 1 {
			t.Error("Test 6: Failed to identify wrong image tag. Refer: fakePodImg2")
		}
	}

}

func TestReplicationControllerImg(t *testing.T) {
	fakeReplicationControllers := fakeaudit.GetReplicationControllers("fakeReplicationControllerImg")
	wg.Add(2)
	results := auditImages("fakeContainerImg:1.5", kubeAuditReplicationControllers{list: fakeReplicationControllers})

	if len(results) != 2 {
		t.Error("Test 1: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeReplicationControllerImg1" && result.err != 1 {
			t.Error("Test 2: Failed to identify that image tag is missing. Refer: fakeReplicationControllerImg1.yml")
		}

		if result.name == "fakeReplicationControllerImg2" && result.img_tag != "1.5" {
			t.Error("Test 3: Failed to identify the correct image tag which is present. Refer: fakeReplicationControllerImg2.yml")
		}
	}

	results = auditImages("fakeContainerImg:1.6", kubeAuditReplicationControllers{list: fakeReplicationControllers})

	if len(results) != 2 {
		t.Error("Test 4: Failed to identify all bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeReplicationControllerImg1" && result.err == 2 {
			t.Error("Test 5: Failed to identify that image tag is missing. Refer: fakeReplicationControllerImg1")
		}

		if result.name == "fakeReplicationControllerImg2" && result.err == 1 {
			t.Error("Test 6: Failed to identify wrong image tag. Refer: fakeReplicationControllerImg2")
		}
	}

}
