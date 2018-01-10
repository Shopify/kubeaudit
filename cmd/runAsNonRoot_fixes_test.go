package cmd

import "testing"

func TestFixRunAsNonRoot(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_false.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsNonRootNil(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_nil.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsNonRootFalse(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_false.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsRootFalseAllowed(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_false_allowed.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsNonRootMisconfiguredAllow(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_misconfigured_allow.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}
