package mounts

type Config struct {
	SensitivePaths []string `yaml:"paths"`
}

func (config *Config) GetSensitivePaths() []string {
	if config == nil || len(config.SensitivePaths) == 0 {
		return DefaultSensitivePaths
	}
	return config.SensitivePaths
}
