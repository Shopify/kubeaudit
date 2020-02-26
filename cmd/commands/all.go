package commands

import (
	"os"

	"github.com/Shopify/kubeaudit/auditors/all"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/auditors/limits"
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
	conf = setConfigFromFlags(conf)

	allAuditors, err := all.Auditors(conf)
	if err != nil {
		log.WithError(err).Fatal("Error creating auditors")
	}

	runAudit(allAuditors...)(cmd, args)
}

func setConfigFromFlags(conf config.KubeauditConfig) config.KubeauditConfig {
	if imageConfig != (image.Config{}) {
		conf.AuditorConfig.Image = imageConfig
	}

	if limitsConfig != (limits.Config{}) {
		conf.AuditorConfig.Limits = limitsConfig
	}

	if capabilitiesConfig != (customCapabilitiesConfig{}) {
		conf.AuditorConfig.Capabilities = capabilitiesConfig.ToConfig()
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
