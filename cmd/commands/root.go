package commands

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/config"
	"github.com/Shopify/kubeaudit/internal/k8s"
)

var rootConfig rootFlags

type rootFlags struct {
	format      string
	kubeConfig  string
	manifest    string
	namespace   string
	minSeverity string
	exitCode    int
}

// RootCmd defines the shell command usage for kubeaudit.
var RootCmd = &cobra.Command{
	Use:   "kubeaudit",
	Short: "A Kubernetes security auditor",
	Long: `Kubeaudit audits Kubernetes clusters for common security controls.

kubeaudit has three modes:
  1. Manifest mode: If a Kubernetes manifest file is provided using the -f/--manifest flag, kubeaudit will audit the manifest file. Kubeaudit also supports autofixing in manifest mode using the 'autofix' command. This will fix the manifest in-place. The fixed manifest can be written to a different file using the -o/--out flag.
  2. Cluster mode: If kubeaudit detects it is running in a cluster, it will audit the other resources in the cluster.
  3. Local mode: kubeaudit will try to connect to a cluster using the local kubeconfig file ($HOME/.kube/config). A different kubeconfig location can be specified using the -c/--kubeconfig flag
`,
}

// Execute is a wrapper for the RootCmd.Execute method which will exit the program if there is an error.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&rootConfig.kubeConfig, "kubeconfig", "c", "", "Path to local Kubernetes config file. Only used in local mode (default is $HOME/.kube/config)")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.minSeverity, "minseverity", "m", "info", "Set the lowest severity level to report (one of \"error\", \"warning\", \"info\")")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.format, "format", "p", "pretty", "The output format to use (one of \"pretty\", \"logrus\", \"json\")")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.namespace, "namespace", "n", apiv1.NamespaceAll, "Only audit resources in the specified namespace. Not currently supported in manifest mode.")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.manifest, "manifest", "f", "", "Path to the yaml configuration to audit. Only used in manifest mode.")
	RootCmd.PersistentFlags().IntVarP(&rootConfig.exitCode, "exitcode", "e", 2, "Exit code to use if there are results with severity of \"error\". Conventionally, 0 is used for success and all non-zero codes for an error.")
}

// KubeauditLogLevels represents an enum for the supported log levels.
var KubeauditLogLevels = map[string]kubeaudit.SeverityLevel{
	"error":   kubeaudit.Error,
	"warn":    kubeaudit.Warn,
	"warning": kubeaudit.Warn,
	"info":    kubeaudit.Info,
}

func runAudit(auditable ...kubeaudit.Auditable) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		report := getReport(auditable...)

		printOptions := []kubeaudit.PrintOption{
			kubeaudit.WithMinSeverity(KubeauditLogLevels[strings.ToLower(rootConfig.minSeverity)]),
		}

		switch rootConfig.format {
		case "json":
			printOptions = append(printOptions, kubeaudit.WithFormatter(&log.JSONFormatter{}))
		case "logrus":
			printOptions = append(printOptions, kubeaudit.WithFormatter(&log.TextFormatter{}))
		}

		report.PrintResults(printOptions...)

		if report.HasErrors() {
			os.Exit(rootConfig.exitCode)
		}
	}
}

func getReport(auditors ...kubeaudit.Auditable) *kubeaudit.Report {
	auditor := initKubeaudit(auditors...)

	if rootConfig.manifest != "" {
		manifest, err := os.Open(rootConfig.manifest)
		if err != nil {
			log.WithError(err).Fatal("Error opening manifest file")
		}
		report, err := auditor.AuditManifest(manifest)
		if err != nil {
			log.WithError(err).Fatal("Error auditing manifest")
		}
		return report
	}

	if ( k8s.IsRunningInCluster(k8s.DefaultClient) && rootConfig.kubeConfig == "" ) {
		report, err := auditor.AuditCluster(k8s.ClientOptions{Namespace: rootConfig.namespace})
		if err != nil {
			log.WithError(err).Fatal("Error auditing cluster")
		}
		return report
	}

	report, err := auditor.AuditLocal(rootConfig.kubeConfig, k8s.ClientOptions{Namespace: rootConfig.namespace})
	if err != nil {
		log.WithError(err).Fatal("Error auditing cluster in local mode")
	}
	return report
}

func initKubeaudit(auditable ...kubeaudit.Auditable) *kubeaudit.Kubeaudit {
	if len(auditable) == 0 {
		allAuditors, err := all.Auditors(config.KubeauditConfig{})
		if err != nil {
			log.WithError(err).Fatal("Error initializing auditors")
		}
		auditable = allAuditors
	}

	auditor, err := kubeaudit.New(auditable)
	if err != nil {
		log.WithError(err).Fatal("Error creating auditor")
	}

	return auditor
}
