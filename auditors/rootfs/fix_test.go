package rootfs

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
)

func TestFixReadOnlyRootFilesystem(t *testing.T) {
	cases := []struct {
		file          string
		fixtureDir    string
		expectedValue bool
	}{
		{"read-only-root-filesystem-nil.yml", fixtureDir, true},
		{"read-only-root-filesystem-false.yml", fixtureDir, true},
		{"read-only-root-filesystem-false-allowed.yml", fixtureDir, false},
		{"read-only-root-filesystem-redundant-override.yml", fixtureDir, true},
		{"read-only-root-filesystem-false-allowed-multi-labels.yml", fixtureDir, false},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, tc.fixtureDir, tc.file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				for _, container := range containers {
					assert.Equal(t, tc.expectedValue, *container.SecurityContext.ReadOnlyRootFilesystem)
				}
			}
		})
	}

	file := "read-only-root-filesystem-false-allowed-single-label.yml"
	t.Run(file, func(t *testing.T) {
		resources, _ := test.FixSetup(t, fixtureDir, file, New())
		for _, resource := range resources {
			containers := k8s.GetContainers(resource)
			for _, container := range containers {
				switch container.Name {
				case "container1":
					assert.False(t, *container.SecurityContext.ReadOnlyRootFilesystem)
				case "container2":
					assert.True(t, *container.SecurityContext.ReadOnlyRootFilesystem)
				default:
					assert.Failf(t, "unexpected container name", container.Name)
				}
			}
		}
	})
}
