package cmd

import "testing"

func TestFixAppArmorDisabledV1(t *testing.T) {
	testFixAppArmor(t, "apparmor_disabled_v1.yml")
}

func TestFixAppArmorAnnotationMissingV1(t *testing.T) {
	testFixAppArmor(t, "apparmor_annotation_missing_v1.yml")
}

func testFixAppArmor(t *testing.T, configFile string) {
	assert, resource := FixTestSetup(t, configFile, auditAppArmor)
	containers := getContainers(resource)
	annotations := getPodAnnotations(resource)

	for _, container := range containers {
		containerAnnotation := ContainerAnnotationKeyPrefix + container.Name
		assert.Equal(ProfileRuntimeDefault, annotations[containerAnnotation])
	}
}
