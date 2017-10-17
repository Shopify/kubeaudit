package cmd

import (
	"testing"

	"github.com/Shopify/kubeaudit/fakeaudit"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeDeploymentSC")
	fakeaudit.CreateFakeDeploymentSC("fakeDeploymentSC")
	fakeaudit.CreateFakeNamespace("fakeStatefulSetSC")
	fakeaudit.CreateFakeStatefulSetSC("fakeStatefulSetSC")
	fakeaudit.CreateFakeNamespace("fakeDaemonSetSC")
	fakeaudit.CreateFakeDaemonSetSC("fakeDaemonSetSC")
	fakeaudit.CreateFakeNamespace("fakePodSC")
	fakeaudit.CreateFakePodSC("fakePodSC")
	fakeaudit.CreateFakeNamespace("fakeReplicationControllerSC")
	fakeaudit.CreateFakeReplicationControllerSC("fakeReplicationControllerSC")
}

func TestDeploymentSC(t *testing.T) {
	fakeDeployments := fakeaudit.GetDeployments("fakeDeploymentSC")
	wg.Add(1)
	results := auditSecurityContext(kubeAuditDeployments{list: fakeDeployments})

	if len(results) != 4 {
		t.Error("Test 1: Failed to catch all the bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeDeploymentSC1" && result.err != ErrorSecurityContextNIL {
			t.Error("Test 2: Failed to recognize security context missing. Refer: fakeDeploymentSC1.yml")
		}

		if result.name == "fakeDeploymentSC2" && result.err != ErrorCapabilitiesNIL {
			t.Error("Test 3: Failed to recognize capabilities field missing. Refer: fakeDeploymentSC2.yml")
		}

		if result.name == "fakeDeploymentSC3" && (result.err != ErrorCapabilitiesAddedOrNotDropped || result.capsAdded == nil) {
			t.Error("Test 4: Failed to identify new capabilities were added. Refer: fakeDeploymentSC3.yml")
		}

		if result.name == "fakeDeploymentSC3" && (result.err != ErrorCapabilitiesAddedOrNotDropped || result.capsDropped) {
			t.Error("Test 5: Failed to identify no capabilities were droped. Refer: fakeDeploymentsSC3.yml")
		}

		if result.name == "fakeDeploymentSC4" && (result.err != ErrorCapabilitiesAddedOrNotDropped || !result.capsDropped || result.capsAdded == nil) {
			t.Error("Test 6: Failed to identify caps were added. Refer: fakeDeploymentSC4.yml")
		}
	}
}

func TestStatefulSetSC(t *testing.T) {
	fakeStatefulSets := fakeaudit.GetStatefulSets("fakeStatefulSetSC")
	wg.Add(1)
	results := auditSecurityContext(kubeAuditStatefulSets{list: fakeStatefulSets})

	if len(results) != 4 {
		t.Error("Test 1: Failed to catch all the bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeStatefulSetSC1" && result.err != ErrorSecurityContextNIL {
			t.Error("Test 2: Failed to recognize security context missing. Refer: fakeStatefulSetSC1.yml")
		}

		if result.name == "fakeStatefulSetSC2" && result.err != ErrorCapabilitiesNIL {
			t.Error("Test 3: Failed to recognize capabilities field missing. Refer: fakeStatefulSetSC2.yml")
		}

		if result.name == "fakeStatefulSetSC3" && (result.err != ErrorCapabilitiesAddedOrNotDropped || result.capsAdded == nil) {
			t.Error("Test 4: Failed to identify new capabilities were added. Refer: fakeStatefulSetSC3.yml")
		}

		if result.name == "fakeStatefulSetSC3" && (result.err != ErrorCapabilitiesAddedOrNotDropped || result.capsDropped) {
			t.Error("Test 5: Failed to identify no capabilities were droped. Refer: fakeStatefulSetSC3.yml")
		}

		if result.name == "fakeStatefulSetSC4" && (result.err != ErrorCapabilitiesAddedOrNotDropped || !result.capsDropped || result.capsAdded == nil) {
			t.Error("Test 6: Failed to identify caps were added. Refer: fakeStatefulSetSC4.yml")
		}
	}
}

func TestDaemonSetSC(t *testing.T) {
	fakeDaemonSets := fakeaudit.GetDaemonSets("fakeDaemonSetSC")
	wg.Add(1)
	results := auditSecurityContext(kubeAuditDaemonSets{list: fakeDaemonSets})

	if len(results) != 4 {
		t.Error("Test 1: Failed to catch all the bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeDaemonSetSC1" && result.err != ErrorSecurityContextNIL {
			t.Error("Test 2: Failed to recognize security context missing. Refer: fakeDaemonSetSC1.yml")
		}

		if result.name == "fakeDaemonSetSC2" && result.err != ErrorCapabilitiesNIL {
			t.Error("Test 3: Failed to recognize capabilities field missing. Refer: fakeDaemonSetSC2.yml")
		}

		if result.name == "fakeDaemonSetSC3" && (result.err != ErrorCapabilitiesAddedOrNotDropped || result.capsAdded == nil) {
			t.Error("Test 4: Failed to identify new capabilities were added. Refer: fakeDaemonSetSC3.yml")
		}

		if result.name == "fakeDaemonSetSC3" && (result.err != ErrorCapabilitiesAddedOrNotDropped || result.capsDropped) {
			t.Error("Test 5: Failed to identify no capabilities were droped. Refer: fakeDaemonSetSC3.yml")
		}

		if result.name == "fakeDaemonSetSC4" && (result.err != ErrorCapabilitiesAddedOrNotDropped || !result.capsDropped || result.capsAdded == nil) {
			t.Error("Test 6: Failed to identify caps were added. Refer: fakeDaemonSetSC4.yml")
		}
	}
}

func TestPodSC(t *testing.T) {
	fakePods := fakeaudit.GetPods("fakePodSC")
	wg.Add(1)
	results := auditSecurityContext(kubeAuditPods{list: fakePods})

	if len(results) != 4 {
		t.Error("Test 1: Failed to catch all the bad configurations")
	}

	for _, result := range results {
		if result.name == "fakePodSC1" && result.err != ErrorSecurityContextNIL {
			t.Error("Test 2: Failed to recognize security context missing. Refer: fakePodSC1.yml")
		}

		if result.name == "fakePodSC2" && result.err != ErrorCapabilitiesNIL {
			t.Error("Test 3: Failed to recognize capabilities field missing. Refer: fakePodSC2.yml")
		}

		if result.name == "fakePodSC3" && (result.err != ErrorCapabilitiesAddedOrNotDropped || result.capsAdded == nil) {
			t.Error("Test 4: Failed to identify new capabilities were added. Refer: fakePodSC3.yml")
		}

		if result.name == "fakePodSC3" && (result.err != ErrorCapabilitiesAddedOrNotDropped || result.capsDropped) {
			t.Error("Test 5: Failed to identify no capabilities were droped. Refer: fakePodSC3.yml")
		}

		if result.name == "fakePodSC4" && (result.err != ErrorCapabilitiesAddedOrNotDropped || !result.capsDropped || result.capsAdded == nil) {
			t.Error("Test 6: Failed to identify caps were added. Refer: fakePodSC4.yml")
		}
	}
}

func TestReplicationControllerSC(t *testing.T) {
	fakeReplicationControllers := fakeaudit.GetReplicationControllers("fakeReplicationControllerSC")
	wg.Add(1)
	results := auditSecurityContext(kubeAuditReplicationControllers{list: fakeReplicationControllers})

	if len(results) != 4 {
		t.Error("Test 1: Failed to catch all the bad configurations")
	}

	for _, result := range results {
		if result.name == "fakeReplicationControllerSC1" && result.err != ErrorSecurityContextNIL {
			t.Error("Test 2: Failed to recognize security context missing. Refer: fakeReplicationControllerSC1.yml")
		}

		if result.name == "fakeReplicationControllerSC2" && result.err != ErrorCapabilitiesNIL {
			t.Error("Test 3: Failed to recognize capabilities field missing. Refer: fakeReplicationControllerSC2.yml")
		}

		if result.name == "fakeReplicationControllerSC3" && (result.err != ErrorCapabilitiesAddedOrNotDropped || result.capsAdded == nil) {
			t.Error("Test 4: Failed to identify new capabilities were added. Refer: fakeReplicationControllerSC3.yml")
		}

		if result.name == "fakeReplicationControllerSC3" && (result.err != ErrorCapabilitiesAddedOrNotDropped || result.capsDropped) {
			t.Error("Test 5: Failed to identify no capabilities were droped. Refer: fakeReplicationControllerSC3.yml")
		}

		if result.name == "fakeReplicationControllerSC4" && (result.err != ErrorCapabilitiesAddedOrNotDropped || !result.capsDropped || result.capsAdded == nil) {
			t.Error("Test 6: Failed to identify caps were added. Refer: fakeReplicationControllerSC4.yml")
		}
	}
}
