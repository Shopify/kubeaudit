package privesc

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
)

func TestFixPrivilegeEscalation(t *testing.T) {
	cases := []struct {
		file          string
		fixtureDir    string
		expectedValue bool
	}{
		{"allow-privilege-escalation-nil.yml", fixtureDir, false},
		{"allow-privilege-escalation-redundant-override.yml", fixtureDir, false},
		{"allow-privilege-escalation-true-allowed.yml", fixtureDir, true},
		{"allow-privilege-escalation-true-multi-allowed-multi-containers.yml", fixtureDir, true},
		{"allow-privilege-escalation-true.yml", fixtureDir, false},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, tc.fixtureDir, tc.file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				for _, container := range containers {
					assert.Equal(t, tc.expectedValue, *container.SecurityContext.AllowPrivilegeEscalation)
				}
			}
		})
	}

	file := "allow-privilege-escalation-true-single-allowed-multi-containers.yml"
	t.Run(file, func(t *testing.T) {
		resources, _ := test.FixSetup(t, fixtureDir, file, New())
		for _, resource := range resources {
			containers := k8s.GetContainers(resource)
			for _, container := range containers {
				switch container.Name {
				case "container1":
					assert.False(t, *container.SecurityContext.AllowPrivilegeEscalation)
				case "container2":
					assert.True(t, *container.SecurityContext.AllowPrivilegeEscalation)
				default:
					assert.Failf(t, "unexpected container name", container.Name)
				}
			}
		}
	})
}
