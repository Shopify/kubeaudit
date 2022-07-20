package sarif

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/owenrumney/go-sarif/v2/sarif"
	"github.com/xeipuuv/gojsonschema"
)

// Create generates new sarif Report or returns an error
func Create(kubeauditReport *kubeaudit.Report) (*sarif.Report, error) {
	// create a new report object
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return nil, err
	}

	// create a run for kubeaudit
	run := sarif.NewRunWithInformationURI("kubeaudit", "https://github.com/Shopify/kubeaudit")

	report.AddRun(run)

	var results []*kubeaudit.AuditResult

	for _, reportResult := range kubeauditReport.Results() {
		r := reportResult.GetAuditResults()
		results = append(results, r...)
	}

	for _, result := range results {
		severityLevel := result.Severity.String()
		auditor := strings.ToLower(result.Auditor)

		var docsURL string

		auditor, ok := violationsToRules[result.Rule]
		if ok {
			docsURL = "https://github.com/Shopify/kubeaudit/blob/main/docs/auditors/" + auditor + ".md"
		}

		helpMessage := fmt.Sprintf("**Type**: kubernetes\n**Docs**: %s\n**Description:** %s", docsURL, allAuditors[auditor])

		// we only add rules to the report based on the result findings
		run.AddRule(result.Rule).
			WithName(result.Auditor).
			WithMarkdownHelp(helpMessage).
			WithHelp(&sarif.MultiformatMessageString{Text: &docsURL}).
			WithShortDescription(&sarif.MultiformatMessageString{Text: &result.Rule}).
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

	var reportBytes bytes.Buffer

	err = report.Write(&reportBytes)
	if err != nil {
		return nil, nil
	}

	err, errs := validate(&reportBytes)
	if err != nil {
		return nil, fmt.Errorf("error validating SARIF schema: %s", err)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("SARIF schema validation errors: %s", errs)
	}

	return report, nil
}

// Validates that the SARIF file is valid as per sarif spec
// https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Documents/CommitteeSpecifications/2.1.0/sarif-schema-2.1.0.json
func validate(report io.Reader) (error, []gojsonschema.ResultError) {
	schemaLoader := gojsonschema.NewReferenceLoader("http://json.schemastore.org/sarif-2.1.0")
	var reportLoader gojsonschema.JSONLoader

	_, ok := report.(*bytes.Buffer)
	if ok {
		reportLoader = gojsonschema.NewStringLoader(report.(*bytes.Buffer).String())
	}

	result, err := gojsonschema.Validate(schemaLoader, reportLoader)

	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		return nil, nil
	} else {
		return nil, result.Errors()
	}
}
