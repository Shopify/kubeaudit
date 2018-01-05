package cmd

import "testing"

func TestServiceAccountTokenDeprecated(t *testing.T) {
	runAuditTest(t, "service_account_token_deprecated.yml", auditAutomountServiceAccountToken, []int{ErrorServiceAccountTokenDeprecated})
}

func TestServiceAccountTokenTrueAndNoName(t *testing.T) {
	runAuditTest(t, "service_account_token_true_and_no_name.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenTrueAndNoName})
}

func TestServiceAccountTokenNILAndNoName(t *testing.T) {
	runAuditTest(t, "service_account_token_nil_and_no_name.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenNILAndNoName})
}

func TestServiceAccountTokenTrueAllowed(t *testing.T) {
	runAuditTest(t, "service_account_token_true_allowed.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenTrueAllowed})
}

func TestServiceAccountTokenMisconfiguredAllow(t *testing.T) {
	runAuditTest(t, "service_account_token_misconfigured_allow.yml", auditAutomountServiceAccountToken, []int{ErrorMisconfiguredKubeauditAllow})
}
