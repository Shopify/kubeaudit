package capabilities

type Config struct {
	AllowAddList []string `yaml:"allowAddList"`
}

func (config *Config) GetAllowAddList() []string {
	if config == nil || len(config.AllowAddList) == 0 {
		return DefaultAllowAddList
	}

	return config.AllowAddList
}
