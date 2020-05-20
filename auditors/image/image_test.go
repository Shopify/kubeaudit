package image

import (
	"fmt"
	"strings"
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

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			image, tag := splitImageString(tc.image)
			assert.Equal(t, tc.expectedName, image)
			assert.Equal(t, tc.expectedTag, tag)
		})
	}
}

func TestAuditImage(t *testing.T) {
	cases := []struct {
		file           string
		image          string
		expectedErrors []string
	}{
		{"image-tag-missing.yml", "scratch:1.6", []string{ImageTagMissing}},
		{"image-tag-missing.yml", "", []string{ImageTagMissing}},
		{"image-tag-present.yml", "scratch:1.6", []string{ImageTagIncorrect}},
		{"image-tag-present.yml", "", []string{}},
		{"image-tag-present.yml", "scratch:1.5", []string{ImageCorrect}},
	}

	for i, tc := range cases {
		// These lines are needed because of how scopes work with parallel tests (see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721)
		tc := tc
		i := i
		t.Run(tc.file+" "+tc.image, func(t *testing.T) {
			t.Parallel()
			test.AuditManifest(t, fixtureDir, tc.file, New(Config{Image: tc.image}), tc.expectedErrors)
			test.AuditLocal(t, fixtureDir, tc.file, New(Config{Image: tc.image}), fmt.Sprintf("%s%d", strings.Split(tc.file, ".")[0], i), tc.expectedErrors)
		})
	}
}
