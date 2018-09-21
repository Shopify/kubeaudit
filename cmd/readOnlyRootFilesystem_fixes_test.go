package cmd

import "testing"

func TestFixReadOnlyRootFilesystemNil(t *testing.T) {
	assert, resource := FixTestSetup(t, "read_only_root_filesystem_nil.yml", auditReadOnlyRootFS)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.ReadOnlyRootFilesystem)
	}
}

func TestFixReadOnlyRootFilesystemFalse(t *testing.T) {
	assert, resource := FixTestSetup(t, "read_only_root_filesystem_false.yml", auditReadOnlyRootFS)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.ReadOnlyRootFilesystem)
	}
}

func TestFixReadOnlyRootFilesystemFalseAllowed(t *testing.T) {
	assert, resource := FixTestSetup(t, "read_only_root_filesystem_false_allowed.yml", auditReadOnlyRootFS)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.ReadOnlyRootFilesystem)
	}
}

func TestFixReadOnlyRootFilesystemMisconfiguredAllow(t *testing.T) {
	assert, resource := FixTestSetup(t, "read_only_root_filesystem_misconfigured_allow.yml", auditReadOnlyRootFS)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.ReadOnlyRootFilesystem)
	}
}
