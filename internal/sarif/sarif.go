package sarif

import (
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/Shopify/kubeaudit/auditors/asat"
	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/Shopify/kubeaudit/auditors/hostns"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/auditors/limits"
	"github.com/Shopify/kubeaudit/auditors/mounts"
	"github.com/Shopify/kubeaudit/auditors/netpols"
	"github.com/Shopify/kubeaudit/auditors/nonroot"
	"github.com/Shopify/kubeaudit/auditors/privesc"
	"github.com/Shopify/kubeaudit/auditors/privileged"
	"github.com/Shopify/kubeaudit/auditors/rootfs"
	"github.com/Shopify/kubeaudit/auditors/seccomp"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

var Auditors = map[string]string{
	apparmor.Name:     "Finds containers that do not have AppArmor enabled",
	asat.Name:         "Finds containers where the deprecated SA field is used or with a mounted default SA",
	capabilities.Name: "Finds containers that do not drop the recommended capabilities or add new ones",
	hostns.Name:       "Finds containers that have HostPID, HostIPC or HostNetwork enabled",
	image.Name:        "Finds containers which do not use the desired version of an image (via the tag) or use an image without a tag",
	limits.Name:       "Finds containers which exceed the specified CPU and memory limits or do not specify any",
	mounts.Name:       "Finds containers that have sensitive host paths mounted",
	netpols.Name:      "Finds namespaces that do not have a default-deny network policy",
	nonroot.Name:      "Finds containers allowed to run as root",
	privesc.Name:      "Finds containers that allow privilege escalation",
	privileged.Name:   "Finds containers running as privileged",
	rootfs.Name:       "Finds containers which do not have a read-only filesystem",
	seccomp.Name:      "Finds containers running without seccomp",
}

// New creates a new sarif Report and Run or returns an error
func New() (*sarif.Report, *sarif.Run, error) {
	// create a new report object
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return nil, nil, err
	}

	// create a run for kubeaudit
	run := sarif.NewRunWithInformationURI("kubeaudit", "https://github.com/Shopify/kubeaudit")

	report.AddRun(run)

	return report, run, nil
}

// Create adds SARIF rules to the run and adds results to the report
func Create(kubeauditReport *kubeaudit.Report, run *sarif.Run) {
	var results []*kubeaudit.AuditResult

	for _, reportResult := range kubeauditReport.Results() {
		r := reportResult.GetAuditResults()
		results = append(results, r...)
	}

	for _, result := range results {
		severityLevel := result.Severity.String()
		auditor := strings.ToLower(result.Auditor)
		ruleID := strings.ToLower(result.Rule)

		var docsURL string
		if strings.Contains(ruleID, auditor) {
			docsURL = "https://github.com/Shopify/kubeaudit/blob/main/docs/auditors/" + auditor + ".md"
		}

		helpMessage := fmt.Sprintf("**Type**: kubernetes\n**Docs**: %s\n**Description:** %s", docsURL, Auditors[auditor])

		// we only add rules to the report based on the result findings
		run.AddRule(result.Rule).
			WithName(result.Auditor).
			WithMarkdownHelp(helpMessage).
			WithProperties(sarif.Properties{
				"tags": []string{
					"security",
					"kubernetes",
					"infrastructure",
				},
				"precision": "very-high",
			})

		// SARIF specifies the following severity levels: warning, error, note and none
		// https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-v2.1.0.html
		// so we're converting info to none here so we get valid SARIF output
		if result.Severity.String() == "info" {
			severityLevel = "note"
		}

		location := sarif.NewPhysicalLocation().
			WithArtifactLocation(sarif.NewSimpleArtifactLocation(result.FilePath).WithUriBaseId("ROOTPATH")).
			WithRegion(sarif.NewRegion().WithStartLine(1))
		result := sarif.NewRuleResult(result.Rule).
			WithMessage(sarif.NewTextMessage(result.Message)).
			WithLevel(severityLevel).
			WithLocations([]*sarif.Location{sarif.NewLocation().WithPhysicalLocation(location)})
		run.AddResult(result)
	}
}
