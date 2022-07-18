package sarif

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/Shopify/kubeaudit/auditors/seccomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	capabilitiesAuditable := capabilities.New(capabilities.Config{})
	apparmorAuditable := apparmor.New()
	seccompAuditable := seccomp.New()

	cases := []struct {
		file          string
		auditorName   string
		auditors      []kubeaudit.Auditable
		expectedRules []string
	}{
		{
			"apparmor-disabled.yaml",
			apparmor.Name,
			[]kubeaudit.Auditable{apparmorAuditable},
			[]string{"AppArmorInvalidAnnotation"},
		},
		{
			"capabilities-added.yaml",
			capabilities.Name,
			[]kubeaudit.Auditable{capabilitiesAuditable, seccompAuditable},
			[]string{"CapabilityAdded", "SeccompAnnotationMissing"},
		},
		{
			"capabilities-added.yaml",
			capabilities.Name,
			[]kubeaudit.Auditable{capabilitiesAuditable},
			[]string{"CapabilityAdded"},
		},
	}

	for _, tc := range cases {
		fixture := filepath.Join("fixtures", tc.file)
		auditor, err := kubeaudit.New(tc.auditors)
		require.NoError(t, err)

		manifest, openErr := os.Open(fixture)
		require.NoError(t, openErr)

		kubeAuditReport, err := auditor.AuditManifest(fixture, manifest)
		require.NoError(t, err)

		// we're only appending sarif to the path here for testing purposes
		// this allows us to visualize the sarif output preview correctly
		for _, reportResult := range kubeAuditReport.Results() {
			r := reportResult.GetAuditResults()

			for _, auditResult := range r {
				auditResult.FilePath = filepath.Join("sarif/", auditResult.FilePath)
			}

		}
		// verify that the file path is correct
		assert.Contains(t, kubeAuditReport.Results()[0].GetAuditResults()[0].FilePath, "sarif/fixtures")

		sarifReport, err := Create(kubeAuditReport)
		require.NoError(t, err)

		assert.Equal(t, "https://github.com/Shopify/kubeaudit",
			*sarifReport.Runs[0].Tool.Driver.InformationURI)

		// verify that the rules have been added as per report findings
		assert.Len(t, sarifReport.Runs[0].Tool.Driver.Rules, len(tc.expectedRules))

		var ruleNames []string
		// check for rules occurrences
		for _, sarifRule := range sarifReport.Runs[0].Tool.Driver.Rules {
			assert.NotEqual(t, *sarifRule.Help.Markdown, "")

			assert.Equal(t, sarifRule.Properties["tags"], []string{
				"security",
				"kubernetes",
				"infrastructure",
			})

			ruleNames = append(ruleNames, sarifRule.ID)
		}

		for _, expectedRule := range tc.expectedRules {
			assert.Contains(t, ruleNames, expectedRule)
		}

		// for _, sarifResult := range sarifReport.Runs[0].Results {
		// 	// TODO: test for severity, message and location
		// 	// assert.Equal(sarifResult.Level,
		// }

		// also add a fixture with info level so that we capture the conversion to note

	}
}

// Validates that the given file path refers to a valid SARIF file.
// Throws an error if the file is invalid.

func TestValidate(t *testing.T) {
	var reportBytes bytes.Buffer
	testSarif, err := ioutil.ReadFile("fixtures/valid.sarif")
	require.NoError(t, err)

	reportBytes.Write(testSarif)

	err, errs := Validate(&reportBytes)
	require.NoError(t, err)
	if len(errs) > 0 {
		fmt.Println(errs)
	}

	assert.Len(t, errs, 0)

}
