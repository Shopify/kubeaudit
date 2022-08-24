package sarif

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/auditors/limits"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateWithResults(t *testing.T) {
	cases := []struct {
		description        string
		auditResults       []*kubeaudit.AuditResult
		expectedRule       string
		expectedErrorLevel string
		expectedMessage    string
		expectedURI        string
		expectedFilePath   string
	}{
		{
			"apparmor invalid",
			[]*kubeaudit.AuditResult{{
				Auditor:  apparmor.Name,
				Rule:     apparmor.AppArmorInvalidAnnotation,
				Severity: kubeaudit.Error,
				Message:  "AppArmor annotation key refers to a container that doesn't exist",
				FilePath: "apparmorPath",
			}},
			apparmor.AppArmorInvalidAnnotation,
			"error",
			"AppArmor annotation key refers to a container that doesn't exist",
			"https://github.com/Shopify/kubeaudit/blob/main/docs/auditors/apparmor.md",
			"apparmorPath",
		},
		{
			"capabilities added",
			[]*kubeaudit.AuditResult{{
				Auditor:  capabilities.Name,
				Rule:     capabilities.CapabilityAdded,
				Severity: kubeaudit.Error,
				Message:  "It should be removed from the capability add list",
				FilePath: "capsPath",
			}},
			capabilities.CapabilityAdded,
			"error",
			"It should be removed from the capability add list",
			"https://github.com/Shopify/kubeaudit/blob/main/docs/auditors/capabilities.md",
			"capsPath",
		},
		{
			"image tag is present",
			[]*kubeaudit.AuditResult{{
				Auditor:  image.Name,
				Rule:     image.ImageCorrect,
				Severity: kubeaudit.Info,
				Message:  "Image tag is correct",
				FilePath: "imagePath",
			}},
			image.ImageCorrect,
			"note",
			"Image tag is correct",
			"https://github.com/Shopify/kubeaudit/blob/main/docs/auditors/image.md",
			"imagePath",
		},
		{
			"limits is nil",
			[]*kubeaudit.AuditResult{{
				Auditor:  limits.Name,
				Rule:     limits.LimitsNotSet,
				Severity: kubeaudit.Warn,
				Message:  "Resource limits not set",
				FilePath: "limitsPath",
			}},
			limits.LimitsNotSet,
			"warning",
			"Resource limits not set",
			"https://github.com/Shopify/kubeaudit/blob/main/docs/auditors/limits.md",
			"limitsPath",
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			kubeAuditReport := kubeaudit.NewReport([]kubeaudit.Result{&kubeaudit.WorkloadResult{
				AuditResults: tc.auditResults,
			}})

			sarifReport, err := Create(kubeAuditReport)
			require.NoError(t, err)

			assert.Equal(t, repoURL,
				*sarifReport.Runs[0].Tool.Driver.InformationURI)

			// verify that the rules have been added as per report findings
			assert.Equal(t, tc.expectedRule, sarifReport.Runs[0].Tool.Driver.Rules[0].ID)

			var ruleNames []string

			//check for rules occurrences
			for _, sarifRule := range sarifReport.Runs[0].Tool.Driver.Rules {
				assert.Equal(t, []string{
					"security",
					"kubernetes",
					"infrastructure",
				},
					sarifRule.Properties["tags"],
				)

				ruleNames = append(ruleNames, sarifRule.ID)

				assert.Contains(t, *sarifRule.Help.Text, tc.expectedURI)
			}

			for _, sarifResult := range sarifReport.Runs[0].Results {
				assert.Contains(t, ruleNames, *sarifResult.RuleID)
				assert.Equal(t, tc.expectedErrorLevel, *sarifResult.Level)
				assert.Contains(t, *sarifResult.Message.Text, tc.expectedMessage)
				assert.Contains(t, tc.expectedFilePath, *sarifResult.Locations[0].PhysicalLocation.ArtifactLocation.URI)
			}
		})
	}
}

func TestCreateWithNoResults(t *testing.T) {
	sarifReport, err := Create(&kubeaudit.Report{})
	require.NoError(t, err)
	require.NotEmpty(t, *sarifReport.Runs[0])
	// verify that the rules are only added as per report findings
	assert.Len(t, sarifReport.Runs[0].Tool.Driver.Rules, 0)
}
