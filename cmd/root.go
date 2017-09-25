package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootConfig rootFlags

type rootFlags struct {
	kubeConfig string
	localMode  bool
	verbose    bool
	allPods    bool
	json       bool
	version    bool
}

var RootCmd = &cobra.Command{
	Use:   "kubeaudit",
	Short: "A Kubernetes security auditor",
	Long: `kubeaudit is a program that will help you audit
your Kubernetes clusters. Specify -l to run kubeaudit using ~/.kube/config
otherwise it will attempt to create an in-cluster client.

#patcheswelcome`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&rootConfig.kubeConfig, "kubeconfig", "c", "", "config file (default is $HOME/.kube/config")

	RootCmd.PersistentFlags().BoolVarP(&rootConfig.localMode, "local", "l", false, "Local mode, uses ~/.kube/config as configuration")
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.version, "version", "t", false, "Print kubeaudit version/tag")
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.verbose, "verbose", "v", false, "Enable debug (verbose) logging")
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.json, "json", "j", false, "Enable json logging")
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.allPods, "allPods", "a", false, "Audit againsts pods in all the phases (default Running Phase)")
}
