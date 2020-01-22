package privesc

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestFixPrivilegeEscalation(t *testing.T) {
	cases := []struct {
		file          string
		fixtureDir    string
		expectedValue bool
	}{
		{"allow_privilege_escalation_nil_v1.yml", fixtureDir, false},
		{"allow_privilege_escalation_nil_v1beta1.yml", fixtureDir, false},
		{"allow_privilege_escalation_redundant_override_v1.yml", fixtureDir, false},
		{"allow_privilege_escalation_redundant_override_v1beta1.yml", fixtureDir, false},
		{"allow_privilege_escalation_true_allowed_v1.yml", fixtureDir, true},
		{"allow_privilege_escalation_true_allowed_v1beta1.yml", fixtureDir, true},
		{"allow_privilege_escalation_true_multiple_allowed_multiple_containers_v1beta.yml", fixtureDir, true},
		{"allow_privilege_escalation_true_v1.yml", fixtureDir, false},
		{"allow_privilege_escalation_true_v1beta1.yml", fixtureDir, false},

		// Shared fixtures
		{"security_context_nil_v1.yml", test.SharedFixturesDir, false},
		{"security_context_nil_v1beta1.yml", test.SharedFixturesDir, false},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, tt.fixtureDir, tt.file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				for _, container := range containers {
					assert.Equal(t, tt.expectedValue, *container.SecurityContext.AllowPrivilegeEscalation)
				}
			}
		})
	}

	file := "allow_privilege_escalation_true_single_allowed_multiple_containers_v1beta.yml"
	t.Run(file, func(t *testing.T) {
		resources, _ := test.FixSetup(t, fixtureDir, file, New())
		for _, resource := range resources {
			containers := k8s.GetContainers(resource)
			for _, container := range containers {
				switch container.Name {
				case "fakeContainerAPE":
					assert.False(t, *container.SecurityContext.AllowPrivilegeEscalation)
				case "fakeContainerAPE2":
					assert.True(t, *container.SecurityContext.AllowPrivilegeEscalation)
				}
			}
		}
	})
}
