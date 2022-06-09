// Package kubeaudit provides methods to find and fix security issues in Kubernetes resources.
//
// Modes
//
// Kubeaudit supports three different modes. The mode used depends on the audit method used.
//
// 1. Manifest mode: Audit a manifest file
//
// 2. Local mode: Audit resources in a local kubeconfig file
//
// 3. Cluster mode: Audit resources in a running cluster (kubeaudit must be invoked from a container within the cluster)
//
// In manifest mode, kubeaudit can automatically fix security issues.
//
// Follow the instructions below to use kubeaudit:
//
// First initialize the security auditors
//
// The auditors determine which security issues kubeaudit will look for. Each auditor is responsible for a different
// security issue. For an explanation of what each auditor checks for, see https://github.com/Shopify/kubeaudit#auditors.
//
// To initialize all available auditors:
//
//   import "github.com/Shopify/kubeaudit/auditors/all"
//
//   auditors, err := all.Auditors(config.KubeauditConfig{})
//
// Or, to initialize specific auditors, import each one:
//
//   import (
//     "github.com/Shopify/kubeaudit/auditors/apparmor"
//     "github.com/Shopify/kubeaudit/auditors/image"
//   )
//
//   auditors := []kubeaudit.Auditable{
//     apparmor.New(),
//     image.New(image.Config{Image: "myimage:mytag"}),
//   }
//
// Initialize Kubeaudit
//
// Create a new instance of kubeaudit:
//
//   kubeAuditor, err := kubeaudit.New(auditors)
//
// Run the audit
//
// To run the audit in manifest mode:
//
//   import "os"
//
//   manifest, err := os.Open("/path/to/manifest.yaml")
//   if err != nil {
//     ...
//   }
//
//   report, err := kubeAuditor.AuditManifest(manifest)
//
// Or, to run the audit in local mode:
//
//   report, err := kubeAuditor.AuditLocal("/path/to/kubeconfig.yml", kubeaudit.AuditOptions{})
//
// Or, to run the audit in cluster mode (pass it a namespace name as a string to only audit resources in that namespace, or an empty string to audit resources in all namespaces):
//
//   report, err := auditor.AuditCluster(kubeaudit.AuditOptions{})
//
// Get the results
//
// To print the results in a human readable way:
//
//   report.PrintResults()
//
// Results are printed to standard out by default. To print to a string instead:
//
//   var buf bytes.Buffer
//   report.PrintResults(kubeaudit.WithWriter(&buf), kubeaudit.WithColor(false))
//   resultsString := buf.String()
//
// Or, to get the result objects:
//
//   results := report.Results()
//
// Autofix
//
// Note that autofixing is only supported in manifest mode.
//
// To print the plan (what will be fixed):
//
//  report.PrintPlan(os.Stdout)
//
// To automatically fix the security issues and print the fixed manifest:
//
//   err = report.Fix(os.Stdout)
//
// Override Errors
//
// Overrides can be used to ignore specific auditors for specific containers or pods.
// See the documentation for the specific auditor you wish to override at https://github.com/Shopify/kubeaudit#auditors.
//
// Custom Auditors
//
// Kubeaudit supports custom auditors. See the Custom Auditor example.
//
package kubeaudit

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/Shopify/kubeaudit/internal/k8sinternal"
	"github.com/Shopify/kubeaudit/pkg/k8s"
)

// Kubeaudit provides functions to audit and fix Kubernetes manifests
type Kubeaudit struct {
	auditors []Auditable
}

type AuditOptions = k8sinternal.ClientOptions

// New returns a new Kubeaudit instance
func New(auditors []Auditable, opts ...Option) (*Kubeaudit, error) {
	if len(auditors) == 0 {
		return nil, errors.New("no auditors enabled")
	}

	auditor := &Kubeaudit{
		auditors: auditors,
	}

	if err := auditor.parseOptions(opts); err != nil {
		return nil, err
	}

	return auditor, nil
}

// AuditManifest audits the Kubernetes resources in the provided manifest
func (a *Kubeaudit) AuditManifest(manifest io.Reader) (*Report, error) {
	manifestBytes, err := ioutil.ReadAll(manifest)
	if err != nil {
		return nil, err
	}

	resources, err := getResourcesFromManifest(manifestBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources from manifest: %w", err)
	}

	results, err := auditResources(resources, a.auditors)
	if err != nil {
		return nil, err
	}

	report := &Report{results: results}

	return report, nil
}

