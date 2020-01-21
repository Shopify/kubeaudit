package apparmor

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

const fixtureDir = "fixtures"

func TestAuditAppArmor(t *testing.T) {
	cases := []struct {
		file           string
		expectedErrors []string
	}{
		{"apparmor_enabled_v1.yml", nil},
		{"apparmor_annotation_missing_v1.yml", []string{AppArmorAnnotationMissing}},
		{"apparmor_annotation_missing_multiple_resources_v1.yml", []string{AppArmorAnnotationMissing}},
		{"apparmor_disabled_v1.yml", []string{AppArmorDisabled}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, fixtureDir, tt.file, New(), tt.expectedErrors)
		})
	}
}

func TestFixAppArmor(t *testing.T) {
	cases := []struct {
		file                    string
		expectedAnnotationValue string
	}{
		{"apparmor_enabled_v1.yml", "localhost/something"},
		{"apparmor_annotation_missing_v1.yml", ProfileRuntimeDefault},
		{"apparmor_annotation_missing_multiple_resources_v1.yml", ProfileRuntimeDefault},
		{"apparmor_disabled_v1.yml", ProfileRuntimeDefault},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			resources, _ := test.FixSetup(t, fixtureDir, tt.file, New())
			for _, resource := range resources {
				containers := k8s.GetContainers(resource)
				annotations := k8s.GetAnnotations(resource)

				for _, container := range containers {
					containerAnnotation := getContainerAnnotation(container)
					assert.Equal(t, tt.expectedAnnotationValue, annotations[containerAnnotation])
				}
			}
		})
	}
}
