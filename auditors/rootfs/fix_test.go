package rootfs

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestFixReadOnlyRootFilesystem(t *testing.T) {
	cases := []struct {
		file          string
		fixtureDir    string
		expectedValue bool
	}{
		{"read_only_root_filesystem_nil_v1.yml", fixtureDir, true},
		{"read_only_root_filesystem_false_v1.yml", fixtureDir, true},
		{"read_only_root_filesystem_false_allowed_v1.yml", fixtureDir, false},
		{"read_only_root_filesystem_redundant_override_v1.yml", fixtureDir, true},
		{"read_only_root_filesystem_false_allowed_multi_container_multi_labels_v1.yml", fixtureDir, false},

		// Shared fixtures
		{"security_context_nil_v1.yml", test.SharedFixturesDir, true},
		{"security_context_nil_v1beta1.yml", test.SharedFixturesDir, true},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, tt.fixtureDir, tt.file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				for _, container := range containers {
					assert.Equal(t, tt.expectedValue, *container.SecurityContext.ReadOnlyRootFilesystem)
				}
			}
		})
	}

	file := "read_only_root_filesystem_false_allowed_multi_container_single_label_v1.yml"
	t.Run(file, func(t *testing.T) {
		resources, _ := test.FixSetup(t, fixtureDir, file, New())
		for _, resource := range resources {
			containers := k8s.GetContainers(resource)
			for _, container := range containers {
				switch container.Name {
				case "fakeContainerRORF1":
					assert.False(t, *container.SecurityContext.ReadOnlyRootFilesystem)
				case "fakeContainerRORF2":
					assert.True(t, *container.SecurityContext.ReadOnlyRootFilesystem)
				}
			}
		}
	})
}
