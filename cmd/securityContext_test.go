package cmd

import (
	"testing"

	"github.com/Shopify/kubeaudit/fakeaudit"
	"github.com/stretchr/testify/assert"
)

func init() {
	fakeaudit.CreateFakeNamespace("fakeDeploymentSC")
	fakeaudit.CreateFakeDeploymentSC("fakeDeploymentSC")
}

func TestRecommendedCapabilitiesToBeDropped(t *testing.T) {
	assert := assert.New(t)
	capabilities, err := recommendedCapabilitiesToBeDropped()
	assert.Nil(err)
	assert.Equal([]Capability{"AUDIT_WRITE", "CHOWN", "DAC_OVERRIDE", "FOWNER", "FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE", "NET_RAW", "SETFCAP", "SETGID", "SETUID", "SETPCAP", "SYS_CHROOT"}, capabilities, "")
}

func TestCapsNotDropped(t *testing.T) {
	assert := assert.New(t)
	caps := []Capability{"CHOWN", "DAC_OVERRIDE", "FOWNER", "FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE", "NET_RAW", "SETFCAP", "SETGID", "SETUID", "SETPCAP", "SYS_CHROOT"}
	notDropped, err := capsNotDropped(caps)
	assert.Nil(err)
	assert.Equal([]Capability{"AUDIT_WRITE"}, notDropped, "")
}

func TestDeploymentSC(t *testing.T) {
	fakeDeployments := fakeaudit.GetDeployments("fakeDeploymentSC")
	results := auditSecurityContext(kubeAuditDeployments{list: fakeDeployments})

	if len(results) != 4 {
		t.Error("Test 1: Failed to catch all the bad configurations")
	}

	for _, result := range results {
		if result.Name == "fakeDeploymentSC1" && result.Occurrences[0].id != ErrorSecurityContextNIL {
			t.Error("Test 2: Failed to recognize security context missing. Refer: fakeDeploymentSC1.yml")
		}

		if result.Name == "fakeDeploymentSC2" && result.Occurrences[0].id != ErrorCapabilitiesNIL {
			t.Error("Test 3: Failed to recognize capabilities field missing. Refer: fakeDeploymentSC2.yml")
		}

		if result.Name == "fakeDeploymentSC3" && (result.Occurrences[0].id != ErrorCapabilitiesAdded) {
			t.Error("Test 4: Failed to identify new capabilities were added. Refer: fakeDeploymentSC3.yml")
		}

		if result.Name == "fakeDeploymentSC3" && (result.Occurrences[1].id != ErrorCapabilitiesNoneDropped) {
			t.Error("Test 5: Failed to identify no capabilities were droped. Refer: fakeDeploymentsSC3.yml")
		}

		if result.Name == "fakeDeploymentSC4" && (result.Occurrences[0].id != ErrorCapabilitiesAdded) {
			t.Error("Test 6: Failed to identify caps were added. Refer: fakeDeploymentSC4.yml")
		}
	}
}
