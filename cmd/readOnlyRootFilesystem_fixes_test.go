package cmd

import "testing"

func TestFixReadOnlyRootFilesystemNilV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "read_only_root_filesystem_nil_v1.yml", auditReadOnlyRootFS)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.ReadOnlyRootFilesystem)
	}
}

func TestFixReadOnlyRootFilesystemFalseV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "read_only_root_filesystem_false_v1.yml", auditReadOnlyRootFS)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.ReadOnlyRootFilesystem)
	}
}

func TestFixReadOnlyRootFilesystemFalseAllowedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "read_only_root_filesystem_false_allowed_v1.yml", auditReadOnlyRootFS)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.ReadOnlyRootFilesystem)
	}
}

func TestFixReadOnlyRootFilesystemMisconfiguredAllowV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "read_only_root_filesystem_misconfigured_allow_v1.yml", auditReadOnlyRootFS)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.ReadOnlyRootFilesystem)
	}
}
