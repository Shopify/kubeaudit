package kubeaudit

import "github.com/Shopify/kubeaudit/k8stypes"

// AuditResult severity levels. They also correspond to log levels
const (
	// Info is used for informational audit results where no action is required
	Info = iota
	// Warn is used for audit results where there may be security concerns. If an auditor is disabled for a resource
	// using an override label, the audit results will be warnings instead of errors. Kubeaudit will NOT attempt to
	// fix these
	Warn
	// Error is used for audit results where action is required. Kubeaudit will attempt to fix these
	Error
)

// Result contains the audit results for a single Kubernetes resource
type Result interface {
	GetResource() KubeResource
	GetAuditResults() []*AuditResult
}

// AuditResult represents a potential security issue. There may be multiple AuditResults per resource and audit
type AuditResult struct {
	Name       string     // Name uniquely identifies a type of audit result
	Severity   int        // Severity is one of Error, Warn, or Info
	Message    string     // Message is a human-readable description of the audit result
	PendingFix PendingFix // PendingFix is the fix that will be applied to automatically fix the security issue
	Metadata   Metadata   // Metadata includes additional context for an audit result
}

func (result *AuditResult) Fix(resource k8stypes.Resource) (newResources []k8stypes.Resource) {
	if result.PendingFix == nil {
		return nil
	}

	return result.PendingFix.Apply(resource)
}

func (result *AuditResult) FixPlan() (ok bool, plan string) {
	if result.PendingFix == nil {
		return false, ""
	}

	return true, result.PendingFix.Plan()
}

// PendingFix includes the logic to automatically fix the issues caught by auditing
type PendingFix interface {
	// Plan returns a human-readable description of what Apply() will do
	Plan() string
	// Apply applies the proposed fix to the resource and returns any new resources that were created. Note that
	// Apply is expected to modify the passed in resource
	Apply(k8stypes.Resource) []k8stypes.Resource
}

// Metadata holds metadata for a potential security issue
type Metadata = map[string]string

// Implements Result
type workloadResult struct {
	Resource     KubeResource
	AuditResults []*AuditResult
}

func (wlResult *workloadResult) GetResource() KubeResource {
	return wlResult.Resource
}

func (wlResult *workloadResult) GetAuditResults() []*AuditResult {
	return wlResult.AuditResults
}
