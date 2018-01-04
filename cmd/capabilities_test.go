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
	runAuditTest(t, "security_context_nil.yml", auditCapabilities, []int{ErrorSecurityContextNIL})
}

func TestCapabilitiesNIL(t *testing.T) {
	runAuditTest(t, "capabilities_nil.yml", auditCapabilities, []int{ErrorCapabilitiesNIL})
}

func TestCapabilitiesAdded(t *testing.T) {
	runAuditTest(t, "capabilities_added.yml", auditCapabilities, []int{ErrorCapabilityAdded})
}

func TestCapabilitiesSomeAllowed(t *testing.T) {
	runAuditTest(t, "capabilities_some_allowed.yml", auditCapabilities, []int{ErrorCapabilityAllowed})
}

func TestCapabilitiesSomeDropped(t *testing.T) {
	runAuditTest(t, "capabilities_some_dropped.yml", auditCapabilities, []int{ErrorCapabilityNotDropped})
}

func TestCapabilitiesMisconfiguredAllow(t *testing.T) {
	runAuditTest(t, "capabilities_misconfigured_allow.yml", auditCapabilities, []int{ErrorMisconfiguredKubeauditAllow})
}
