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
		{"run-as-non-root-nil.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-non-root-false.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-non-root-false-allowed.yml", fixtureDir, k8s.NewFalse()},
		{"run-as-non-root-redundant-override-container.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-non-root-redundant-override-pod.yml", fixtureDir, nil},
		{"run-as-non-root-psc-false.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-non-root-psc-true-csc-false.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-non-root-psc-false-csc-false.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-non-root-psc-false-allowed.yml", fixtureDir, nil},
		{"run-as-non-root-psc-false-csc-true.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-non-root-psc-false-csc-nil-multiple-cont.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-non-root-psc-false-csc-true-multiple-cont.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-non-root-psc-false-allowed-multi-containers-single-label.yml", fixtureDir, k8s.NewTrue()},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, tc.fixtureDir, tc.file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				for _, container := range containers {
					if tc.expectedValue == nil {
						assert.True(t, (container.SecurityContext == nil || container.SecurityContext.RunAsNonRoot == nil))
					} else {
						assert.Equal(t, *tc.expectedValue, *container.SecurityContext.RunAsNonRoot)
					}
				}
			}
		})
	}

	file := "run-as-non-root-psc-false-allowed-multi-containers-multi-labels.yml"
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
