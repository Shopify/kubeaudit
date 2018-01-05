package cmd

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

var path = "../fixtures/"

func FixTestSetup(t *testing.T, file string, auditFunction func(k8sRuntime.Object) []Result, fixFunction interface{}) (*assert.Assertions, k8sRuntime.Object) {
	assert := assert.New(t)
	file = filepath.Join(path, file)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	assert.Equal(1, len(resources))
	resource := resources[0]
	results := getResults(resources, auditFunction)
	assert.Equal(1, len(results))
	result := results[0]
	for _, occurrence := range result.Occurrences {
		switch function := fixFunction.(type) {
		case func(k8sRuntime.Object) k8sRuntime.Object:
			resource = function(resource)
		case func(k8sRuntime.Object, Occurrence) k8sRuntime.Object:
			resource = function(resource, occurrence)
		}
	}
	return assert, resource
}

func runAuditTest(t *testing.T, file string, function interface{}, errCodes []int, argStr ...string) (results []Result) {
	assert := assert.New(t)
	file = filepath.Join(path, file)
	var image imgFlags
	var limits limitFlags
	switch function.(type) {
	case (func(imgFlags, k8sRuntime.Object) []Result):
		if len(argStr) != 1 {
			log.Fatal("Incorrect number of images specified")
		}
		image = imgFlags{img: argStr[0]}
		image.splitImageString()
	case (func(limitFlags, k8sRuntime.Object) []Result):
		if len(argStr) == 2 {
			limits = limitFlags{cpuArg: argStr[0], memoryArg: argStr[1]}
		} else if len(argStr) == 0 {
			limits = limitFlags{cpuArg: "", memoryArg: ""}
		} else {
			log.Fatal("Incorrect number of images specified")
		}
		limits.parseLimitFlags()
	}

	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)

	for _, resource := range resources {
		var currentResults []Result
		switch f := function.(type) {
		case (func(k8sRuntime.Object) []Result):
			currentResults = f(resource)
		case (func(imgFlags, k8sRuntime.Object) []Result):
			currentResults = f(image, resource)
		case (func(limitFlags, k8sRuntime.Object) []Result):
			currentResults = f(limits, resource)
		default:
			name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			log.Fatal("Invalid audit function provided: ", name)
		}
		for _, currentResult := range currentResults {
			results = append(results, currentResult)
		}
	}
	errors := map[int]bool{}
	for _, result := range results {
		for _, occurrence := range result.Occurrences {
			errors[occurrence.id] = true
		}
	}

	assert.Equal(len(errCodes), len(errors))
	for _, errCode := range errCodes {
		assert.True(errors[errCode])
	}
	return
}

func runAuditTestInNamespace(t *testing.T, namespace string, file string, function interface{}, errCodes []int) {
	rootConfig.namespace = namespace
	runAuditTest(t, file, function, errCodes)
	rootConfig.namespace = apiv1.NamespaceAll
}

func NewPod() *Pod {
	return &Pod{
		ObjectMeta: ObjectMeta{
			Namespace: "namespace",
			Name:      "name",
			Labels: map[string]string{
				"label": "label",
			},
			Annotations: map[string]string{
				"annotation": "annotation",
			},
		},
		Spec: PodSpec{
			Containers: []Container{
				{
					Image:   "image",
					Name:    "container",
					Command: []string{"sh", "-c", "command"},
				},
			},
		},
	}
}
