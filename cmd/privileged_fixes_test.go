package cmd

import "testing"

func TestFixPrivilegeEscalationV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_nil_v1.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.Privileged)
	}
}

func TestFixPrivilegedNilV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_nil_v1.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.Privileged)
	}
}

func TestFixPrivilegedTrueV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_true_v1.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.Privileged)
	}
}

func TestFixPrivilegedTrueAllowedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_true_allowed_v1.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.Privileged)
	}
}

func TestFixPrivilegedMisconfiguredAllowV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_misconfigured_allow_v1.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.Privileged)
	}
}
