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
		{"service-account-token-deprecated.yml", "", "deprecated", nil},
		{"service-account-token-nil-and-no-name.yml", "", "", k8s.NewFalse()},
		{"service-account-token-redundant-override.yml", "", "", k8s.NewFalse()},
		{"service-account-token-true-allowed.yml", "", "", k8s.NewTrue()},
		{"service-account-token-true-and-default-name.yml", "", "default", k8s.NewFalse()},
		{"service-account-token-true-and-no-name.yml", "", "", k8s.NewFalse()},
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
}
