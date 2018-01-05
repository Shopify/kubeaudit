package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecommendedCapabilitiesToBeDropped(t *testing.T) {
	assert := assert.New(t)
	capabilities, err := recommendedCapabilitiesToBeDropped()
	assert.Nil(err)
	assert.Equal(NewCapSetFromArray([]Capability{"AUDIT_WRITE", "CHOWN", "DAC_OVERRIDE", "FOWNER", "FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE", "NET_RAW", "SETFCAP", "SETGID", "SETUID", "SETPCAP", "SYS_CHROOT"}), capabilities, "")
}

func TestSecurityContextNIL_SC(t *testing.T) {
	runTest(t, "security_context_nil.yml", auditCapabilities, ErrorCapabilityNotDropped)
}

func TestCapabilitiesNIL(t *testing.T) {
	runTest(t, "capabilities_nil.yml", auditCapabilities, ErrorCapabilityNotDropped)
}

func TestCapabilitiesAdded(t *testing.T) {
	runTest(t, "capabilities_added.yml", auditCapabilities, ErrorCapabilityAdded)
}

func TestCapabilitiesSomeAllowed(t *testing.T) {
	runTest(t, "capabilities_some_allowed.yml", auditCapabilities, ErrorCapabilityAllowed)
}

func TestCapabilitiesSomeDropped(t *testing.T) {
	runTest(t, "capabilities_some_dropped.yml", auditCapabilities, ErrorCapabilityNotDropped)
}

func TestCapabilitiesMisconfiguredAllow(t *testing.T) {
	runTest(t, "capabilities_misconfigured_allow.yml", auditCapabilities, ErrorMisconfiguredKubeauditAllow)
}
