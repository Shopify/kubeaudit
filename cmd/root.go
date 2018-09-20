package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
)

var rootConfig rootFlags

type rootFlags struct {
	allPods       bool
	json          bool
	kubeConfig    string
	localMode     bool
	manifest      string
	namespace     string
	verbose       string
	dropCapConfig string
}

// RootCmd defines the shell command useage for kubeaudit.
var RootCmd = &cobra.Command{
	Use:   "kubeaudit",
	Short: "A Kubernetes security auditor",
	Long: `kubeaudit is a program that will help you audit
your Kubernetes clusters. Specify -l to run kubeaudit using ~/.kube/config
otherwise it will attempt to create an in-cluster client.

#patcheswelcome`,
}

// Execute is a wrapper for the RootCmd.Execute method which will exit the program if there is an error.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&rootConfig.kubeConfig, "kubeconfig", "c", "", "config file (default is $HOME/.kube/config")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.verbose, "verbose", "v", "INFO", "Set the debug level")
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.localMode, "local", "l", false, "Local mode, uses ~/.kube/config as configuration")
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.json, "json", "j", false, "Enable json logging")
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.allPods, "allPods", "a", false, "Audit againsts pods in all the phases (default Running Phase)")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.namespace, "namespace", "n", apiv1.NamespaceAll, "Specify the namespace scope to audit")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.manifest, "manifest", "f", "", "yaml configuration to audit")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.dropCapConfig, "dropCapConfig", "d", "", "yaml configuration to audit")
}
