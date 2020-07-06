package commands

import (
	"os"

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
	json        bool
	kubeConfig  string
	manifest    string
	namespace   string
	minSeverity string
}

// RootCmd defines the shell command usage for kubeaudit.
var RootCmd = &cobra.Command{
	Use:   "kubeaudit",
	Short: "A Kubernetes security auditor",
	Long: `kubeaudit is a program that makes sure all your containers are secure #patcheswelcome

kubeaudit has three modes:
1. Manifest mode: If a Kubernetes manifest file is provided using the -f/--manifest flag, kubeaudit will audit the manifest file.
	 Kubeaudit also supports autofixing in manifest mode using the 'autofix' command. This will fix the manifest in-place.
	 The fixed manfiest can be written to a different file using the -o/--out flag.
2. Cluster mode: If kubeaudit detects it is running within a container, it will try to audit the cluster it is contained in.
3. Local mode: kubeaudit will audit the resources specified by the local kubeconfig file ($HOME/.kube/config). A different
     kubeaconfig location can be specified using the -c/--kubeconfig flag
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
	RootCmd.PersistentFlags().StringVarP(&rootConfig.minSeverity, "minseverity", "m", "INFO", "Set the lowest severity level to report (one of \"ERROR\", \"WARN\", \"INFO\")")
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.json, "json", "j", false, "Output audit results in JSON")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.namespace, "namespace", "n", apiv1.NamespaceAll, "Only audit resources in the specified namespace. Only used in cluster mode.")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.manifest, "manifest", "f", "", "Path to the yaml configuration to audit. Only used in manifest mode.")
}

// KubeauditLogLevels represents an enum for the supported log levels.
var KubeauditLogLevels = map[string]int{
	"ERROR": kubeaudit.Error,
	"WARN":  kubeaudit.Warn,
	"INFO":  kubeaudit.Info,
}

func runAudit(auditable ...kubeaudit.Auditable) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		report := getReport(auditable...)

		minSeverity := KubeauditLogLevels[rootConfig.minSeverity]

		var formatter log.Formatter
		if rootConfig.json {
			formatter = &log.JSONFormatter{}
		}

		report.PrintResults(os.Stdout, minSeverity, formatter)
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

	if k8s.IsRunningInCluster(k8s.DefaultClient) {
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
