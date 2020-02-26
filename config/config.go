package config

import (
	"io"
	"io/ioutil"

	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/auditors/limits"
	"gopkg.in/yaml.v1"
)

func New(configData io.Reader) (KubeauditConfig, error) {
	configBytes, err := ioutil.ReadAll(configData)
	if err != nil {
		return KubeauditConfig{}, err
	}

	config := KubeauditConfig{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return KubeauditConfig{}, err
	}

	return config, nil
}

type KubeauditConfig struct {
	EnabledAuditors map[string]bool `yaml:"enabledAuditors"`
	AuditorConfig   AuditorConfig   `yaml:"auditors"`
}

func (conf *KubeauditConfig) GetEnabledAuditors() []string {
	if conf == nil {
		return []string{}
	}
	enabledAuditors := make([]string, 0, len(conf.EnabledAuditors))
	for auditorName, enabled := range conf.EnabledAuditors {
		if enabled {
			enabledAuditors = append(enabledAuditors, auditorName)
		}
	}
	return enabledAuditors
}

func (conf *KubeauditConfig) GetAuditorConfigs() AuditorConfig {
	if conf == nil {
		return AuditorConfig{}
	}
	return conf.AuditorConfig
}

type AuditorConfig struct {
	Capabilities capabilities.Config `yaml:"capabilities"`
	Image        image.Config        `yaml:"image"`
	Limits       limits.Config       `yaml:"limits"`
}
