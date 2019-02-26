package cmd

import "testing"

func TestFixRunAsNonRootV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_false_v1.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsNonRootNilV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_nil_v1.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsNonRootFalseV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_false_v1.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsRootFalseAllowedV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_false_allowed_v1.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.False(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsNonRootMisconfiguredAllowV1(t *testing.T) {
	assert, resource := FixTestSetup(t, "run_as_non_root_misconfigured_allow_container_v1.yml", auditRunAsNonRoot)
	for _, container := range getContainers(resource) {
		assert.True(*container.SecurityContext.RunAsNonRoot)
	}
}

func TestFixRunAsRootFalseAllowedMultiContainersMultiLabelsV1(t *testing.T) {
	assert, resources := FixTestSetupMultipleResources(t, "run_as_non_root_psc_false_allowed_multi_containers_multi_labels_v1.yml", auditRunAsNonRoot)
	for _, resource := range resources {
		for _, container := range getContainers(resource) {
			switch container.Name {
			case "fakeContainerRANR":
				assert.Nil(container.SecurityContext.RunAsNonRoot)
			case "fakeContainerRANR2":
				assert.True(*container.SecurityContext.RunAsNonRoot)
			}
		}
	}
}

func TestFixRunAsRootFalseAllowedMultiContainersSingleLabelV1(t *testing.T) {
	assert, resources := FixTestSetupMultipleResources(t, "run_as_non_root_psc_false_allowed_multi_containers_single_label_v1.yml", auditRunAsNonRoot)
	for _, resource := range resources {
		for _, container := range getContainers(resource) {
			assert.True(*container.SecurityContext.RunAsNonRoot)
		}
	}
}
