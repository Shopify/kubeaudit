package seccomp

import (
	"fmt"
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/k8stypes"
	"github.com/stretchr/testify/assert"
)

const fixtureDir = "fixtures"

func TestAuditSeccomp(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
	}{
		{"seccomp_annotation_missing_v1.yml", []string{SeccompAnnotationMissing}},
		{"seccomp_deprecated_pod_v1.yml", []string{SeccompDeprecatedPod}},
		{"seccomp_deprecated_v1.yml", []string{SeccompDeprecatedContainer, SeccompAnnotationMissing}},
		{"seccomp_disabled_pod_v1.yml", []string{SeccompDisabledPod}},
		{"seccomp_disabled_v1.yml", []string{SeccompDisabledContainer}},
		{"seccomp_enabled_pod_v1.yml", nil},
		{"seccomp_enabled_v1.yml", nil},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, fixtureDir, tt.file, New(), tt.expectedErrors)
		})
	}
}

func TestFixSeccomp(t *testing.T) {
	cases := []struct {
		testName            string
		containerNames      []string
		annotations         map[string]string
		expectedAnnotations map[string]string
	}{
		{
			testName:            "Annotation missing",
			containerNames:      []string{"container1"},
			annotations:         map[string]string{},
			expectedAnnotations: map[string]string{PodAnnotationKey: ProfileRuntimeDefault},
		},
		{
			testName:            "Annotation missing container",
			containerNames:      []string{"container1"},
			annotations:         map[string]string{PodAnnotationKey: ProfileRuntimeDefault},
			expectedAnnotations: map[string]string{PodAnnotationKey: ProfileRuntimeDefault},
		},
		{
			testName:            "Seccomp enabled pod",
			containerNames:      []string{"container1"},
			annotations:         map[string]string{PodAnnotationKey: ProfileRuntimeDefault},
			expectedAnnotations: map[string]string{PodAnnotationKey: ProfileRuntimeDefault},
		},
		{
			testName:       "Seccomp enabled container",
			containerNames: []string{"container1"},
			annotations: map[string]string{
				ContainerAnnotationKeyPrefix + "container1": ProfileRuntimeDefault,
			},
			expectedAnnotations: map[string]string{
				ContainerAnnotationKeyPrefix + "container1": ProfileRuntimeDefault,
			},
		},
		{
			testName:       "Seccomp enabled pod and container",
			containerNames: []string{"container1", "container2", "container3"},
			annotations: map[string]string{
				PodAnnotationKey: ProfileNamePrefix + "myprofile",
				ContainerAnnotationKeyPrefix + "container1": ProfileRuntimeDefault,
				ContainerAnnotationKeyPrefix + "container2": ProfileNamePrefix + "containerprofile",
			},
			expectedAnnotations: map[string]string{
				PodAnnotationKey: ProfileNamePrefix + "myprofile",
				ContainerAnnotationKeyPrefix + "container1": ProfileRuntimeDefault,
				ContainerAnnotationKeyPrefix + "container2": ProfileNamePrefix + "containerprofile",
			},
		},
		{
			testName:            "Seccomp disabled pod",
			containerNames:      []string{"container1"},
			annotations:         map[string]string{PodAnnotationKey: "badprofile"},
			expectedAnnotations: map[string]string{PodAnnotationKey: ProfileRuntimeDefault},
		},
		{
			testName:       "Seccomp disabled container",
			containerNames: []string{"container1", "container2"},
			annotations: map[string]string{
				ContainerAnnotationKeyPrefix + "container1": "badprofile",
			},
			expectedAnnotations: map[string]string{
				PodAnnotationKey: ProfileRuntimeDefault,
			},
		},
		{
			testName:       "Seccomp disabled pod",
			containerNames: []string{"container1", "container2"},
			annotations: map[string]string{
				PodAnnotationKey: "badprofile",
			},
			expectedAnnotations: map[string]string{
				PodAnnotationKey: ProfileRuntimeDefault,
			},
		},
		{
			testName:       "Seccomp disabled container enabled pod",
			containerNames: []string{"container1"},
			annotations: map[string]string{
				PodAnnotationKey: ProfileRuntimeDefault,
				ContainerAnnotationKeyPrefix + "container1": "badprofile",
			},
			expectedAnnotations: map[string]string{
				PodAnnotationKey: ProfileRuntimeDefault,
			},
		},
		{
			testName:       "Seccomp disabled container and pod",
			containerNames: []string{"container1", "container2"},
			annotations: map[string]string{
				PodAnnotationKey: "badprofile",
				ContainerAnnotationKeyPrefix + "container1": "badprofile",
			},
			expectedAnnotations: map[string]string{
				PodAnnotationKey: ProfileRuntimeDefault,
			},
		},
		{
			testName:       "Seccomp enabled container disabled pod",
			containerNames: []string{"container1", "container2"},
			annotations: map[string]string{
				PodAnnotationKey: "badprofile",
				ContainerAnnotationKeyPrefix + "container1": ProfileRuntimeDefault,
			},
			expectedAnnotations: map[string]string{
				PodAnnotationKey: ProfileRuntimeDefault,
				ContainerAnnotationKeyPrefix + "container1": ProfileRuntimeDefault,
			},
		},
	}

	auditor := New()
	for _, tt := range cases {
		t.Run(tt.testName, func(t *testing.T) {
			resource := newPod(tt.containerNames, tt.annotations)
			auditResults, err := auditor.Audit(resource, nil)
			if !assert.Nil(t, err) {
				return
			}

			for _, auditResult := range auditResults {
				auditResult.Fix(resource)
				ok, plan := auditResult.FixPlan()
				if ok {
					fmt.Println(plan)
				}
			}

			fixedAnnotations := k8s.GetAnnotations(resource)
			assert.Equal(t, tt.expectedAnnotations, fixedAnnotations)
		})
	}
}

func newPod(containerNames []string, annotations map[string]string) k8stypes.Resource {
	pod := k8stypes.NewPod()
	containers := make([]k8stypes.ContainerV1, 0, len(containerNames))
	for _, containerName := range containerNames {
		containers = append(containers, k8stypes.ContainerV1{
			Name: containerName,
		})
	}
	k8s.GetObjectMeta(pod).SetAnnotations(annotations)
	k8s.GetPodSpec(pod).Containers = containers
	return pod
}
