package kubeaudit

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
	log "github.com/sirupsen/logrus"
)

// Kubeaudit provides functions to audit and fix Kubernetes manifests
type Kubeaudit struct {
	auditors []Auditable
}

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
func (a *Kubeaudit) AuditCluster(namespace string) (*Report, error) {
	if !k8s.IsRunningInCluster(k8s.DefaultClient) {
		return nil, errors.New("failed to audit resources in cluster mode: not running in cluster")
	}

	clientset, err := k8s.NewKubeClientCluster(k8s.DefaultClient)
	if err != nil {
		return nil, err
	}

	resources := getResourcesFromClientset(clientset, namespace)
	results, err := auditResources(resources, a.auditors)
	if err != nil {
		return nil, err
	}

	report := &Report{results: results}

	return report, nil
}

// AuditLocal audits the Kubernetes resources found in the provided Kubernetes config file
func (a *Kubeaudit) AuditLocal(configpath string) (*Report, error) {
	if _, err := os.Stat(configpath); err != nil {
		err = fmt.Errorf("failed to open kubeconfig file %s", configpath)
		return nil, err
	}

	clientset, err := k8s.NewKubeClientLocal(configpath)
	if err != nil {
		return nil, err
	}

	resources := getResourcesFromClientset(clientset, "")
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

// RawResults returns all of the results for each Kubernetes resource
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

// PrintResults writes the audit results with a severity greater than or matching minSeverity in a human-readable
// way to the provided writer
func (r *Report) PrintResults(writer io.Writer, minSeverity int, formatter log.Formatter) {
	resultLogger := log.New()

	resultLogger.SetOutput(writer)
	if formatter != nil {
		resultLogger.SetFormatter(formatter)
	}

	// We manually manage what severity levels to log, lorgus should let everything through
	resultLogger.SetLevel(log.DebugLevel)

	for _, workloadResult := range r.Results() {
		for _, auditResult := range workloadResult.GetAuditResults() {
			if auditResult.Severity >= minSeverity {
				logAuditResult(auditResult, resultLogger)
			}
		}
	}
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
// TODO include metadata with the plan?
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
	Audit(resource k8stypes.Resource, resources []k8stypes.Resource) ([]*AuditResult, error)
}
