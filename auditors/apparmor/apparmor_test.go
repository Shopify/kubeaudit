package apparmor

import (
	"strings"
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
)

const fixtureDir = "fixtures"

func TestAuditAppArmor(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
		testLocalMode  bool
	}{
		{"apparmor-enabled.yml", nil, true},
		{"apparmor-annotation-missing.yml", []string{AppArmorAnnotationMissing}, true},
		{"apparmor-annotation-init-container-enabled.yml", nil, true},
		{"apparmor-annotation-init-container-missing.yml", []string{AppArmorAnnotationMissing}, true},
		// These are invalid manifests so we should only test it in manifest mode as kubernetes will fail to apply it
		{"apparmor-disabled.yml", []string{AppArmorDisabled}, false},
		{"apparmor-invalid-annotation.yml", []string{AppArmorInvalidAnnotation}, false},
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

func TestFixAppArmor(t *testing.T) {
	cases := []struct {
		file                    string
		expectedAnnotationValue string
	}{
		{"apparmor-enabled.yml", "localhost/something"},
		{"apparmor-annotation-missing.yml", ProfileRuntimeDefault},
		{"apparmor-disabled.yml", ProfileRuntimeDefault},
		{"apparmor-invalid-annotation.yml", ProfileRuntimeDefault},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, fixtureDir, tc.file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				annotations := k8s.GetAnnotations(resource)

				for _, container := range containers {
					containerAnnotation := getContainerAnnotation(container)
					assert.Equal(t, tc.expectedAnnotationValue, annotations[containerAnnotation])
				}
			}
		})
	}
}
