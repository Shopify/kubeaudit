package asat

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
)

func TestFixAutomountServiceAccountToken(t *testing.T) {
	cases := []struct {
		file                             string
		expectedDeprecatedServiceAccount string
		expectedServiceAccountName       string
		expectedAutomountToken           *bool
	}{
		{"service-account-token-deprecated.yml", "", "deprecated", nil},
		{"service-account-token-nil-and-no-name.yml", "", "", k8s.NewFalse()},
		{"service-account-token-redundant-override.yml", "", "", k8s.NewFalse()},
		{"service-account-token-true-allowed.yml", "", "", k8s.NewTrue()},
		{"service-account-token-true-and-default-name.yml", "", "default", k8s.NewFalse()},
		{"service-account-token-true-and-no-name.yml", "", "", k8s.NewFalse()},
		{"service-account-token-false.yml", "", "", k8s.NewFalse()},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, fixtureDir, tc.file, New())
			for _, resource := range resources {
				podSpec := k8s.GetPodSpec(resource)
				assert.Equal(t, tc.expectedDeprecatedServiceAccount, podSpec.DeprecatedServiceAccount)
				assert.Equal(t, tc.expectedServiceAccountName, podSpec.ServiceAccountName)
				if tc.expectedAutomountToken == nil {
					assert.Nil(t, podSpec.AutomountServiceAccountToken)
				} else {
					assert.Equal(t, *tc.expectedAutomountToken, *podSpec.AutomountServiceAccountToken)
				}

			}
		})
	}

	// Test that if a default ServiceAccount was found, its 'automountServiceAccountToken' is set to false
	// instead of on the PodSpec
	files := []string{
		"service-account-token-nil-and-no-name-and-default-sa.yml",
		"service-account-token-true-and-default-sa.yml",
	}
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, fixtureDir, file, New())
			for _, resource := range resources {
				if serviceAccount, ok := resource.(*k8s.ServiceAccountV1); ok {
					assert.Equal(t, serviceAccount.AutomountServiceAccountToken, k8s.NewFalse())
					continue
				}

				podSpec := k8s.GetPodSpec(resource)
				assert.Equal(t, "", podSpec.DeprecatedServiceAccount)
				assert.Equal(t, "", podSpec.ServiceAccountName)
				assert.Nil(t, podSpec.AutomountServiceAccountToken)
			}
		})
	}
}
