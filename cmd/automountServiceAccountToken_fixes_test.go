package cmd

import "testing"

func TestFixServiceAccountTokenDeprecated(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_deprecated.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationController:
		assert.Equal("", typ.Spec.Template.Spec.DeprecatedServiceAccount)
	}
}

func TestFixServiceAccountTokenTrueAndNoName(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_true_and_no_name.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationController:
		assert.False(*typ.Spec.Template.Spec.AutomountServiceAccountToken)
	}
}

func TestFixServiceAccountTokenNILAndNoName(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_nil_and_no_name.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationController:
		assert.False(*typ.Spec.Template.Spec.AutomountServiceAccountToken)
	}
}

func TestFixServiceAccountTokenTrueAllowed(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_true_allowed.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationController:
		assert.True(*typ.Spec.Template.Spec.AutomountServiceAccountToken)
	}
}

func TestFixServiceAccountTokenMisconfiguredAllow(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_misconfigured_allow.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationController:
		assert.False(*typ.Spec.Template.Spec.AutomountServiceAccountToken)
	}
}
