package kubeaudit_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/config"
	"github.com/Shopify/kubeaudit/internal/test"
	"github.com/stretchr/testify/assert"
)

const fixtureDir = "internal/test/fixtures"

func TestFix(t *testing.T) {
	cases := []struct {
		origFile  string
		fixedFile string
	}{
		{"all-resources_v1.yml", "all-resources-fixed_v1.yml"},
		{"all-resources_v1beta1.yml", "all-resources-fixed_v1beta1.yml"},
		{"all-resources_v1beta2.yml", "all-resources-fixed_v1beta2.yml"},
	}

	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	assert.NoError(t, err)

	for _, tt := range cases {
		t.Run(tt.origFile+" <=> "+tt.fixedFile, func(t *testing.T) {
			assert := assert.New(t)

			report := test.AuditManifest(t, fixtureDir, tt.origFile, allAuditors)

			fixed := bytes.NewBuffer(nil)
			report.Fix(fixed)

			expected, err := ioutil.ReadFile(filepath.Join(fixtureDir, tt.fixedFile))
			assert.Nil(err)

			assert.Equal(string(expected), fixed.String())
		})
	}
}
