package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	k8sResource "k8s.io/apimachinery/pkg/api/resource"
)

type limitFlags struct {
	cpuArg    string
	cpu       k8sResource.Quantity
	memoryArg string
	memory    k8sResource.Quantity
}

var limitConfig limitFlags

func (limit *limitFlags) parseLimitFlags() {
	if len(limit.cpuArg) != 0 {
		quantity, err := k8sResource.ParseQuantity(limit.cpuArg)
		if err != nil {
			log.Error("Wrong cpu argument: " + err.Error())
			return
		}
		limit.cpu = quantity
	}

	if len(limit.memoryArg) != 0 {
		quantity, err := k8sResource.ParseQuantity(limit.memoryArg)
		if err != nil {
			log.Error("Wrong cpu argument: " + err.Error())
			return
		}
		limit.memory = quantity
	}
}

func checkLimits(container ContainerV1, limits limitFlags, result *Result) {
	if container.Resources.Limits == nil {
		occ := Occurrence{id: ErrorResourcesLimitsNil, kind: Warn, message: "Resource limit not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	checkCPULimit(container, limits, result)
	checkMemoryLimit(container, limits, result)
}

func checkCPULimit(container ContainerV1, limits limitFlags, result *Result) {
	cpu := container.Resources.Limits.Cpu()
	if cpu == nil || cpu.IsZero() {
		occ := Occurrence{id: ErrorResourcesLimitsCPUNil, kind: Warn, message: "CPU limit not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if limits.cpu.MilliValue() > 0 && cpu.MilliValue() > limits.cpu.MilliValue() {
		result.CPULimitActual = cpu.String()
		result.CPULimitMax = limits.cpu.String()
		message := fmt.Sprintf("CPU limit exceeded, it is set to %s but it must not exceed %s. Please adjust it!", cpu.String(), limits.cpu.String())
		occ := Occurrence{id: ErrorResourcesLimitsCPUExceeded, kind: Warn, message: message}
		result.Occurrences = append(result.Occurrences, occ)
	}
}

func checkMemoryLimit(container ContainerV1, limits limitFlags, result *Result) {
	memory := container.Resources.Limits.Memory()
	if memory == nil || memory.IsZero() {
		occ := Occurrence{id: ErrorResourcesLimitsMemoryNil, kind: Warn, message: "Memory limit not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if limits.memory.Value() > 0 && memory.Value() > limits.memory.Value() {
		result.MEMLimitActual = memory.String()
		result.MEMLimitMax = limits.memory.String()
		message := fmt.Sprintf("Memory limit exceeded, it is set to %s but it must not exceed %s. Please adjust it!", memory.String(), limits.memory.String())
		occ := Occurrence{id: ErrorResourcesLimitsMemoryExceeded, kind: Warn, message: message}
		result.Occurrences = append(result.Occurrences, occ)
	}
}

func auditLimits(limits limitFlags, resource Resource) (results []Result) {
	limits.parseLimitFlags()

	for _, container := range getContainers(resource) {
		result, err, warn := newResultFromResource(resource)
		if warn != nil {
			log.Warn(warn)
			return
		}
		if err != nil {
			log.Error(err)
			return
		}

		checkLimits(container, limits, result)
		if len(result.Occurrences) > 0 {
			results = append(results, *result)
		}
	}
	return
}

var limitsCmd = &cobra.Command{
	Use:   "limits",
	Short: "Audit containers running with limits",
	Long: `This command determines which containers in a kubernetes cluster have and do not exceed specified cpu and memory limits.

A PASS is given when a container has cpu and memory limits
A FAIL is given when a container does not have cpu and memory limits

Example usage:
kubeaudit limits
kubeaudit limits --cpu 500m --memory 256Mi`,
	Run: runAudit(auditLimits),
}

func init() {
	RootCmd.AddCommand(limitsCmd)
	limitsCmd.Flags().StringVar(&limitConfig.cpuArg, "cpu", "", "max cpu limit")
	limitsCmd.Flags().StringVar(&limitConfig.memoryArg, "memory", "", "max memory limit")
}
