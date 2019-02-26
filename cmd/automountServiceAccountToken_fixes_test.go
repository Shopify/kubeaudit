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

func TestFixServiceAccountTokenDeprecatedV2(t *testing.T) {
	assert, resources := FixTestSetupMultipleResources(t, "service_account_token_deprecated_multiple_resources_v1.yml", auditAutomountServiceAccountToken)
	for index := range resources {
		switch t := resources[index].(type) {
		case *CronJobV1Beta1:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.JobTemplate.Spec.Template.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.JobTemplate.Spec.Template.Spec.DeprecatedServiceAccount)
		case *DaemonSetV1:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.Template.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.Template.Spec.DeprecatedServiceAccount)
		case *DaemonSetV1Beta1:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.Template.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.Template.Spec.DeprecatedServiceAccount)
		case *DeploymentExtensionsV1Beta1:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.Template.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.Template.Spec.DeprecatedServiceAccount)
		case *DeploymentV1:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.Template.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.Template.Spec.DeprecatedServiceAccount)
		case *DeploymentV1Beta1:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.Template.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.Template.Spec.DeprecatedServiceAccount)
		case *DeploymentV1Beta2:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.Template.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.Template.Spec.DeprecatedServiceAccount)
		case *PodV1:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.DeprecatedServiceAccount)
		case *ReplicationControllerV1:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.Template.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.Template.Spec.DeprecatedServiceAccount)
		case *StatefulSetV1:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.Template.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.Template.Spec.DeprecatedServiceAccount)
		case *StatefulSetV1Beta1:
			assert.Equal("fakeDeprecatedServiceAccount", t.Spec.Template.Spec.ServiceAccountName)
			assert.Equal("", t.Spec.Template.Spec.DeprecatedServiceAccount)
		}
	}
}

func TestFixServiceAccountTokenTrueAndNoNameV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "service_account_token_true_and_no_name_v1.yml", auditAutomountServiceAccountToken)
	switch typ := resource.(type) {
	case *ReplicationControllerV1:
		assert.False(*typ.Spec.Template.Spec.AutomountServiceAccountToken)
	}
}

func TestFixServiceAccountTokenTrueAndNoNameV2(t *testing.T) {
	assert, resources := FixTestSetupMultipleResources(t, "service_account_token_true_and_no_name_multiple_resources_v1.yml", auditAutomountServiceAccountToken)
	for index := range resources {
		switch t := resources[index].(type) {
		case *CronJobV1Beta1:
			assert.False(*t.Spec.JobTemplate.Spec.Template.Spec.AutomountServiceAccountToken)
		case *DaemonSetV1:
			assert.False(*t.Spec.Template.Spec.AutomountServiceAccountToken)
		case *DaemonSetV1Beta1:
			assert.False(*t.Spec.Template.Spec.AutomountServiceAccountToken)
		case *DeploymentExtensionsV1Beta1:
			assert.False(*t.Spec.Template.Spec.AutomountServiceAccountToken)
		case *DeploymentV1:
			assert.False(*t.Spec.Template.Spec.AutomountServiceAccountToken)
		case *DeploymentV1Beta1:
			assert.False(*t.Spec.Template.Spec.AutomountServiceAccountToken)
		case *DeploymentV1Beta2:
			assert.False(*t.Spec.Template.Spec.AutomountServiceAccountToken)
		case *PodV1:
			assert.False(*t.Spec.AutomountServiceAccountToken)
		case *ReplicationControllerV1:
			assert.False(*t.Spec.Template.Spec.AutomountServiceAccountToken)
		case *StatefulSetV1:
			assert.False(*t.Spec.Template.Spec.AutomountServiceAccountToken)
		case *StatefulSetV1Beta1:
			assert.False(*t.Spec.Template.Spec.AutomountServiceAccountToken)
		}
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