// AuditCluster audits the Kubernetes resources found in the cluster in which Kubeaudit is running
func (a *Kubeaudit) AuditCluster(options AuditOptions) (*Report, error) {
	if !k8sinternal.IsRunningInCluster(k8sinternal.DefaultClient) {
		return nil, errors.New("failed to audit resources in cluster mode: not running in cluster")
	}

	client, err := k8sinternal.NewKubeClientCluster(k8sinternal.DefaultClient)
	if err != nil {
		return nil, err
	}

	resources, err := getResourcesFromClient(client, options)
	if err != nil {
		return nil, err
	}
	results, err := auditResources(resources, a.auditors)
	if err != nil {
		return nil, err
	}

	report := &Report{results: results}

	return report, nil
}

// AuditLocal audits the Kubernetes resources found in the provided Kubernetes config file
func (a *Kubeaudit) AuditLocal(configpath string, options AuditOptions) (*Report, error) {
	client, err := k8sinternal.NewKubeClientLocal(configpath)
	if err == k8sinternal.ErrNoReadableKubeConfig {
		return nil, fmt.Errorf("failed to open kubeconfig file %s", configpath)
	} else if err != nil {
		return nil, err
	}

	resources, err := getResourcesFromClient(client, options)
	if err != nil {
		return nil, err
	}
	results, err := auditResources(resources, a.auditors)
	if err != nil {
		return nil, err
	}

	report := &Report{results: results}

	return report, nil
}

// Report contains the results after auditing
type Report struct {
	results []Result
}

// RawResults returns all of the results for each Kubernetes resource, including ones that had no audit results.
// Generally, you will want to use Results() instead.
func (r *Report) RawResults() []Result {
	return r.results
}

// Results returns the audit results for each Kubernetes resource
func (r *Report) Results() []Result {
	results := make([]Result, 0, len(r.results))
	for _, result := range r.results {
		if len(result.GetAuditResults()) > 0 {
			results = append(results, result)
		}
	}
	return results
}

// ResultsWithMinSeverity returns the audit results for each Kubernetes resource with a minimum severity
func (r *Report) ResultsWithMinSeverity(minSeverity SeverityLevel) []Result {
	var results []Result
	for _, result := range r.results {
		var filteredAuditResults []*AuditResult
		for _, auditResult := range result.GetAuditResults() {
			if auditResult.Severity >= minSeverity {
				filteredAuditResults = append(filteredAuditResults, auditResult)
			}
		}
		if len(filteredAuditResults) > 0 {
			results = append(results, &workloadResult{
				Resource:     result.GetResource(),
				AuditResults: filteredAuditResults,
			})
		}
	}
	return results
}

// HasErrors returns true if any findings have the level of Error
func (r *Report) HasErrors() (errorsFound bool) {
	for _, workloadResult := range r.Results() {
		for _, auditResult := range workloadResult.GetAuditResults() {
			if auditResult.Severity >= Error {
				return true
			}
		}
	}
	return false
}

// PrintResults writes the audit results to the specified writer. Defaults to printing results to stdout
func (r *Report) PrintResults(printOptions ...PrintOption) {
	printer := NewPrinter(printOptions...)
	printer.PrintReport(r)
}

// Fix tries to automatically patch any security concerns and writes the resulting manifest to the provided writer.
// Only applies when audit was performed on a manifest (not local or cluster)
func (r *Report) Fix(writer io.Writer) error {
	fixed, err := fix(r.results)
	if err != nil {
		return err
	}

	_, err = writer.Write(fixed)
	return err
}

// PrintPlan writes the actions that will be performed by the Fix() function in a human-readable way to the
// provided writer. Only applies when audit was performed on a manifest (not local or cluster)
func (r *Report) PrintPlan(writer io.Writer) {
	for _, result := range r.Results() {
		for _, auditResult := range result.GetAuditResults() {
			ok, plan := auditResult.FixPlan()
			if ok {
				fmt.Fprintln(writer, "* ", plan)
			}
		}
	}
}

// Auditable is an interface which is implemented by auditors
type Auditable interface {
	Audit(resource k8s.Resource, resources []k8s.Resource) ([]*AuditResult, error)
}
