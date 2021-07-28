package seccomp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
)

const fixtureDir = "fixtures"

func TestAuditSeccomp(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
		testLocalMode  bool
	}{
		{"seccomp-annotation-missing.yml", []string{SeccompAnnotationMissing}, true},
		{"seccomp-deprecated-pod.yml", []string{SeccompDeprecatedPod}, true},
		{"seccomp-deprecated.yml", []string{SeccompDeprecatedContainer, SeccompAnnotationMissing}, true},
		{"seccomp-disabled-pod.yml", []string{SeccompDisabledPod}, true},
		{"seccomp-disabled.yml", []string{SeccompDisabledContainer}, false},
		{"seccomp-enabled-pod.yml", nil, true},
		{"seccomp-enabled.yml", nil, true},
	}

	for _, tc := range cases {
		// This line is needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			test.AuditManifest(t, fixtureDir, tc.file, New(), tc.expectedErrors)
			if tc.testLocalMode {
				test.AuditLocal(t, fixtureDir, tc.file, New(), strings.Split(tc.file, ".")[0], tc.expectedErrors)
			}
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
	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			resource := newPod(tc.containerNames, tc.annotations)
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
			assert.Equal(t, tc.expectedAnnotations, fixedAnnotations)
		})
	}
}

func newPod(containerNames []string, annotations map[string]string) k8s.Resource {
	pod := k8s.NewPod()
	containers := make([]k8s.ContainerV1, 0, len(containerNames))
	for _, containerName := range containerNames {
		containers = append(containers, k8s.ContainerV1{
			Name: containerName,
		})
	}
	k8s.GetObjectMeta(pod).SetAnnotations(annotations)
	k8s.GetPodSpec(pod).Containers = containers
	return pod
}
