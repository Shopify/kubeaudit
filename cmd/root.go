package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Shopify/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
)

var rootConfig rootFlags

type rootFlags struct {
	allPods         bool
	json            bool
	kubeConfig      string
	localMode       bool
	manifest        string
	namespace       string
	verbose         string
	dropCapConfig   string
	kubeauditConfig string
}

var kubeauditConfig = &KubeauditConfig{}

// RootCmd defines the shell command usage for kubeaudit.
var RootCmd = &cobra.Command{
	Use:   "kubeaudit",
	Short: "A Kubernetes security auditor",
	Long: `kubeaudit is a program that checks security settings on your Kubernetes clusters.
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
	cobra.OnInitialize(processFlags)
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.localMode, "local", "l", false, "[DEPRECATED] Local mode, uses $HOME/.kube/config as configuration")
	RootCmd.Flags().MarkHidden("local")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.kubeConfig, "kubeconfig", "c", "", "Specify local config file (default is $HOME/.kube/config)")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.verbose, "verbose", "v", "INFO", "Set the debug level")
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.json, "json", "j", false, "Enable json logging")
	RootCmd.PersistentFlags().BoolVarP(&rootConfig.allPods, "allPods", "a", false, "Audit againsts pods in all the phases (default Running Phase)")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.namespace, "namespace", "n", apiv1.NamespaceAll, "Specify the namespace scope to audit")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.manifest, "manifest", "f", "", "yaml configuration to audit")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.dropCapConfig, "dropCapConfig", "d", "", "filepath for process capabilities to drop")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.kubeauditConfig, "kubeauditConfig", "k", "", "filepath for kubeaudit config file")
}

func processFlags() {
	if rootConfig.json {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if rootConfig.localMode == true {
		log.Warn("-l/-local is deprecated! kubeaudit will default to local mode if it's not running in a cluster. ")
		if rootConfig.kubeConfig != "" {
			return
		}

		log.Warn("To use a local kubeconfig file from inside a cluster specify '-c $HOME/.kube/config'.")
		home, ok := os.LookupEnv("HOME")
		if !ok {
			log.Fatal("Local mode selected but $HOME not set.")
		}
		rootConfig.kubeConfig = filepath.Join(home, ".kube", "config")
	}
	if rootConfig.kubeauditConfig != "" {
		data, err := ioutil.ReadFile(rootConfig.kubeauditConfig)
		if err != nil {
			log.Warn("Unable to find file at set kubeauditConfig path, auditing without any config")
			return
		}
		err = yaml.Unmarshal(data, kubeauditConfig)
		if err != nil {
			log.Println(err)
			log.Fatal("Unable to parse given kubeauditConfig file, please check the syntax of your config file")
		}
		if !kubeauditConfig.Audit {
			log.Warn("kubeaudit set to no-audit mode in kubeauditConfig!")
			os.Exit(0)
		}
	}
}
