package sarif

import (
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

func TestNew(t *testing.T) {
	sarifReport, _, err := New()
	require.NoError(t, err)
	require.Len(t, sarifReport.Runs, 1)
	assert.Equal(t, "https://github.com/Shopify/kubeaudit",
		*sarifReport.Runs[0].Tool.Driver.InformationURI)
}

func TestCreate(t *testing.T) {
	capabilitiesAuditable := capabilities.New(capabilities.Config{})
	apparmorAuditable := apparmor.New()
	seccompAuditable := seccomp.New()

	cases := []struct {
		file              string
		auditorName       string
		auditors          []kubeaudit.Auditable
		expectedRuleCount int
		expectedRules     []string
	}{
		{
			"apparmor-disabled.yaml",
			apparmor.Name,
			[]kubeaudit.Auditable{apparmorAuditable},
			1,
			[]string{"AppArmorInvalidAnnotation"},
		},
		{
			"capabilities-added.yaml",
			capabilities.Name,
			[]kubeaudit.Auditable{capabilitiesAuditable, seccompAuditable},
			2,
			[]string{"CapabilityAdded, SeccompAnnotationMissing"},
		},
		{
			"capabilities-added.yaml",
			capabilities.Name,
			[]kubeaudit.Auditable{capabilitiesAuditable},
			1,
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

		sarifReport, sarifRun, err := New()
		require.NoError(t, err)

		Create(kubeAuditReport, sarifRun)

		// verify that the rules have been added as per report findings
		assert.Len(t, sarifReport.Runs[0].Tool.Driver.Rules, tc.expectedRuleCount)

		// check for rules occurrences
		for _, expectedRule := range tc.expectedRules {
			assert.Contains(t, expectedRule, sarifReport.Runs[0].Tool.Driver.Rules)
		}
	}
}
