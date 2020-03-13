package commands

import (
	"os"

	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var auditAllConfig struct {
	configFile string
}

func auditAll(cmd *cobra.Command, args []string) {
	conf := loadConfigFromFile(auditAllConfig.configFile)

	// Config options set via flags override the config file
	conf = setConfigFromFlags(cmd, conf)

	allAuditors, err := all.Auditors(conf)
	if err != nil {
		log.WithError(err).Fatal("Error creating auditors")
	}

	runAudit(allAuditors...)(cmd, args)
}

func setConfigFromFlags(cmd *cobra.Command, conf config.KubeauditConfig) config.KubeauditConfig {
	flagset := cmd.Flags()
	for _, item := range []struct {
		flag      string
		flagVal   string
		configVal *string
	}{
		{"image", imageConfig.Image, &conf.AuditorConfig.Image.Image},
		{"cpu", limitsConfig.CPU, &conf.AuditorConfig.Limits.CPU},
		{"memory", limitsConfig.Memory, &conf.AuditorConfig.Limits.Memory},
	} {
		if flagset.Changed(item.flag) {
			*item.configVal = item.flagVal
		}
	}

	if flagset.Changed("drop") {
		conf.AuditorConfig.Capabilities.DropList = capabilitiesConfig.DropList
	}

	return conf
}

func loadConfigFromFile(configFile string) config.KubeauditConfig {
	if auditAllConfig.configFile == "" {
		return config.KubeauditConfig{}
	}

	reader, err := os.Open(auditAllConfig.configFile)
	if err != nil {
		log.WithError(err).Fatal("Unable to open config file ", auditAllConfig.configFile)
	}

	conf, err := config.New(reader)
	if err != nil {
		log.WithError(err).Fatal("Error parsing config file ", auditAllConfig.configFile)
	}

	return conf
}

var auditAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Run all audits",
	Long: `Run all audits

Example usage:
kubeaudit all -f /path/to/yaml`,
	Run: auditAll,
}

func init() {
	RootCmd.AddCommand(auditAllCmd)
	auditAllCmd.Flags().StringVarP(&auditAllConfig.configFile, "kconfig", "k", "", "Path to kubeaudit config")

	// Set flags for the auditors that have them
	setImageFlags(auditAllCmd)
	setLimitsFlags(auditAllCmd)
	setCapabilitiesFlags(auditAllCmd)
}
