package config

import (
	"io"
	"io/ioutil"

	"github.com/Shopify/kubeaudit/auditors/mounts"

	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/auditors/limits"
	"gopkg.in/yaml.v3"
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

func (conf *KubeauditConfig) GetEnabledAuditors() map[string]bool {
	if conf == nil {
		return map[string]bool{}
	}
	return conf.EnabledAuditors
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
	Mounts       mounts.Config       `yaml:"mounts"`
}
