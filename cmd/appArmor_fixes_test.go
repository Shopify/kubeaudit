package cmd

import "testing"

func TestFixAppArmorDisabled(t *testing.T) {
	testFixAppArmor(t, "apparmor_disabled.yml")
}

func TestFixAppArmorAnnotationMissing(t *testing.T) {
	testFixAppArmor(t, "apparmor_annotation_missing.yml")
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
