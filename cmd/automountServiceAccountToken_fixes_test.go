package cmd

import "testing"

func TestFixServiceAccountToken(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_true_and_no_name.yml", auditAutomountServiceAccountToken, fixServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationController:
		assert.Equal(false, *typ.Spec.Template.Spec.AutomountServiceAccountToken)
	}
}

func TestFixDeprecatedServiceAccount(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_deprecated.yml", auditAutomountServiceAccountToken, fixDeprecatedServiceAccount)
	switch typ := resource.(type) {
	case *ReplicationController:
		assert.Equal("", typ.Spec.Template.Spec.DeprecatedServiceAccount)
	}
}
