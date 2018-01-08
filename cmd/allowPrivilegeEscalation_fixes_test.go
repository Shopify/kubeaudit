package cmd

import "testing"

func TestFixAllowPrivilegeEscalation(t *testing.T) {
	assert, resource := FixTestSetup(t, "allow_privilege_escalation_nil.yml", fixAllowPrivilegeEscalation)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.AllowPrivilegeEscalation)
	}
}
