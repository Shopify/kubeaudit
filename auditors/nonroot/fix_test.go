package nonroot

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
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
		{"run-as-user-0.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-user-0-allowed.yml", fixtureDir, nil},
		{"run-as-user-psc-0.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-user-psc-0-allowed.yml", fixtureDir, nil},
		{"run-as-user-psc-1.yml", fixtureDir, nil},
		{"run-as-user-psc-1-csc-0.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-user-psc-0-csc-0.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-user-psc-0-csc-1.yml", fixtureDir, nil},
		{"run-as-user-redundant-override-container.yml", fixtureDir, nil},
		{"run-as-user-redundant-override-pod.yml", fixtureDir, nil},
		{"run-as-user-psc-0-csc-1-multiple-cont.yml", fixtureDir, nil},
		{"run-as-user-psc-0-allowed-multi-containers-multi-labels.yml", fixtureDir, nil},
		{"run-as-user-psc-0-allowed-multi-containers-single-label.yml", fixtureDir, nil},
		{"run-as-user-0-run-as-non-root-true.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-user-0-run-as-non-root-false.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-user-psc-0-run-as-non-root-psc-true.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-user-psc-0-run-as-non-root-psc-false.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-user-1-run-as-non-root-true.yml", fixtureDir, k8s.NewTrue()},
		{"run-as-user-1-run-as-non-root-false.yml", fixtureDir, k8s.NewFalse()},
		{"run-as-user-psc-1-run-as-non-root-psc-true.yml", fixtureDir, nil},
		{"run-as-user-psc-1-run-as-non-root-psc-false.yml", fixtureDir, nil},
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

	files := []string{
		"run-as-non-root-psc-false-allowed-multi-containers-multi-labels.yml",
		"run-as-user-psc-0-csc-nil-multiple-cont.yml",
	}
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, fixtureDir, file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				for _, container := range containers {
					switch container.Name {
					case "container1":
						assert.True(t, (container.SecurityContext == nil || container.SecurityContext.RunAsNonRoot == nil))
					case "container2":
						assert.True(t, *container.SecurityContext.RunAsNonRoot)
					}
				}
			}
		})
	}
}
