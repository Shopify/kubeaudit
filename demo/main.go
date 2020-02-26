// Run the examples with
// $ make run-demo

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	demoManifest()
	demoLocal()
	demoCluster()
	demoAuditorSubset()
	demoAuditWithConfig()
	demoCustomAuditor()
}

var manifest = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myAuditor 
  namespace: myNamespace
spec:
  strategy: {}
  template:
    metadata:
      labels:
        apps: myLabel
    spec:
      containers:
      - name: myContainer
        resources: {}
status: {}
`

// Audit a manifest file
func demoManifest() {
	fmt.Println("\n--- Manifest ---")

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
	report.PrintResults(os.Stdout, kubeaudit.Error, nil)

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

// Audit using a local kubeconfig file
func demoLocal() {
	fmt.Println("\n--- Local ---")

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
	report, err := auditor.AuditLocal("/path/to/kubeconfig.yml")
	if err != nil {
		fmt.Println("Local file doesn't exist")
		return
	}

	// Print the audit results to screen
	report.PrintResults(os.Stdout, kubeaudit.Info, nil)
}

// Audit a cluster from inside the cluster
func demoCluster() {
	fmt.Println("\n--- Cluster ---")

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
	report, err := auditor.AuditCluster("")
	if err != nil {
		fmt.Println("Not running in a cluster")
		return
	}

	// Print the audit results to screen
	report.PrintResults(os.Stdout, kubeaudit.Info, nil)
}

// Audit using a subset of security auditors
func demoAuditorSubset() {
	fmt.Println("\n--- Auditor Subset ---")

	// Initialize the image auditor
	auditor, err := kubeaudit.New([]kubeaudit.Auditable{
		apparmor.New(),
		image.New(image.Config{Image: "myimage:mytag"}),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Run the audit in manifest mode
	report, err := auditor.AuditManifest(strings.NewReader(manifest))
	if err != nil {
		log.Fatal(err)
	}

	// Print the audit results to screen
	report.PrintResults(os.Stdout, kubeaudit.Info, &log.JSONFormatter{})
}

// Audit using a config file
// A kubeaudit config can be used to specify which security auditors to run, and to specify configuration
// for those auditors.
func demoAuditWithConfig() {
	fmt.Println("\n--- Config File ---")

	configFile := "demo/config.yaml"

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
	report.PrintResults(os.Stdout, kubeaudit.Error, nil)
}

// Audit using a custom auditor
func demoCustomAuditor() {
	fmt.Println("\n--- Custom Auditor ---")

	// Initialize kubeaudit with your auditor (see demo/custom_auditor.go for auditor definition)
	auditor, err := kubeaudit.New([]kubeaudit.Auditable{NewCustomAuditor()})
	if err != nil {
		log.Fatal(err)
	}

	// Run the audit in the mode of your choosing. Here we use manifest mode.
	report, err := auditor.AuditManifest(strings.NewReader(manifest))
	if err != nil {
		log.Fatal(err)
	}

	// Print the results to screen
	report.PrintResults(os.Stdout, kubeaudit.Info, nil)
}
