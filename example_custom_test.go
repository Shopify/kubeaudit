package kubeaudit_test

import (
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"

	log "github.com/sirupsen/logrus"
)

func NewCustomAuditor() kubeaudit.Auditable {
	return &myAuditor{}
}

// Your auditor must implement the Auditable interface, which requires only one method: Audit().
type myAuditor struct{}

// The Audit function takes in a resource to audit and returns audit results for that resource.
//
// Params
//   resource: Read-only. The resource to audit.
//   resources: Read-only. A reference to all resources. Can be used for context though most auditors don't need this.
//
// Return
//   auditResults: The results for the audit. Each result can optionally include a PendingFix object to
//     define autofix behaviour (see below).
func (a *myAuditor) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	return []*kubeaudit.AuditResult{
		{
			Name:     "MyAudit",
			Severity: kubeaudit.Error,
			Message:  "My custom error",
			PendingFix: &myAuditorFix{
				newVal: "bye",
			},
		},
	}, nil
}

// To provide autofix behaviour for an audit result, implement the PendingFix interface. The PendingFix interface
// has two methods: Plan() and Apply().
type myAuditorFix struct {
	newVal string
}

// The Plan method explains what fix will be applied by Apply().
//
// Return
//   plan: A human-friendly explanation of what Apply() will do
func (f *myAuditorFix) Plan() string {
	return fmt.Sprintf("Set label 'hi' to '%s'", f.newVal)
}

// The Apply method applies a fix to a resource.
//
// Params
//   resource: A reference to the resource that should be fixed.
//
// Return
//   newResources: New resources created as part of the fix. Generally, it should not be necessary to create
//     new resources, only modify the passed in resource.
func (f *myAuditorFix) Apply(resource k8s.Resource) []k8s.Resource {
	setLabel(resource, "hi", f.newVal)
	return nil
}

// This is just a helper function
func setLabel(resource k8s.Resource, key, value string) {
	switch kubeType := resource.(type) {
	case *k8s.PodV1:
		kubeType.Labels[key] = value
	case *k8s.DeploymentV1:
		kubeType.Labels[key] = value
	}
}

// A sample Kubernetes manifest file
var manifest = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myAuditor 
spec:
  template:
    spec:
      containers:
      - name: myContainer
`

// ExampleCustomAuditor shows how to use a custom auditor
func Example_customAuditor() {
	// Initialize kubeaudit with your custom auditor
	auditor, err := kubeaudit.New([]kubeaudit.Auditable{NewCustomAuditor()})
	if err != nil {
		log.Fatal(err)
	}

	// Run the audit in the mode of your choosing. Here we use manifest mode.
	report, err := auditor.AuditManifest(strings.NewReader(manifest))
	if err != nil {
		log.Fatal(err)
	}

	// Print the results to screen
	report.PrintResults()
}
