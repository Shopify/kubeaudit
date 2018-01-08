package cmd

import "testing"

func TestFixPrivilegeEscalation(t *testing.T) {
	assert, resource := FixTestSetup(t, "privileged_nil.yml", auditPrivileged)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.Privileged)
	}
}
