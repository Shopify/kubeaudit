package cmd

import "testing"

func TestFixPrivilegeEscalation(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_nil.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.Privileged)
	}
}

func TestFixPrivilegedNil(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_nil.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.Privileged)
	}
}

func TestFixPrivilegedTrue(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_true.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.Privileged)
	}
}

func TestFixPrivilegedTrueAllowed(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_true_allowed.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.Privileged)
	}
}

func TestFixPrivilegedMisconfiguredAllow(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_misconfigured_allow.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.Privileged)
	}
}
