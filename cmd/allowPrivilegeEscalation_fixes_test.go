package cmd

import "testing"

func TestFixAllowPrivilegeEscalation(t *testing.T) {
	assert, resource := FixTestSetup(t, "allow_privilege_escalation_nil.yml", auditAllowPrivilegeEscalation)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.AllowPrivilegeEscalation)
	}
}

func TestFixAllowPrivilegeEscalationTrue(t *testing.T) {
	assert, resource := FixTestSetup(t, "allow_privilege_escalation_true.yml", auditAllowPrivilegeEscalation)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.AllowPrivilegeEscalation)
	}
}

func TestFixAllowPrivilegeEscalationTrueAllowed(t *testing.T) {
	assert, resource := FixTestSetup(t, "allow_privilege_escalation_true_allowed.yml", auditAllowPrivilegeEscalation)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.AllowPrivilegeEscalation)
	}
}

func TestFixAllowPrivilegeEscalationMisconfiguredAllow(t *testing.T) {
	assert, resource := FixTestSetup(t, "allow_privilege_escalation_misconfigured_allow.yml", auditAllowPrivilegeEscalation)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.AllowPrivilegeEscalation)
	}
}
