package nonroot

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestFixRunAsNonRoot(t *testing.T) {
	cases := []struct {
		file          string
		fixtureDir    string
		expectedValue *bool
	}{
		{"run_as_non_root_nil_v1.yml", fixtureDir, k8s.NewTrue()},
		{"run_as_non_root_false_v1.yml", fixtureDir, k8s.NewTrue()},
		{"run_as_non_root_false_allowed_v1.yml", fixtureDir, k8s.NewFalse()},
		{"run_as_non_root_redundant_override_container_v1.yml", fixtureDir, k8s.NewTrue()},
		{"run_as_non_root_redundant_override_pod_v1.yml", fixtureDir, nil},
		{"run_as_non_root_psc_false_v1.yml", fixtureDir, k8s.NewTrue()},
		{"run_as_non_root_psc_true_csc_false_v1.yml", fixtureDir, k8s.NewTrue()},
		{"run_as_non_root_psc_false_csc_false_v1.yml", fixtureDir, k8s.NewTrue()},
		{"run_as_non_root_psc_false_allowed_v1.yml", fixtureDir, nil},
		{"run_as_non_root_psc_false_csc_true_v1.yml", fixtureDir, k8s.NewTrue()},
		{"run_as_non_root_psc_false_csc_nil_multiple_cont_v1.yml", fixtureDir, k8s.NewTrue()},
		{"run_as_non_root_psc_false_csc_true_multiple_cont_v1.yml", fixtureDir, k8s.NewTrue()},
		{"run_as_non_root_psc_false_allowed_multi_containers_single_label_v1.yml", fixtureDir, k8s.NewTrue()},

		// Shared fixtures
		{"security_context_nil_v1.yml", test.SharedFixturesDir, k8s.NewTrue()},
		{"security_context_nil_v1beta1.yml", test.SharedFixturesDir, k8s.NewTrue()},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, tt.fixtureDir, tt.file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				for _, container := range containers {
					if tt.expectedValue == nil {
						assert.True(t, (container.SecurityContext == nil || container.SecurityContext.RunAsNonRoot == nil))
					} else {
						assert.Equal(t, *tt.expectedValue, *container.SecurityContext.RunAsNonRoot)
					}
				}
			}
		})
	}

	file := "run_as_non_root_psc_false_allowed_multi_containers_multi_labels_v1.yml"
	t.Run(file, func(t *testing.T) {
		resources, _ := test.FixSetup(t, fixtureDir, file, New())
		for _, resource := range resources {
			containers := k8s.GetContainers(resource)
			for _, container := range containers {
				switch container.Name {
				case "fakeContainerRANR":
					assert.True(t, (container.SecurityContext == nil || container.SecurityContext.RunAsNonRoot == nil))
				case "fakeContainerRANR2":
					assert.True(t, *container.SecurityContext.RunAsNonRoot)
				}
			}
		}
	})
}
