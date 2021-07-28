package kubeaudit

import (
	"fmt"
	"io"
	"os"

	"github.com/Shopify/kubeaudit/internal/color"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	log "github.com/sirupsen/logrus"
)

type Printer struct {
	writer      io.Writer
	minSeverity SeverityLevel
	formatter   log.Formatter
	color       bool
}

type PrintOption func(p *Printer)

// WithMinSeverity sets the minimum severity of results that will be printed.
func WithMinSeverity(minSeverity SeverityLevel) PrintOption {
	return func(p *Printer) {
		p.minSeverity = minSeverity
	}
}

// WithWriter sets the writer where results will be written to.
func WithWriter(writer io.Writer) PrintOption {
	return func(p *Printer) {
		p.writer = writer
	}
}

// WithFormatter sets a logrus formatter to use to format results.
func WithFormatter(formatter log.Formatter) PrintOption {
	return func(p *Printer) {
		p.formatter = formatter
	}
}

// WithColor specifies whether or not to colorize output. You will likely want to set this to false if
// not writing to standard out.
func WithColor(color bool) PrintOption {
	return func(p *Printer) {
		p.color = color
	}
}

func (p *Printer) parseOptions(opts ...PrintOption) {
	for _, opt := range opts {
		opt(p)
	}
}

func NewPrinter(opts ...PrintOption) Printer {
	p := Printer{
		writer:      os.Stdout,
		minSeverity: Info,
		color:       true,
	}
	p.parseOptions(opts...)
	return p
}

func (p *Printer) PrintReport(report *Report) {
	if p.formatter == nil {
		p.prettyPrintReport(report)
	} else {
		p.logReport(report)
	}
}

func (p *Printer) prettyPrintReport(report *Report) {
	if len(report.ResultsWithMinSeverity(p.minSeverity)) < 1 {
		p.printColor(color.GreenColor, "All checks completed. 0 high-risk vulnerabilities found\n")
		return
	}

	for _, workloadResult := range report.ResultsWithMinSeverity(p.minSeverity) {
		resource := workloadResult.GetResource().Object()
		objectMeta := k8s.GetObjectMeta(resource)
		resouceApiVersion, resourceKind := resource.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()

		p.printColor(color.CyanColor, "\n---------------- Results for ---------------\n\n")
		p.printColor(color.CyanColor, "  apiVersion: "+resouceApiVersion+"\n")
		p.printColor(color.CyanColor, "  kind: "+resourceKind+"\n")
		if objectMeta != nil && (objectMeta.GetName() != "" || objectMeta.GetNamespace() != "") {
			p.printColor(color.CyanColor, "  metadata:\n")
			if objectMeta.GetName() != "" {
				p.printColor(color.CyanColor, "    name: "+objectMeta.GetName()+"\n")
			}
			if objectMeta.GetNamespace() != "" {
				p.printColor(color.CyanColor, "    namespace: "+objectMeta.GetNamespace()+"\n")
			}
		}
		p.printColor(color.CyanColor, "\n--------------------------------------------\n\n")

		for _, auditResult := range workloadResult.GetAuditResults() {
			severityColor := color.YellowColor
			switch auditResult.Severity {
			case Info:
				severityColor = color.CyanColor
			case Warn:
				severityColor = color.YellowColor
			case Error:
				severityColor = color.RedColor
			}
			p.print("-- ")
			p.printColor(severityColor, "["+auditResult.Severity.String()+"] ")
			p.print(auditResult.Name + "\n")
			p.print("   Message: " + auditResult.Message + "\n")
			if len(auditResult.Metadata) > 0 {
				p.print("   Metadata:\n")
			}
			for k, v := range auditResult.Metadata {
				p.print(fmt.Sprintf("      %s: %s\n", k, v))
			}
			p.print("\n")
		}
	}
}

func (p *Printer) print(s string) {
	fmt.Fprint(p.writer, s)
}

func (p *Printer) printColor(c string, s string) {
	if p.color {
		fmt.Fprint(p.writer, color.Colored(c, s))
	} else {
		p.print(s)
	}
}

func (p *Printer) logReport(report *Report) {
	resultLogger := log.New()
	resultLogger.SetOutput(p.writer)
	resultLogger.SetFormatter(p.formatter)

	// We manually manage what severity levels to log, logrus should let everything through
	resultLogger.SetLevel(log.DebugLevel)

	for _, workloadResult := range report.ResultsWithMinSeverity(p.minSeverity) {
		for _, auditResult := range workloadResult.GetAuditResults() {
			p.logAuditResult(workloadResult.GetResource().Object(), auditResult, resultLogger)
		}
	}
}

func (p *Printer) logAuditResult(resource k8s.Resource, result *AuditResult, baseLogger *log.Logger) {
	logger := baseLogger.WithFields(p.getLogFieldsForResult(resource, result))
	switch result.Severity {
	case Info:
		logger.Info(result.Message)
	case Warn:
		logger.Warn(result.Message)
	case Error:
		logger.Error(result.Message)
	}
}

func (p *Printer) getLogFieldsForResult(resource k8s.Resource, result *AuditResult) log.Fields {
	apiVersion, kind := resource.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
	objectMeta := k8s.GetObjectMeta(resource)

	fields := log.Fields{
		"AuditResultName":    result.Name,
		"ResourceKind":       kind,
		"ResourceApiVersion": apiVersion,
	}

	if objectMeta != nil {
		if objectMeta.GetNamespace() != "" {
			fields["ResourceNamespace"] = objectMeta.GetNamespace()
		}

		if objectMeta.GetName() != "" {
			fields["ResourceName"] = objectMeta.GetName()
		}
	}

	for k, v := range result.Metadata {
		fields[k] = v
	}

	return fields
}
