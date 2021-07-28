package privileged

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
)

func TestFixPrivileged(t *testing.T) {
	cases := []struct {
		file          string
		fixtureDir    string
		expectedValue bool
	}{
		{"privileged-nil.yml", fixtureDir, false},
		{"privileged-true.yml", fixtureDir, false},
		{"privileged-true-allowed.yml", fixtureDir, true},
		{"privileged-redundant-override.yml", fixtureDir, false},
		{"privileged-true-allowed-multi-containers-multi-labels.yml", fixtureDir, true},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, tc.fixtureDir, tc.file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				for _, container := range containers {
					assert.Equal(t, tc.expectedValue, *container.SecurityContext.Privileged)
				}
			}
		})
	}

	file := "privileged-true-allowed-multi-containers-single-label.yml"
	t.Run(file, func(t *testing.T) {
		resources, _ := test.FixSetup(t, fixtureDir, file, New())
		for _, resource := range resources {
			containers := k8s.GetContainers(resource)
			for _, container := range containers {
				switch container.Name {
				case "container1":
					assert.False(t, *container.SecurityContext.Privileged)
				case "container2":
					assert.True(t, *container.SecurityContext.Privileged)
				default:
					assert.Failf(t, "unexpected container name", container.Name)
				}
			}
		}
	})
}
