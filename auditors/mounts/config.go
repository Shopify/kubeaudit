package mounts

type Config struct {
	SensitivePaths []string `yaml:"paths"`
}

func (config *Config) GetSensitivePaths() []string {
	if config == nil {
		return DefaultSensitivePaths
	}
	return config.SensitivePaths
}
