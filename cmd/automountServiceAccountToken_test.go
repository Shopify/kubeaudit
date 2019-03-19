package cmd

import "testing"

func TestServiceAccountTokenDeprecatedV1(t *testing.T) {
	runAuditTest(t, "service_account_token_deprecated_v1.yml", auditAutomountServiceAccountToken, []int{ErrorServiceAccountTokenDeprecated})
}

func TestServiceAccountTokenTrueAndNoNameV1(t *testing.T) {
	runAuditTest(t, "service_account_token_true_and_no_name_v1.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenTrueAndNoName})
}

func TestServiceAccountTokenNilAndNoNameV1(t *testing.T) {
	runAuditTest(t, "service_account_token_nil_and_no_name_v1.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenNilAndNoName})
}

func TestServiceAccountTokenTrueAllowedV1(t *testing.T) {
	runAuditTest(t, "service_account_token_true_allowed_v1.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenTrueAllowed})
}

func TestServiceAccountTokenMisconfiguredAllowV1(t *testing.T) {
	runAuditTest(t, "service_account_token_misconfigured_allow_v1.yml", auditAutomountServiceAccountToken, []int{ErrorMisconfiguredKubeauditAllow})
}

func TestServiceAccountTokenTrueAndDefaultNameV1(t *testing.T) {
	runAuditTest(t, "service_account_token_true_and_default_name_v1.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenTrueAndNoName})
}

func TestAutomountServiceAccountTokenFromConfig(t *testing.T) {
	rootConfig.kubeauditConfig = "../configs/allow_Automount_Service_Account_Token_From_Config.yml"
	runAuditTest(t, "service_account_token_deprecated_v1.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenTrueAllowed})
	runAuditTest(t, "service_account_token_true_and_no_name_v1.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenTrueAllowed})
	runAuditTest(t, "service_account_token_nil_and_no_name_v1.yml", auditAutomountServiceAccountToken, []int{ErrorAutomountServiceAccountTokenTrueAllowed})
}
