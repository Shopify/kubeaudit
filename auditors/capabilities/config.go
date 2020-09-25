package capabilities

type Config struct {
	DropList []string `yaml:"drop"`
}

func (config *Config) GetDropList() []string {
	// Default Drop contains only "ALL"
	if config == nil || len(config.DropList) == 0 {
		return DefaultDropList
	}

	return config.DropList
}
