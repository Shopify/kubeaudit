package image

type Config struct {
	Image string `yaml:"image"`
}

func (config *Config) GetImage() string {
	if config == nil {
		return ""
	}
	return config.Image
}
