package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestSecurityContextNIL_SC(t *testing.T) {
	runTest(t, "security_context_nil.yml", auditCapabilities, ErrorSecurityContextNIL)
}

func TestCapabilitiesNIL(t *testing.T) {
	runTest(t, "capabilities_nil.yml", auditCapabilities, ErrorCapabilitiesNIL)
}

func TestCapabilitiesAdded(t *testing.T) {
	runTest(t, "capabilities_added.yml", auditCapabilities, ErrorCapabilitiesAdded)
}

func TestCapabilitiesNoneDropped(t *testing.T) {
	runTest(t, "capabilities_none_dropped.yml", auditCapabilities, ErrorCapabilitiesNoneDropped)
}

func TestCapabilitiesSomeDropped(t *testing.T) {
	runTest(t, "capabilities_some_dropped.yml", auditCapabilities, ErrorCapabilitiesSomeDropped)
}
