package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8sinternal"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SharedFixturesDir contains fixtures used by multiple tests
const SharedFixturesDir = "../../internal/test/fixtures"

const MANIFEST_MODE = "manifest"
const LOCAL_MODE = "local"

func AuditManifest(t *testing.T, fixtureDir, fixture string, auditable kubeaudit.Auditable, expectedErrors []string) {
	AuditMultiple(t, fixtureDir, fixture, []kubeaudit.Auditable{auditable}, expectedErrors, "", MANIFEST_MODE)
}

func AuditLocal(t *testing.T, fixtureDir, fixture string, auditable kubeaudit.Auditable, namespace string, expectedErrors []string) {
	AuditMultiple(t, fixtureDir, fixture, []kubeaudit.Auditable{auditable}, expectedErrors, namespace, LOCAL_MODE)
}

func AuditMultiple(t *testing.T, fixtureDir, fixture string, auditables []kubeaudit.Auditable, expectedErrors []string, namespace string, mode string) {
	if mode == LOCAL_MODE && !UseKind() {
		return
	}

	expected := make(map[string]bool, len(expectedErrors))
	for _, err := range expectedErrors {
		expected[err] = true
	}

	report := GetReport(t, fixtureDir, fixture, auditables, namespace, mode)
	require.NotNil(t, report)

	errors := make(map[string]bool)
	for _, result := range report.Results() {
		for _, auditResult := range result.GetAuditResults() {
			errors[auditResult.Name] = true
		}
	}

	assert.Equal(t, expected, errors)
}

func FixSetup(t *testing.T, fixtureDir, fixture string, auditable kubeaudit.Auditable) (fixedResources []k8s.Resource, report *kubeaudit.Report) {
	return FixSetupMultiple(t, fixtureDir, fixture, []kubeaudit.Auditable{auditable})
}

// FixSetup runs Fix() on a given manifest and returns the resulting resources
func FixSetupMultiple(t *testing.T, fixtureDir, fixture string, auditables []kubeaudit.Auditable) (fixedResources []k8s.Resource, report *kubeaudit.Report) {
	require := require.New(t)

	report = GetReport(t, fixtureDir, fixture, auditables, "", MANIFEST_MODE)
	require.NotNil(report)

	// This increases code coverage by calling the Plan() method on each PendingFix object. Plan() returns a human
	// readable description of what Apply() will do and does not need to be otherwise tested for correct logic
	report.PrintPlan(os.Stdout)

	// New resources that are created to fix the manifest are not added to the results, only the fixed manifest. By
	// running the fixed manifest through another audit run, the new resource is treated the same as any other
	// resource and parsed into a Result
	fixedManifest := bytes.NewBuffer(nil)
	err := report.Fix(fixedManifest)
	require.Nil(err)

	auditor, err := kubeaudit.New(auditables)
	require.Nil(err)

	report, err = auditor.AuditManifest(fixedManifest)
	require.Nil(err)

	results := report.RawResults()
	fixedResources = make([]k8s.Resource, 0, len(results))

	for _, result := range results {
		resource := result.GetResource().Object()
		if resource != nil {
			fixedResources = append(fixedResources, resource)
		}
	}

	return fixedResources, report
}

func GetReport(t *testing.T, fixtureDir, fixture string, auditables []kubeaudit.Auditable, namespace string, mode string) *kubeaudit.Report {
	require := require.New(t)

	fixture = filepath.Join(fixtureDir, fixture)
	auditor, err := kubeaudit.New(auditables)
	require.NoError(err)

	var report *kubeaudit.Report
	switch mode {
	case MANIFEST_MODE:
		manifest, openErr := os.Open(fixture)
		require.NoError(openErr)
		report, err = auditor.AuditManifest(manifest)
	case LOCAL_MODE:
		defer DeleteNamespace(t, namespace)
		CreateNamespace(t, namespace)
		ApplyManifest(t, fixture, namespace)
		report, err = auditor.AuditLocal("", k8sinternal.ClientOptions{Namespace: namespace})
	}

	require.NoError(err)

	return report
}

// GetAllFileNames returns all file names in the given directory
// It can be used to retrieve all of the resource manifests from the test/fixtures/all_resources directory
// This directory is not hardcoded because the working directory for tests is relative to the test
func GetAllFileNames(t *testing.T, directory string) []string {
	files, err := ioutil.ReadDir(directory)
	require.Nil(t, err)

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames
}

// UseKind returns true if tests which utilize Kind should run
func UseKind() bool {
	return os.Getenv("USE_KIND") != "false"
}

func ApplyManifest(t *testing.T, manifestPath, namespace string) {
	t.Helper()
	runCmd(t, exec.Command("kubectl", "apply", "-f", manifestPath, "-n", namespace))
}

func CreateNamespace(t *testing.T, namespace string) {
	t.Helper()
	runCmd(t, exec.Command("kubectl", "create", "namespace", namespace))
}

func DeleteNamespace(t *testing.T, namespace string) {
	t.Helper()
	runCmd(t, exec.Command("kubectl", "delete", "namespace", namespace))
}

func runCmd(t *testing.T, cmd *exec.Cmd) {
	t.Helper()
	stdoutStderr, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", stdoutStderr)
	require.NoError(t, err)
}
