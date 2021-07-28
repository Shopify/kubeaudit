package limits

import (
	"fmt"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	v1 "k8s.io/api/core/v1"
	k8sResource "k8s.io/apimachinery/pkg/api/resource"
)

const Name = "limits"

const (
	// LimitsNotSet occurs when there are no cpu and memory limits specified for a container
	LimitsNotSet = "LimitsNotSet"
	// LimitsCPUNotSet occurs when there is no cpu limit specified for a container
	LimitsCPUNotSet = "LimitsCPUNotSet"
	// LimitsMemoryNotSet occurs when there is no memory limit specified for a container
	LimitsMemoryNotSet = "LimitsMemoryNotSet"
	// LimitsCPUExceeded occurs when the CPU limit specified for a container is higher than the specified max CPU limit
	LimitsCPUExceeded = "LimitsCPUExceeded"
	// LimitsMemoryExceeded occurs when the memory limit specified for a container is higher than the specified max memory limit
	LimitsMemoryExceeded = "LimitsMemoryExceeded"
)

// Limits implements Auditable
type Limits struct {
	maxCPU    k8sResource.Quantity
	maxMemory k8sResource.Quantity
}

func New(config Config) (*Limits, error) {
	maxCPU, err := config.GetCPU()
	if err != nil {
		return nil, fmt.Errorf("error creating Limits auditor: %w", err)
	}

	maxMemory, err := config.GetMemory()
	if err != nil {
		return nil, fmt.Errorf("error creating Limits auditor: %w", err)
	}

	return &Limits{
		maxCPU:    maxCPU,
		maxMemory: maxMemory,
	}, nil
}

// Audit checks that the container cpu and memory limits do not exceed specified limits
func (limits *Limits) Audit(resource k8s.Resource, _ []k8s.Resource) ([]*kubeaudit.AuditResult, error) {
	var auditResults []*kubeaudit.AuditResult

	for _, container := range k8s.GetContainers(resource) {
		for _, auditResult := range limits.auditContainer(container) {
			if auditResult != nil {
				auditResults = append(auditResults, auditResult)
			}
		}
	}

	return auditResults, nil
}

func (limits *Limits) auditContainer(container *k8s.ContainerV1) (auditResults []*kubeaudit.AuditResult) {
	if isLimitsNil(container) {
		auditResult := &kubeaudit.AuditResult{
			Name:     LimitsNotSet,
			Severity: kubeaudit.Warn,
			Message:  "Resource limits not set.",
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
		return []*kubeaudit.AuditResult{auditResult}
	}

	containerLimits := getLimits(container)
	cpu := containerLimits.Cpu().String()
	memory := containerLimits.Memory().String()

	if isCPULimitUnset(container) {
		auditResult := &kubeaudit.AuditResult{
			Name:     LimitsCPUNotSet,
			Severity: kubeaudit.Warn,
			Message:  "Resource CPU limit not set.",
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
		auditResults = append(auditResults, auditResult)
	} else if exceedsCPULimit(container, limits) {
		maxCPU := limits.maxCPU.String()
		auditResult := &kubeaudit.AuditResult{
			Name:     LimitsCPUExceeded,
			Severity: kubeaudit.Warn,
			Message:  fmt.Sprintf("CPU limit exceeded. It is set to '%s' which exceeds the max CPU limit of '%s'.", cpu, maxCPU),
			Metadata: kubeaudit.Metadata{
				"Container":         container.Name,
				"ContainerCpuLimit": cpu,
				"MaxCPU":            maxCPU,
			},
		}
		auditResults = append(auditResults, auditResult)
	}

	if isMemoryLimitUnset(container) {
		auditResult := &kubeaudit.AuditResult{
			Name:     LimitsMemoryNotSet,
			Severity: kubeaudit.Warn,
			Message:  "Resource Memory limit not set.",
			Metadata: kubeaudit.Metadata{
				"Container": container.Name,
			},
		}
		auditResults = append(auditResults, auditResult)
	} else if exceedsMemoryLimit(container, limits) {
		maxMemory := limits.maxMemory.String()
		auditResult := &kubeaudit.AuditResult{
			Name:     LimitsMemoryExceeded,
			Severity: kubeaudit.Warn,
			Message:  fmt.Sprintf("Memory limit exceeded. It is set to '%s' which exceeds the max Memory limit of '%s'.", memory, maxMemory),
			Metadata: kubeaudit.Metadata{
				"Container":            container.Name,
				"ContainerMemoryLimit": memory,
				"MaxMemory":            maxMemory,
			},
		}
		auditResults = append(auditResults, auditResult)
	}

	return
}

func exceedsCPULimit(container *k8s.ContainerV1, limits *Limits) bool {
	containerLimits := getLimits(container)
	cpuLimit := containerLimits.Cpu().MilliValue()
	maxCPU := limits.maxCPU.MilliValue()
	return maxCPU > 0 && cpuLimit > maxCPU
}

func exceedsMemoryLimit(container *k8s.ContainerV1, limits *Limits) bool {
	containerLimits := getLimits(container)
	memoryLimit := containerLimits.Memory().Value()
	maxMemory := limits.maxMemory.Value()
	return maxMemory > 0 && memoryLimit > maxMemory
}

func isLimitsNil(container *k8s.ContainerV1) bool {
	return container.Resources.Limits == nil
}

func isCPULimitUnset(container *k8s.ContainerV1) bool {
	limits := getLimits(container)
	cpu := limits.Cpu()
	return cpu == nil || cpu.IsZero()
}

func isMemoryLimitUnset(container *k8s.ContainerV1) bool {
	limits := getLimits(container)
	memory := limits.Memory()
	return memory == nil || memory.IsZero()
}

func getLimits(container *k8s.ContainerV1) v1.ResourceList {
	if isLimitsNil(container) {
		return v1.ResourceList{}
	}

	return container.Resources.Limits
}
