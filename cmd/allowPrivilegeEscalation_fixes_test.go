package cmd

import "testing"

func TestFixAllowPrivilegeEscalationV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "allow_privilege_escalation_nil_v1.yml", auditAllowPrivilegeEscalation)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.AllowPrivilegeEscalation)
	}
}

func TestFixAllowPrivilegeEscalationTrueV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "allow_privilege_escalation_true_v1.yml", auditAllowPrivilegeEscalation)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.AllowPrivilegeEscalation)
	}
}

func TestFixAllowPrivilegeEscalationTrueAllowedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "allow_privilege_escalation_true_allowed_v1.yml", auditAllowPrivilegeEscalation)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.AllowPrivilegeEscalation)
	}
}

func TestFixAllowPrivilegeEscalationMisconfiguredAllowV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "allow_privilege_escalation_misconfigured_allow_v1.yml", auditAllowPrivilegeEscalation)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.AllowPrivilegeEscalation)
	}
}
