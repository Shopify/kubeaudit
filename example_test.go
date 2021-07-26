package kubeaudit_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/config"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Example shows how to audit and fix a Kubernetes manifest file
func Example() {
	// A sample Kubernetes manifest file
	manifest := `
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

	// Initialize all the security auditors using default configuration
	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize kubeaudit
	auditor, err := kubeaudit.New(allAuditors)
	if err != nil {
		log.Fatal(err)
	}

	// Run the audit in manifest mode
	report, err := auditor.AuditManifest(strings.NewReader(manifest))
	if err != nil {
		log.Fatal(err)
	}

	// Print the audit results to screen
	report.PrintResults()

	// Print the plan to screen. These are the steps that will be taken by calling "report.Fix()".
	fmt.Println("\nPlan:")
	report.PrintPlan(os.Stdout)

	// Print the fixed manifest to screen. Note that this leaves the original manifest unmodified.
	fmt.Println("\nFixed manifest:")
	err = report.Fix(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

// ExampleAuditLocal shows how to run kubeaudit in local mode
func Example_auditLocal() {
	// Initialize all the security auditors using default configuration
	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	if err != nil {
		log.WithError(err).Fatal("Error initializing all auditors")
	}

	// Initialize kubeaudit
	auditor, err := kubeaudit.New(allAuditors)
	if err != nil {
		log.Fatal(err)
	}

	// Run the audit in local mode
	report, err := auditor.AuditLocal("", kubeaudit.AuditOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Print the audit results to screen
	report.PrintResults()
}

// ExampleAuditCluster shows how to run kubeaudit in cluster mode (only works if kubeaudit is being run from a container insdie of a cluster)
func Example_auditCluster() {
	// Initialize all the security auditors using default configuration
	allAuditors, err := all.Auditors(config.KubeauditConfig{})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize kubeaudit
	auditor, err := kubeaudit.New(allAuditors)
	if err != nil {
		log.Fatal(err)
	}

	// Run the audit in cluster mode. Note this will fail if kubeaudit is not running within a cluster.
	report, err := auditor.AuditCluster(kubeaudit.AuditOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Print the audit results to screen
	report.PrintResults()
}

// ExampleAuditorSubset shows how to run kubeaudit with a subset of auditors
func Example_auditorSubset() {
	// A sample Kubernetes manifest file
	manifest := `
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

	// Initialize the auditors you want to use
	auditor, err := kubeaudit.New([]kubeaudit.Auditable{
		apparmor.New(),
		image.New(image.Config{Image: "myimage:mytag"}),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Run the audit in the mode of your choosing. Here we use manifest mode.
	report, err := auditor.AuditManifest(strings.NewReader(manifest))
	if err != nil {
		log.Fatal(err)
	}

	// Print the audit results to screen
	report.PrintResults()
}

// ExampleConfig shows how to use a kubeaudit with a config file.
// A kubeaudit config can be used to specify which security auditors to run, and to specify configuration
// for those auditors.
func Example_config() {
	configFile := "config/config.yaml"
	// A sample Kubernetes manifest file
	manifest := `
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

	// Open the configuration file
	reader, err := os.Open(configFile)
	if err != nil {
		log.WithError(err).Fatal("Unable to open config file ", configFile)
	}

	// Load the config
	conf, err := config.New(reader)
	if err != nil {
		log.WithError(err).Fatal("Error parsing config file ", configFile)
	}

	// Initialize security auditors using the configuration
	auditors, err := all.Auditors(conf)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize kubeaudit
	auditor, err := kubeaudit.New(auditors)
	if err != nil {
		log.Fatal(err)
	}

	// Run the audit in the mode of your choosing. Here we use manifest mode.
	report, err := auditor.AuditManifest(strings.NewReader(manifest))
	if err != nil {
		log.Fatal(err)
	}

	// Print the audit results to screen
	report.PrintResults()
}

// ExamplePrintOptions shows how to use different print options for printing audit results.
func Example_printOptions() {
	auditor, err := kubeaudit.New([]kubeaudit.Auditable{apparmor.New()})
	if err != nil {
		log.Fatal(err)
	}

	report, err := auditor.AuditLocal("", kubeaudit.AuditOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Print the audit results to a file
	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	defer os.Remove("output.txt")
	report.PrintResults(kubeaudit.WithWriter(f))

	// Only print audit results with severity of Error (ignore info and warning)
	report.PrintResults(kubeaudit.WithMinSeverity(kubeaudit.Error))

	// Print results as JSON
	report.PrintResults(kubeaudit.WithFormatter(&logrus.JSONFormatter{}))
}
