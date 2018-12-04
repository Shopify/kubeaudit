package cmd

import "testing"

func TestFixRunAsNonRootV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_false_v1.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsNonRootNilV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_nil_v1.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsNonRootFalseV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_false_v1.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsRootFalseAllowedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_false_allowed_v1.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsNonRootMisconfiguredAllowV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_misconfigured_allow_v1.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}
