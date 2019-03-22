package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecommendedCapabilitiesToBeDroppedV1(t *testing.T) {
	assert := assert.New(t)
	capabilities, err := recommendedCapabilitiesToBeDropped()
	assert.Nil(err)
	assert.Equal(NewCapSetFromArray([]CapabilityV1{"AUDIT_WRITE", "CHOWN", "DAC_OVERRIDE", "FOWNER", "FSETID", "KILL", "MKNOD", "NET_BIND_SERVICE", "NET_RAW", "SETFCAP", "SETGID", "SETUID", "SETPCAP", "SYS_CHROOT"}), capabilities, "")
}

func TestSecurityContextNil_SCV1(t *testing.T) {
	runAuditTest(t, "security_context_nil_v1.yml", auditCapabilities, []int{ErrorCapabilityNotDropped})
}

func TestCapabilitiesNilV1Beta2(t *testing.T) {
	runAuditTest(t, "capabilities_nil_v1beta2.yml", auditCapabilities, []int{ErrorCapabilityNotDropped})
}

func TestCapabilitiesAddedV1Beta2(t *testing.T) {
	runAuditTest(t, "capabilities_added_v1beta2.yml", auditCapabilities, []int{ErrorCapabilityAdded})
}

func TestCapabilitiesSomeAllowedV1Beta2(t *testing.T) {
	runAuditTest(t, "capabilities_some_allowed_v1beta2.yml", auditCapabilities, []int{ErrorCapabilityAllowed, ErrorCapabilityAllowed})
}

func TestCapabilitiesSomeDroppedV1Beta2(t *testing.T) {
	runAuditTest(t, "capabilities_some_dropped_v1beta2.yml", auditCapabilities, []int{ErrorCapabilityNotDropped})
}

func TestCapabilitiesMisconfiguredAllowV1Beta2(t *testing.T) {
	runAuditTest(t, "capabilities_misconfigured_allow_v1beta2.yml", auditCapabilities, []int{ErrorMisconfiguredKubeauditAllow})
}

func TestCapabilitiesDroppedAllV1Beta2(t *testing.T) {
	runAuditTest(t, "capabilities_dropped_all_v1beta2.yml", auditCapabilities, []int{})
}

func TestCapabilitiesSomeAllowedMultiContainersAllLabelsV1Beta2(t *testing.T) {
	runAuditTest(t, "capabilities_some_allowed_multi_containers_all_container_labels_v1beta2.yml", auditCapabilities, []int{ErrorCapabilityAllowed, ErrorCapabilityAllowed})
}

func TestCapabilitiesSomeAllowedMultiContainersSomeLabelsV1Beta2(t *testing.T) {
	runAuditTest(t, "capabilities_some_allowed_multi_containers_some_container_labels_v1beta2.yml", auditCapabilities, []int{ErrorCapabilityAdded, ErrorCapabilityAdded, ErrorCapabilityAllowed})
}

func TestCapabilitiesSomeAllowedMultiContainersMixLabelsV1Beta2(t *testing.T) {
	runAuditTest(t, "capabilities_some_allowed_multi_containers_mix_labels_v1beta2.yml", auditCapabilities, []int{ErrorCapabilityAllowed, ErrorCapabilityAllowed})
}

func TestCapabilitiesManualConfigV2(t *testing.T) {
	rootConfig.auditConfig = "../fake/file/path"
	runAuditTest(t, "capabilities_some_dropped_v1beta2.yml", auditCapabilities, []int{KubeauditInternalError})
}

func TestCustomCapabilitiesToBeDroppedV1(t *testing.T) {
	assert := assert.New(t)
	rootConfig.auditConfig = "../configs/custom_capabilities_to_be_dropped_v1.yml"
	capabilities, err := recommendedCapabilitiesToBeDropped()
	assert.Nil(err)
	assert.Equal(NewCapSetFromArray([]CapabilityV1{"MKNOD", "CHOWN", "DAC_OVERRIDE", "FSETID", "SETGID", "NET_BIND_SERVICE", "SETFCAP"}), capabilities, "")
}

func TestCustomCapabilitiesToBeDroppedV2(t *testing.T) {
	assert := assert.New(t)
	rootConfig.auditConfig = "../configs/custom_capabilities_to_be_dropped_v1.yml"
	capabilities, err := recommendedCapabilitiesToBeDropped()
	assert.Nil(err)
	assert.NotEqual(NewCapSetFromArray([]CapabilityV1{"MKNOD", "SYS_CHROOT", "KILL", "CHOWN", "DAC_OVERRIDE", "FSETID", "SETGID", "NET_BIND_SERVICE", "SETFCAP"}), capabilities, "")
}
