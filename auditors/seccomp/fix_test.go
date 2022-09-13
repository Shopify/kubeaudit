package seccomp

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiv1 "k8s.io/api/core/v1"
)

const fixtureDir = "fixtures"
const emptyProfile = apiv1.SeccompProfileType("EMPTY")
const defaultProfile = apiv1.SeccompProfileTypeRuntimeDefault
const localhostProfile = apiv1.SeccompProfileTypeLocalhost

func TestFixSeccomp(t *testing.T) {
	cases := []struct {
		file                             string
		expectedPodSeccompProfile        apiv1.SeccompProfileType
		expectedContainerSeccompProfiles []apiv1.SeccompProfileType
	}{
		{"seccomp-profile-missing.yml", defaultProfile, []apiv1.SeccompProfileType{emptyProfile}},
		{"seccomp-profile-missing-disabled-container.yml", defaultProfile, []apiv1.SeccompProfileType{emptyProfile}},
		{"seccomp-profile-missing-annotations.yml", defaultProfile, []apiv1.SeccompProfileType{emptyProfile}},
		{"seccomp-disabled-pod.yml", defaultProfile, []apiv1.SeccompProfileType{defaultProfile}},
		{"seccomp-disabled.yml", defaultProfile, []apiv1.SeccompProfileType{emptyProfile, emptyProfile}},
		{"seccomp-disabled-localhost.yml", localhostProfile, []apiv1.SeccompProfileType{defaultProfile, emptyProfile}},
	}

	for _, tc := range cases {
		// This line is needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, fixtureDir, tc.file, New())
			require.Len(t, resources, 1)
			resource := resources[0]

			updatedPodSpec := k8s.GetPodSpec(resource)
			checkPodSeccompProfile(t, updatedPodSpec, tc.expectedPodSeccompProfile)
			checkContainerSeccompProfiles(t, updatedPodSpec, tc.expectedContainerSeccompProfiles)
			checkNoSeccompAnnotations(t, resource)
		})
	}
}

func checkPodSeccompProfile(t *testing.T, podSpec *apiv1.PodSpec, expectedPodSeccompProfile apiv1.SeccompProfileType) {
	securityContext := podSpec.SecurityContext
	if expectedPodSeccompProfile == emptyProfile {
		require.Nil(t, securityContext)
	} else {
		assert.Equal(t, expectedPodSeccompProfile, securityContext.SeccompProfile.Type)
	}
}

func checkContainerSeccompProfiles(t *testing.T, podSpec *apiv1.PodSpec, expectedContainerSeccompProfiles []apiv1.SeccompProfileType) {
	for i, container := range podSpec.Containers {
		securityContext := container.SecurityContext
		expectedProfile := expectedContainerSeccompProfiles[i]
		if expectedProfile == emptyProfile {
			require.True(t, securityContext == nil || securityContext.SeccompProfile == nil)
		} else {
			assert.Equal(t, expectedProfile, securityContext.SeccompProfile.Type)
		}
	}
}

func checkNoSeccompAnnotations(t *testing.T, resource k8s.Resource) {
	annotations := k8s.GetAnnotations(resource)
	if annotations == nil {
		return
	}

	seccompAnnotations := []string{}
	for annotation := range annotations {
		if annotation == PodAnnotationKey || strings.HasPrefix(annotation, ContainerAnnotationKeyPrefix) {
			seccompAnnotations = append(seccompAnnotations, annotation)
		}
	}
	assert.Empty(t, seccompAnnotations)
}
