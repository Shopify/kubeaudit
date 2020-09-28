package capabilities

type Config struct {
	AddList []string `yaml:"add"`
}

func (config *Config) GetAddList() []string {
	if config == nil || len(config.AddList) == 0 {
		return []string{}
	}

	return config.AddList
}
