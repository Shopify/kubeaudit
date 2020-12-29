package runasuser

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

var (
	uid0 int64 = 0
	uid1 int64 = 1
)

func TestFixRunAsUser(t *testing.T) {
	cases := []struct {
		file          string
		fixtureDir    string
		expectedValue *int64
	}{
		{"run-as-user-nil.yml", fixtureDir, &uid1},
		{"run-as-user-0.yml", fixtureDir, &uid1},
		{"run-as-user-0-allowed.yml", fixtureDir, &uid0},
		{"run-as-user-redundant-override-container.yml", fixtureDir, &uid1},
		{"run-as-user-redundant-override-pod.yml", fixtureDir, nil},
		{"run-as-user-psc-0.yml", fixtureDir, &uid1},
		{"run-as-user-psc-1-csc-0.yml", fixtureDir, &uid1},
		{"run-as-user-psc-0-csc-0.yml", fixtureDir, &uid1},
		{"run-as-user-psc-0-allowed.yml", fixtureDir, nil},
		{"run-as-user-psc-0-csc-1.yml", fixtureDir, &uid1},
		{"run-as-user-psc-0-csc-nil-multiple-cont.yml", fixtureDir, &uid1},
		{"run-as-user-psc-0-csc-1-multiple-cont.yml", fixtureDir, &uid1},
		{"run-as-user-psc-0-allowed-multi-containers-single-label.yml", fixtureDir, &uid1},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, tc.fixtureDir, tc.file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				for _, container := range containers {
					if tc.expectedValue == nil {
						assert.True(t, (container.SecurityContext == nil || container.SecurityContext.RunAsUser == nil))
					} else {
						assert.Equal(t, *tc.expectedValue, *container.SecurityContext.RunAsUser)
					}
				}
			}
		})
	}

	file := "run-as-user-psc-0-allowed-multi-containers-multi-labels.yml"
	t.Run(file, func(t *testing.T) {
		resources, _ := test.FixSetup(t, fixtureDir, file, New())
		for _, resource := range resources {
			containers := k8s.GetContainers(resource)
			for _, container := range containers {
				switch container.Name {
				case "fakeContainerRANR":
					assert.True(t, (container.SecurityContext == nil || container.SecurityContext.RunAsUser == nil))
				case "fakeContainerRANR2":
					assert.True(t, *container.SecurityContext.RunAsUser > 0)
				}
			}
		}
	})
}
