package cmd

import "testing"

func TestFixRunAsNonRoot(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_false.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}
