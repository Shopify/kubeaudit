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
	conf := loadKubeAuditConfigFromFile(auditAllConfig.configFile)

	// Config options set via flags override the config file
	conf = setConfigFromFlags(cmd, conf)

	auditors, err := all.Auditors(conf)
	if err != nil {
		log.WithError(err).Fatal("Error creating auditors")
	}

	runAudit(auditors...)(cmd, args)
}

func setConfigFromFlags(cmd *cobra.Command, conf config.KubeauditConfig) config.KubeauditConfig {
	flagset := cmd.Flags()
	for _, item := range []struct {
		flag      string
		flagVal   string
		configVal *string
	}{
		{imageFlagName, imageConfig.Image, &conf.AuditorConfig.Image.Image},
		{limitCpuFlagName, limitsConfig.CPU, &conf.AuditorConfig.Limits.CPU},
		{limitMemoryFlagName, limitsConfig.Memory, &conf.AuditorConfig.Limits.Memory},
	} {
		if flagset.Changed(item.flag) {
			*item.configVal = item.flagVal
		}
	}

	if flagset.Changed(capsAddFlagName) {
		conf.AuditorConfig.Capabilities.AllowAddList = capabilitiesConfig.AllowAddList
	}

	if flagset.Changed(sensitivePathsFlagName) {
		conf.AuditorConfig.Mounts.SensitivePaths = mountsConfig.SensitivePaths
	}

	return conf
}

func loadKubeAuditConfigFromFile(configFile string) config.KubeauditConfig {
	if configFile == "" {
		return config.KubeauditConfig{}
	}

	reader, err := os.Open(configFile)
	if err != nil {
		log.WithError(err).Fatal("Unable to open config file ", configFile)
	}

	conf, err := config.New(reader)
	if err != nil {
		log.WithError(err).Fatal("Error parsing config file ", configFile)
	}

	return conf
}

var auditAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Run all audits",
	Long: `Run all audits

Example usage:
kubeaudit all -f /path/to/yaml
kubeaudit all -k /path/to/kubeaudit-config.yaml /path/to/yaml
`,
	Run: auditAll,
}

func init() {
	RootCmd.AddCommand(auditAllCmd)
	auditAllCmd.Flags().StringVarP(&auditAllConfig.configFile, "kconfig", "k", "", "Path to kubeaudit config")

	// Set flags for the auditors that have them
	setImageFlags(auditAllCmd)
	setLimitsFlags(auditAllCmd)
	setCapabilitiesFlags(auditAllCmd)
	setPathsFlags(auditAllCmd)
}
