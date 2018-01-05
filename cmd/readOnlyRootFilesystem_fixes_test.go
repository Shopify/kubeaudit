package cmd

import "testing"

func TestFixReadOnlyRootFilesystem(t *testing.T) {
	assert, resource := FixTestSetup(t, "read_only_root_filesystem_false.yml", auditReadOnlyRootFS, fixReadOnlyRootFilesystem)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.ReadOnlyRootFilesystem)
	}
}
