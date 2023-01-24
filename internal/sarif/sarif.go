package sarif

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

const repoURL = "https://github.com/Shopify/kubeaudit"

// Create generates new sarif Report or returns an error
func Create(kubeauditReport *kubeaudit.Report) (*sarif.Report, error) {
	// create a new report object
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return nil, err
	}

	// create a run for kubeaudit
	run := sarif.NewRunWithInformationURI("kubeaudit", repoURL)

	report.AddRun(run)

	var results []*kubeaudit.AuditResult

	for _, reportResult := range kubeauditReport.Results() {
		r := reportResult.GetAuditResults()
		results = append(results, r...)
	}

	for _, result := range results {
		severityLevel := result.Severity.String()

		auditor := strings.ToLower(result.Auditor)

		var metadataTxt string
		if len(result.Metadata) > 0 {
			formattedMap := make(map[string]string)

			for k, v := range result.Metadata {
				formattedMap[k] = v
			}

			metadata, jsonErr := json.Marshal(formattedMap)
			if jsonErr != nil {
				metadata = []byte(jsonErr.Error())
			}

			metadataTxt = fmt.Sprintf("Metadata: %s\n", string(metadata))
		}

		docsURL := "https://github.com/Shopify/kubeaudit/blob/main/docs/auditors/" + auditor + ".md"

		helpText := fmt.Sprintf("Type: kubernetes\nAuditor Docs: To find out more about the issue and how to fix it, follow [this link](%s)\nDescription: %s\n%s\n\n Note: These audit results are generated with `kubeaudit`, a command line tool and a Go package that checks for potential security concerns in kubernetes manifest specs. You can read more about it at https://github.com/Shopify/kubeaudit ", docsURL, allAuditors[auditor], metadataTxt)

		helpMarkdown := fmt.Sprintf("**Type**: kubernetes\n**Auditor Docs**: To find out more about the issue and how to fix it, follow [this link](%s)\n**Description:** %s\n **Metadata**: %s\n\n *Note*: These audit results are generated with `kubeaudit`, a command line tool and a Go package that checks for potential security concerns in kubernetes manifest specs. You can read more about it at https://github.com/Shopify/kubeaudit ",
			docsURL, allAuditors[auditor], metadataTxt)

		// we only add rules to the report based on the result findings
		run.AddRule(result.Rule).
			WithName(result.Auditor).
			WithHelp(&sarif.MultiformatMessageString{Text: &helpText, Markdown: &helpMarkdown}).
			WithShortDescription(&sarif.MultiformatMessageString{Text: &result.Rule}).
			WithProperties(sarif.Properties{
				"tags": []string{
					"security",
					"kubernetes",
					"infrastructure",
				},
			})

		// SARIF specifies the following severity levels: warning, error, note and none
		// https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-v2.1.0.html
		// so we're converting info to note here so we get valid SARIF output
		if result.Severity.String() == kubeaudit.Info.String() {
			severityLevel = "note"
		}

		details := fmt.Sprintf("Details: %s\n Auditor: %s\nDescription: %s\nAuditor docs: %s ",
			result.Message, result.Auditor, allAuditors[auditor], docsURL)

		location := sarif.NewPhysicalLocation().
			WithArtifactLocation(sarif.NewSimpleArtifactLocation(result.FilePath).WithUriBaseId("ROOTPATH")).
			WithRegion(sarif.NewRegion().WithStartLine(1))
		result := sarif.NewRuleResult(result.Rule).
			WithMessage(sarif.NewTextMessage(details)).
			WithLevel(severityLevel).
			WithLocations([]*sarif.Location{sarif.NewLocation().WithPhysicalLocation(location)})
		run.AddResult(result)
	}

	var reportBytes bytes.Buffer

	err = report.Write(&reportBytes)
	if err != nil {
		return nil, nil
	}

	return report, nil
}
