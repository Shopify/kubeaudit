package limits

import (
	"fmt"

	k8sResource "k8s.io/apimachinery/pkg/api/resource"
)

type Config struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

func (config *Config) GetCPU() (k8sResource.Quantity, error) {
	cpuArg := ""
	if config != nil {
		cpuArg = config.CPU
	}
	if cpuArg != "" {
		CPU, err := k8sResource.ParseQuantity(cpuArg)
		if err != nil {
			return CPU, fmt.Errorf("error parsing max CPU limit: %w", err)
		}
		return CPU, nil
	}
	return k8sResource.Quantity{}, nil
}

func (config *Config) GetMemory() (k8sResource.Quantity, error) {
	memoryArg := ""
	if config != nil {
		memoryArg = config.Memory
	}
	if memoryArg != "" {
		memory, err := k8sResource.ParseQuantity(memoryArg)
		if err != nil {
			return memory, fmt.Errorf("error parsing max memory limit: %w", err)
		}
		return memory, nil
	}
	return k8sResource.Quantity{}, nil
}
