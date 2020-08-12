package kubeaudit

import (
	"fmt"
	"io"
	"os"

	"github.com/Shopify/kubeaudit/internal/color"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
	log "github.com/sirupsen/logrus"
)

type Printer struct {
	writer      io.Writer
	minSeverity SeverityLevel
	formatter   log.Formatter
	color       bool
}

type PrintOption func(p *Printer)

func WithMinSeverity(minSeverity SeverityLevel) PrintOption {
	return func(p *Printer) {
		p.minSeverity = minSeverity
	}
}

func WithWriter(writer io.Writer) PrintOption {
	return func(p *Printer) {
		p.writer = writer
	}
}

func WithFormatter(formatter log.Formatter) PrintOption {
	return func(p *Printer) {
		p.formatter = formatter
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
	}
	p.parseOptions(opts...)
	if p.writer == os.Stdout {
		p.color = true
	}
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
	for _, workloadResult := range report.ResultsWithMinSeverity(p.minSeverity) {
		resource := workloadResult.GetResource().Object()
		groupVersionKind := resource.GetObjectKind().GroupVersionKind()
		resourceName := k8s.GetObjectMeta(resource).GetName()
		resourceNamespace := k8s.GetObjectMeta(resource).GetNamespace()

		p.printColor(color.CyanColor, "--------- Results for ---------------------\n\n")
		p.printColor(color.CyanColor, "  apiVersion: ")
		if groupVersionKind.Group != "" {
			p.printColor(color.CyanColor, groupVersionKind.Group+"/")
		}
		p.printColor(color.CyanColor, groupVersionKind.Version+"\n")
		p.printColor(color.CyanColor, ("  kind: " + groupVersionKind.Kind + "\n"))
		if resourceName != "" || resourceNamespace != "" {
			p.printColor(color.CyanColor, "  metadata:\n")
			if resourceName != "" {
				p.printColor(color.CyanColor, "    name: "+resourceName+"\n")
			}
			if resourceNamespace != "" {
				p.printColor(color.CyanColor, "    namespace: "+resourceNamespace+"\n")
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
		p.print("\n")
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

	// We manually manage what severity levels to log, lorgus should let everything through
	resultLogger.SetLevel(log.DebugLevel)

	for _, workloadResult := range report.ResultsWithMinSeverity(p.minSeverity) {
		for _, auditResult := range workloadResult.GetAuditResults() {
			p.logAuditResult(workloadResult.GetResource().Object(), auditResult, resultLogger)
		}
	}
}

func (p *Printer) logAuditResult(resource k8stypes.Resource, result *AuditResult, baseLogger *log.Logger) {
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

func (p *Printer) getLogFieldsForResult(resource k8stypes.Resource, result *AuditResult) log.Fields {
	groupVersionKind := resource.GetObjectKind().GroupVersionKind()
	resourceMetadata := k8s.GetObjectMeta(resource)

	fields := log.Fields{
		"AuditResultName":   result.Name,
		"ResourceVersion":   groupVersionKind.Version,
		"ResourceKind":      groupVersionKind.Kind,
		"ResourceGroup":     groupVersionKind.Group,
		"ResourceNamespace": resourceMetadata.GetNamespace(),
	}

	if fields["ResourceGroup"] == "" {
		fields["ResourceGroup"] = "core"
	}

	if resourceMetadata.GetName() != "" {
		fields["ResourceName"] = resourceMetadata.GetName()
	}

	for k, v := range result.Metadata {
		fields[k] = v
	}

	return fields
}
