package asat

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestFixAutomountServiceAccountToken(t *testing.T) {
	cases := []struct {
		file                             string
		expectedDeprecatedServiceAccount string
		expectedServiceAccountName       string
		expectedAutomountToken           *bool
	}{
		{"service_account_token_deprecated_multiple_resources_v1.yml", "", "fakeDeprecatedServiceAccount", nil},
		{"service_account_token_deprecated_v1.yml", "", "fakeDeprecatedServiceAccount", nil},
		{"service_account_token_nil_and_no_name_v1.yml", "", "", k8s.NewFalse()},
		{"service_account_token_redundant_override_v1.yml", "", "", k8s.NewFalse()},
		{"service_account_token_true_allowed_v1.yml", "", "", k8s.NewTrue()},
		{"service_account_token_true_and_default_name_v1.yml", "", "default", k8s.NewFalse()},
		{"service_account_token_true_and_no_name_multiple_resources_v1.yml", "", "", k8s.NewFalse()},
		{"service_account_token_true_and_no_name_v1.yml", "", "", k8s.NewFalse()},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, fixtureDir, tt.file, New())
			for _, resource := range resources {
				podSpec := k8s.GetPodSpec(resource)
				assert.Equal(t, tt.expectedDeprecatedServiceAccount, podSpec.DeprecatedServiceAccount)
				assert.Equal(t, tt.expectedServiceAccountName, podSpec.ServiceAccountName)
				if tt.expectedAutomountToken == nil {
					assert.Nil(t, podSpec.AutomountServiceAccountToken)
				} else {
					assert.Equal(t, *tt.expectedAutomountToken, *podSpec.AutomountServiceAccountToken)
				}

			}
		})
	}
}
