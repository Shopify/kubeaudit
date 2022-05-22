package kubeaudit

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Shopify/kubeaudit/pkg/k8s"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type logEntry struct {
	AuditResultName    string
	Foo                string
	Level              string `json:"level"`
	ResourceKind       string
	ResourceApiVersion string
	ResourceName       string
	ResourceNamespace  string
}

func TestPrintResults(t *testing.T) {
	report := Report{
		results: []Result{
			&workloadResult{
				AuditResults: []*AuditResult{
					newTestAuditResult(Error),
					newTestAuditResult(Warn),
					newTestAuditResult(Info),
				},
				Resource: &kubeResource{
					object: k8s.NewPod(),
				},
			},
		},
	}
	out := bytes.NewBuffer(nil)
	writerOption := WithWriter(out)
	formatterOption := WithFormatter(&log.JSONFormatter{})

	// Error results only
	report.PrintResults(writerOption, WithMinSeverity(Error), formatterOption)
	assert.Equal(t, 1, bytes.Count(out.Bytes(), []byte{'\n'}))
	out.Reset()

	// Error and warn results
	report.PrintResults(writerOption, WithMinSeverity(Warn), formatterOption)
	assert.Equal(t, 2, bytes.Count(out.Bytes(), []byte{'\n'}))
	out.Reset()

	// Error, warn, and info results
	report.PrintResults(writerOption, WithMinSeverity(Info), formatterOption)
	assert.Equal(t, 3, bytes.Count(out.Bytes(), []byte{'\n'}))
}

func newTestAuditResult(severity SeverityLevel) *AuditResult {
	return &AuditResult{
		Name:     "MyAuditResult",
		Severity: severity,
		Metadata: Metadata{"Foo": "bar"},
	}
}

func TestLogAuditResult(t *testing.T) {
	for _, severity := range []SeverityLevel{Error, Warn, Info} {
		// Send all log output as JSON to this byte buffer
		out := bytes.NewBuffer(nil)

		resource := k8s.NewDeployment()
		resource.Name = "mydeployment"
		resource.Namespace = "mynamespace"

		auditResult := newTestAuditResult(severity)
		report := &Report{
			results: []Result{
				&workloadResult{
					AuditResults: []*AuditResult{
						auditResult,
					},
					Resource: &kubeResource{
						object: resource,
					},
				},
			},
		}
		expectedApiVersion, expectedKind := resource.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
		expected := logEntry{
			AuditResultName:    "MyAuditResult",
			Level:              severity.String(),
			Foo:                auditResult.Metadata["Foo"],
			ResourceKind:       expectedKind,
			ResourceApiVersion: expectedApiVersion,
			ResourceName:       resource.GetName(),
			ResourceNamespace:  resource.GetNamespace(),
		}

		// This writes the log to the variable out, parses the JSON into the logEntry struct, and checks the struct
		printer := NewPrinter(WithWriter(out), WithFormatter(&log.JSONFormatter{}))
		printer.PrintReport(report)
		got := logEntry{}
		err := json.Unmarshal(out.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		out.Reset()
	}
}
