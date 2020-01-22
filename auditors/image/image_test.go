package image

import (
	"testing"

	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

const fixtureDir = "fixtures"

func TestSplitImageString(t *testing.T) {
	cases := []struct {
		testName     string
		image        string
		expectedName string
		expectedTag  string
	}{
		{"Correct image and tag", "myimage:mytag", "myimage", "mytag"},
		{"No tag", "myimage", "myimage", ""},
		{"No image", ":mytag", "", "mytag"},
		{"Empty string", "", "", ""},
	}

	for _, tt := range cases {
		t.Run(tt.testName, func(t *testing.T) {
			image, tag := splitImageString(tt.image)
			assert.Equal(t, tt.expectedName, image)
			assert.Equal(t, tt.expectedTag, tag)
		})
	}
}

func TestAuditImage(t *testing.T) {
	cases := []struct {
		file           string
		image          string
		expectedErrors []string
	}{
		{"image_tag_missing_v1.yml", "fakeContainerImg:1.6", []string{ImageTagMissing}},
		{"image_tag_present_v1.yml", "fakeContainerImg:1.6", []string{ImageTagIncorrect}},
		{"image_tag_present_v1.yml", "fakeContainerImg:1.5", []string{ImageCorrect}},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			test.Audit(t, fixtureDir, tt.file, New(tt.image), tt.expectedErrors)
		})
	}
}
