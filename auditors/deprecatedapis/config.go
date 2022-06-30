package deprecatedapis

import (
	"fmt"
	"regexp"
	"strconv"
)

type Config struct {
	CurrentVersion  string `yaml:"currentVersion"`
	TargetedVersion string `yaml:"targetedVersion"`
}

type Version struct {
	Major int
	Minor int
}

func (config *Config) GetCurrentVersion() (*Version, error) {
	if config == nil {
		return nil, nil
	}
	return toMajorMinor(config.CurrentVersion)
}

func (config *Config) GetTargetedVersion() (*Version, error) {
	if config == nil {
		return nil, nil
	}
	return toMajorMinor(config.TargetedVersion)
}

func toMajorMinor(version string) (*Version, error) {
	if len(version) == 0 {
		return nil, nil
	}
	re := regexp.MustCompile(`^(\d{1,2})\.(\d{1,2})$`)
	if !re.MatchString(version) {
		return nil, fmt.Errorf("error parsing version: %s", version)
	}
	major, err := strconv.Atoi(re.FindStringSubmatch(version)[1])
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(re.FindStringSubmatch(version)[2])
	if err != nil {
		return nil, err
	}
	return &Version{major, minor}, nil
}
