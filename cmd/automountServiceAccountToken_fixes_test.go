package cmd

import "testing"

func TestFixServiceAccountTokenDeprecatedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_deprecated_v1.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationControllerV1:
		assert.Equal("fakeDeprecatedServiceAccount", typ.Spec.Template.Spec.ServiceAccountName)
		assert.Equal("", typ.Spec.Template.Spec.DeprecatedServiceAccount)
	}
}

func TestFixServiceAccountTokenTrueAndNoNameV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_true_and_no_name_v1.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationControllerV1:
		assert.False(*typ.Spec.Template.Spec.AutomountServiceAccountToken)
	}
}

func TestFixServiceAccountTokenNilAndNoNameV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_nil_and_no_name_v1.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationControllerV1:
		assert.False(*typ.Spec.Template.Spec.AutomountServiceAccountToken)
	}
}

func TestFixServiceAccountTokenTrueAllowedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_true_allowed_v1.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationControllerV1:
		assert.True(*typ.Spec.Template.Spec.AutomountServiceAccountToken)
	}
}

func TestFixServiceAccountTokenMisconfiguredAllowV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_misconfigured_allow_v1.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationControllerV1:
		assert.False(*typ.Spec.Template.Spec.AutomountServiceAccountToken)
	}
}
