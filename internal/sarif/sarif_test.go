package sarif

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSarifReport(t *testing.T) {
	sarifReport, _ := CreateSarifReport()
	require.Len(t, sarifReport.Runs, 1)
	assert.Equal(t, "https://github.com/Shopify/kubeaudit", *sarifReport.Runs[0].Tool.Driver.InformationURI)
}

func TestAddSarifResultToReport(t *testing.T) {
	capabilitiesAuditable := capabilities.New(capabilities.Config{})
	apparmorAuditable := apparmor.New()

	auditables := []kubeaudit.Auditable{capabilitiesAuditable, apparmorAuditable}

	cases := []struct {
		file    string
		auditor string
	}{
		{"apparmor-disabled.yaml", apparmor.Name},
		{"capabilities-added.yaml", capabilities.Name},
	}

	for _, tc := range cases {
		fixture := filepath.Join("fixtures", tc.file)
		auditor, err := kubeaudit.New(auditables)
		require.NoError(t, err)

		manifest, openErr := os.Open(fixture)
		require.NoError(t, openErr)

		kubeAuditReport, err := auditor.AuditManifest(fixture, manifest)
		require.NoError(t, err)

		// we're only appending sarif to the path here for testing purposes
		for _, reportResult := range kubeAuditReport.Results() {
			r := reportResult.GetAuditResults()

			for _, auditResult := range r {
				auditResult.FilePath = filepath.Join("sarif/", auditResult.FilePath)
			}
		}

		// verify that the file path is correct
		assert.Contains(t, kubeAuditReport.Results()[0].GetAuditResults()[0].FilePath, "sarif/fixtures")

		sarifReport, sarifRun := CreateSarifReport()

		AddSarifRules(kubeAuditReport, sarifRun)

		// verify that the 2 rules (auditors) enabled have been added to the report
		require.Len(t, *&sarifReport.Runs[0].Tool.Driver.Rules, 2)

		AddSarifResult(kubeAuditReport, sarifRun)
	}
}

// todo: verify the contents of the report
