package kubeaudit

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

type logEntry struct {
	AuditResultName string
	Foo             string
	Level           string `json:"level"`
}

var levelString = map[int]string{
	Error: "error",
	Warn:  "warning",
	Info:  "info",
}

func TestGetResourcesFromClientset(t *testing.T) {
	resources := []runtime.Object{k8stypes.NewDeployment(), k8stypes.NewNamespace()}

	expected := make([]KubeResource, 0, len(resources))
	for _, resource := range resources {
		expected = append(expected, &kubeResource{object: resource})
	}

	got := getResourcesFromClientset(fakeclientset.NewSimpleClientset(resources...), k8s.ClientOptions{})
	assert.Equal(t, expected, got)
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
			},
		},
	}
	out := bytes.NewBuffer(nil)

	// Error results only
	report.PrintResults(out, Error, &log.JSONFormatter{})
	assert.Equal(t, 1, bytes.Count(out.Bytes(), []byte{'\n'}))
	out.Reset()

	// Error and warn results
	report.PrintResults(out, Warn, &log.JSONFormatter{})
	assert.Equal(t, 2, bytes.Count(out.Bytes(), []byte{'\n'}))
	out.Reset()

	// Error, warn, and info results
	report.PrintResults(out, Info, &log.JSONFormatter{})
	assert.Equal(t, 3, bytes.Count(out.Bytes(), []byte{'\n'}))
}

func newTestAuditResult(severity int) *AuditResult {
	return &AuditResult{
		Name:     "MyAuditResult",
		Severity: severity,
		Metadata: Metadata{"Foo": "bar"},
	}
}

func TestLogAuditResult(t *testing.T) {
	// Send all log output as JSON to this byte buffer
	out := bytes.NewBuffer(nil)
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(out)

	for _, severity := range []int{Error, Warn, Info} {
		auditResult := newTestAuditResult(severity)
		expected := logEntry{
			AuditResultName: "MyAuditResult",
			Level:           levelString[severity],
			Foo:             auditResult.Metadata["Foo"],
		}

		// This writes the log to the variable out, parses the JSON into the logEntry struct, and checks the struct
		logAuditResult(auditResult, logger)
		got := logEntry{}
		err := json.Unmarshal(out.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		out.Reset()
	}
}
