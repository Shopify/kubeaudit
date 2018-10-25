package cmd

import (
	apiv1 "k8s.io/api/core/v1"
	"testing"
)

func TestFixSeccompDisabled(t *testing.T) {
	fileToFix := "seccomp_disabled.yml"
	fileFixed := "seccomp_disabled-fixed.yml"
	assertEqualYaml(fileToFix, fileFixed, auditSeccomp, t)
}

func TestFixSeccompDisabled2(t *testing.T) {
	fileToFix := "seccomp_disabled_2.yml"
	fileFixed := "seccomp_disabled_2-fixed.yml"
	assertEqualYaml(fileToFix, fileFixed, auditSeccomp, t)
}

func TestFixSeccompDisabledPod(t *testing.T) {
	testFixSeccomp(t, "seccomp_disabled_pod.yml")
}

func TestFixSeccompAnnotationMissing(t *testing.T) {
	testFixSeccomp(t, "seccomp_annotation_missing.yml")
}

func testFixSeccomp(t *testing.T, configFile string) {
	assert, resource := FixTestSetup(t, configFile, auditSeccomp)
	annotations := getPodAnnotations(resource)
	podVal, podOk := annotations[apiv1.SeccompPodAnnotationKey]

	assert.True(podOk)
	assert.False(badSeccompProfileName(podVal))

	for _, container := range getContainers(resource) {
		containerAnnotation := apiv1.SeccompContainerAnnotationKeyPrefix + container.Name
		if val, ok := annotations[containerAnnotation]; ok {
			assert.False(badSeccompProfileName(val))
		}
	}
}
