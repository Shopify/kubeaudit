package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/k8stypes"
	"github.com/stretchr/testify/assert"
)

// SharedFixturesDir contains fixtures used by multiple tests
const SharedFixturesDir = "../../internal/test/fixtures"

// Audit audits a given manifest using the given auditor and returns the results
func Audit(t *testing.T, fixtureDir, fixture string, auditable kubeaudit.Auditable, expectedErrors []string) []kubeaudit.Result {
	return AuditMultiple(t, fixtureDir, fixture, []kubeaudit.Auditable{auditable}, expectedErrors)
}

func AuditMultiple(t *testing.T, fixtureDir, fixture string, auditables []kubeaudit.Auditable, expectedErrors []string) []kubeaudit.Result {
	assert := assert.New(t)

	report := AuditManifest(t, fixtureDir, fixture, auditables)
	if !assert.NotNil(t, report) {
		return nil
	}

	errors := make(map[string]bool)
	for _, result := range report.Results() {
		for _, auditResult := range result.GetAuditResults() {
			errors[auditResult.Name] = true
		}
	}

	expected := make(map[string]bool, len(expectedErrors))
	for _, err := range expectedErrors {
		expected[err] = true
	}

	assert.Equal(expected, errors)

	return report.Results()
}

func FixSetup(t *testing.T, fixtureDir, fixture string, auditable kubeaudit.Auditable) (fixedResources []k8stypes.Resource, report *kubeaudit.Report) {
	return FixSetupMultiple(t, fixtureDir, fixture, []kubeaudit.Auditable{auditable})
}

// FixSetup runs Fix() on a given manifest and returns the resulting resources
func FixSetupMultiple(t *testing.T, fixtureDir, fixture string, auditables []kubeaudit.Auditable) (fixedResources []k8stypes.Resource, report *kubeaudit.Report) {
	assert := assert.New(t)

	report = AuditManifest(t, fixtureDir, fixture, auditables)
	if !assert.NotNil(report) {
		return
	}

	// This increases code coverage by calling the Plan() method on each PendingFix object. Plan() returns a human
	// readable description of what Apply() will do and does not need to be otherwise tested for correct logic
	report.PrintPlan(os.Stdout)

	// New resources that are created to fix the manifest are not added to the results, only the fixed manifest. By
	// running the fixed manifest through another audit run, the new resource is treated the same as any other
	// resource and parsed into a Result
	fixedManifest := bytes.NewBuffer(nil)
	err := report.Fix(fixedManifest)
	if !assert.Nil(err) {
		return
	}

	auditor, err := kubeaudit.New(auditables)
	if !assert.Nil(err) {
		return
	}

	report, err = auditor.AuditManifest(fixedManifest)
	if !assert.Nil(err) {
		return
	}

	results := report.RawResults()
	fixedResources = make([]k8stypes.Resource, 0, len(results))

	for _, result := range results {
		resource := result.GetResource().Object()
		if resource != nil {
			fixedResources = append(fixedResources, resource)
		}
	}

	return fixedResources, report
}

func AuditManifest(t *testing.T, fixtureDir, fixture string, auditables []kubeaudit.Auditable) *kubeaudit.Report {
	assert := assert.New(t)

	fixture = filepath.Join(fixtureDir, fixture)
	manifest, err := os.Open(fixture)
	if !assert.Nil(err) {
		return nil
	}

	auditor, err := kubeaudit.New(auditables)
	if !assert.Nil(err) {
		return nil
	}

	report, err := auditor.AuditManifest(manifest)
	if !assert.Nil(err) {
		return nil
	}

	return report
}
