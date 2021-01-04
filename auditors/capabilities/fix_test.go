package capabilities

import (
	"fmt"
	"testing"

	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/internal/override"
	"github.com/Shopify/kubeaudit/k8stypes"
	v1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/assert"
)

func TestFixCapabilities(t *testing.T) {
	customAddList := []string{"apple", "banana"}

	cases := []struct {
		testName     string
		overrides    []string
		add          []string
		expectedAdd  []string
		drop         []string
		expectedDrop []string
	}{
		{
			testName:     "Capabilities not set to ALL",
			overrides:    []string{},
			add:          []string{},
			expectedAdd:  []string{},
			drop:         []string{"orange", "banana"},
			expectedDrop: []string{"ALL"},
		},
		{
			testName:     "Nothing to fix - no caps added and drop is set to all",
			overrides:    []string{},
			add:          []string{},
			expectedAdd:  []string{},
			drop:         []string{"all"},
			expectedDrop: []string{"all"},
		},
		{
			testName:     "No capabilities specified",
			overrides:    []string{},
			add:          []string{},
			expectedAdd:  []string{},
			drop:         []string{},
			expectedDrop: []string{"ALL"},
		},
		{
			testName:     "Capability Added with with override specified in AddList and 2 capabilities dropped",
			overrides:    []string{},
			add:          []string{customAddList[0], customAddList[1], "orange"},
			expectedAdd:  []string{customAddList[0], customAddList[1]},
			drop:         []string{"pineapple", "pomegranate"},
			expectedDrop: []string{"ALL"},
		},
		{
			testName:     "CapabilityAdded - all",
			overrides:    []string{},
			add:          []string{"ALL"},
			expectedAdd:  []string{},
			drop:         []string{"orange"},
			expectedDrop: []string{"ALL"},
		},
		{
			testName:     "CapabilityAdded",
			overrides:    []string{},
			add:          []string{"orange"},
			expectedAdd:  []string{},
			drop:         []string{},
			expectedDrop: []string{"ALL"},
		},
		{
			testName:     "Pod override",
			overrides:    []string{override.GetPodOverrideLabel(getOverrideLabel("orange"))},
			add:          []string{"orange"},
			expectedAdd:  []string{"orange"},
			drop:         []string{},
			expectedDrop: []string{"ALL"},
		},
		{
			testName:     "Container override",
			overrides:    []string{override.GetContainerOverrideLabel("mycontainer", getOverrideLabel("pear"))},
			add:          []string{customAddList[0], "pear", "orange"},
			expectedAdd:  []string{customAddList[0], "pear"},
			drop:         []string{},
			expectedDrop: []string{"ALL"},
		},
		{
			testName:     "CapabilityAdded with 3 override labels",
			overrides:    []string{override.GetContainerOverrideLabel("mycontainer", getOverrideLabel("blueberries")), override.GetContainerOverrideLabel("mycontainer", getOverrideLabel("strawberries")), override.GetContainerOverrideLabel("mycontainer", getOverrideLabel("raspberries"))},
			add:          []string{customAddList[0], "blueberries", "raspberries", "strawberries", "orange"},
			expectedAdd:  []string{customAddList[0], "blueberries", "raspberries", "strawberries"},
			drop:         []string{},
			expectedDrop: []string{"ALL"},
		},
	}

	auditor := New(Config{AddList: customAddList})

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
		assertCapabilitiesEqual(t, capabilities.Drop, []string{"ALL"})
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
