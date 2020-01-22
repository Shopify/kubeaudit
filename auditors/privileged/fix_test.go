package privileged

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestFixPrivileged(t *testing.T) {
	cases := []struct {
		file          string
		fixtureDir    string
		expectedValue bool
	}{
		{"privileged_nil_v1.yml", fixtureDir, false},
		{"privileged_true_v1.yml", fixtureDir, false},
		{"privileged_true_allowed_v1.yml", fixtureDir, true},
		{"privileged_redundant_override_v1.yml", fixtureDir, false},
		{"privileged_true_allowed_multi_containers_multi_labels_v1.yml", fixtureDir, true},

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
					assert.Equal(t, tt.expectedValue, *container.SecurityContext.Privileged)
				}
			}
		})
	}

	file := "privileged_true_allowed_multi_containers_single_label_v1.yml"
	t.Run(file, func(t *testing.T) {
		resources, _ := test.FixSetup(t, fixtureDir, file, New())
		for _, resource := range resources {
			containers := k8s.GetContainers(resource)
			for _, container := range containers {
				switch container.Name {
				case "fakeContainerPrivileged":
					assert.False(t, *container.SecurityContext.Privileged)
				case "fakeContainerPrivileged2":
					assert.True(t, *container.SecurityContext.Privileged)
				}
			}
		}
	})
}
