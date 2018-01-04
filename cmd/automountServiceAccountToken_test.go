package cmd

import "testing"

func TestServiceAccountTokenDeprecated(t *testing.T) {
	runTest(t, "service_account_token_deprecated.yml", auditAutomountServiceAccountToken, ErrorServiceAccountTokenDeprecated)
}

func TestServiceAccountTokenTrueAndNoName(t *testing.T) {
	runTest(t, "service_account_token_true_and_no_name.yml", auditAutomountServiceAccountToken, ErrorAutomountServiceAccountTokenTrueAndNoName)
}

func TestServiceAccountTokenNILAndNoName(t *testing.T) {
	runTest(t, "service_account_token_nil_and_no_name.yml", auditAutomountServiceAccountToken, ErrorAutomountServiceAccountTokenNILAndNoName)
}

func TestServiceAccountTokenTrueAllowed(t *testing.T) {
	runTest(t, "service_account_token_true_allowed.yml", auditAutomountServiceAccountToken, ErrorAutomountServiceAccountTokenTrueAllowed)
}

func TestServiceAccountTokenMisconfiguredAllow(t *testing.T) {
	runTest(t, "service_account_token_misconfigured_allow.yml", auditAutomountServiceAccountToken, ErrorMisconfiguredKubeauditAllow)
}
