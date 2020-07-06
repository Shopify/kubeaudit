package capabilities

import (
	"fmt"
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestFixCapabilities(t *testing.T) {
	customDropList := []string{"apple", "banana"}

	cases := []struct {
		testName     string
		overrides    []string
		add          []string
		expectedAdd  []string
		drop         []string
		expectedDrop []string
	}{
		{
			testName:     "Nothing to fix",
			overrides:    []string{},
			add:          []string{},
			expectedAdd:  []string{},
			drop:         []string{customDropList[0], customDropList[1]},
			expectedDrop: []string{customDropList[0], customDropList[1]},
		},
		{
			testName:     "Nothing to fix - all",
			overrides:    []string{},
			add:          []string{},
			expectedAdd:  []string{},
			drop:         []string{"all"},
			expectedDrop: []string{"all"},
		},
		{
			testName:     "CapabilityNotDropped",
			overrides:    []string{},
			add:          []string{},
			expectedAdd:  []string{},
			drop:         []string{},
			expectedDrop: []string{customDropList[0], customDropList[1]},
		},
		{
			testName:     "CapabilityAdded",
			overrides:    []string{},
			add:          []string{"orange", "blueberries"},
			expectedAdd:  []string{},
			drop:         []string{customDropList[0], customDropList[1]},
			expectedDrop: []string{customDropList[0], customDropList[1]},
		},
		{
			testName:     "CapabilityAdded - all",
			overrides:    []string{},
			add:          []string{"ALL"},
			expectedAdd:  []string{},
			drop:         []string{customDropList[0]},
			expectedDrop: []string{customDropList[0], customDropList[1]},
		},
		{
			testName:     "CapabilityAdded and CapabilityNotDropped",
			overrides:    []string{},
			add:          []string{customDropList[0]},
			expectedAdd:  []string{},
			drop:         []string{},
			expectedDrop: []string{customDropList[0], customDropList[1]},
		},
		{
			testName:     "Pod override",
			overrides:    []string{override.GetPodOverrideLabel(getOverrideLabel(customDropList[0]))},
			add:          []string{},
			expectedAdd:  []string{},
			drop:         []string{},
			expectedDrop: []string{customDropList[1]},
		},
		{
			testName:     "Container override",
			overrides:    []string{override.GetContainerOverrideLabel("mycontainer", getOverrideLabel(customDropList[0]))},
			add:          []string{customDropList[0], "pear"},
			expectedAdd:  []string{customDropList[0]},
			drop:         []string{},
			expectedDrop: []string{customDropList[1]},
		},
	}

	auditor := New(Config{DropList: customDropList})

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			resource := newPod(tc.add, tc.drop, tc.overrides)
			auditResults, err := auditor.Audit(resource, nil)
			if !assert.Nil(t, err) {
				return
			}

			for _, auditResult := range auditResults {
				auditResult.Fix(resource)
				ok, plan := auditResult.FixPlan()
				if ok {
					fmt.Println(plan)
				}
			}

			capabilities := k8s.GetContainers(resource)[0].SecurityContext.Capabilities
			assertCapabilitiesEqual(t, capabilities.Add, tc.expectedAdd)
			assertCapabilitiesEqual(t, capabilities.Drop, tc.expectedDrop)
		})
	}

	t.Run("Nil security context", func(t *testing.T) {
		resource := &k8stypes.PodV1{
			Spec: v1.PodSpec{
				Containers: []k8stypes.ContainerV1{{}},
			},
		}
		auditResults, err := auditor.Audit(resource, nil)
		if !assert.Nil(t, err) {
			return
		}

		for _, auditResult := range auditResults {
			auditResult.Fix(resource)
			ok, plan := auditResult.FixPlan()
			if ok {
				fmt.Println(plan)
			}
		}

		capabilities := k8s.GetContainers(resource)[0].SecurityContext.Capabilities
		assertCapabilitiesEqual(t, capabilities.Drop, customDropList)
	})
}

func assertCapabilitiesEqual(t *testing.T, capabilities []k8stypes.CapabilityV1, expected []string) {
	assert := assert.New(t)

	if !assert.Equal(len(expected), len(capabilities)) {
		return
	}

	m := make(map[string]bool)

	for _, cap := range capabilities {
		m[string(cap)] = true
	}

	for _, cap := range expected {
		ok, val := m[cap]
		assert.True(ok)
		assert.True(val)
	}
}

func newPod(add, drop, overrides []string) k8stypes.Resource {
	pod := k8stypes.NewPod()

	container := k8stypes.ContainerV1{
		Name: "mycontainer",
		SecurityContext: &k8stypes.SecurityContextV1{
			Capabilities: &k8stypes.CapabilitiesV1{
				Add:  capabilitiesFromStringArray(add),
				Drop: capabilitiesFromStringArray(drop),
			},
		},
	}
	k8s.GetPodSpec(pod).Containers = []k8stypes.ContainerV1{container}

	overrideLabels := make(map[string]string)
	for _, override := range overrides {
		overrideLabels[override] = "SomeReason"
	}
	k8s.GetPodObjectMeta(pod).SetLabels(overrideLabels)

	return pod
}

func capabilitiesFromStringArray(arr []string) []k8stypes.CapabilityV1 {
	capabilities := make([]k8stypes.CapabilityV1, 0, len(arr))
	for _, str := range arr {
		capabilities = append(capabilities, k8stypes.CapabilityV1(str))
	}
	return capabilities
}
