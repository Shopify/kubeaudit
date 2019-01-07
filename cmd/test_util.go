package cmd

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
)

var path = "../fixtures/"

// FixTestSetup allows kubeaudit to be used programmatically instead of via the shell. It is intended to be used for testing.
func FixTestSetup(t *testing.T, file string, auditFunction func(Resource) []Result) (*assert.Assertions, Resource) {
	assert := assert.New(t)
	file = filepath.Join(path, file)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	assert.Equal(1, len(resources))
	resource := resources[0]
	results := getResults(resources, auditFunction)
	assert.Equal(1, len(results))
	result := results[0]
	return assert, fixPotentialSecurityIssue(resource, result)
}

func runAuditTest(t *testing.T, file string, function interface{}, errCodes []int, argStr ...string) (results []Result) {
	assert := assert.New(t)
	file = filepath.Join(path, file)
	var image imgFlags
	var limits limitFlags
	switch function.(type) {
	case (func(imgFlags, Resource) []Result):
		if len(argStr) != 1 {
			log.Fatal("Incorrect number of images specified")
		}
		image = imgFlags{img: argStr[0]}
		image.splitImageString()
	case (func(limitFlags, Resource) []Result):
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
	// Set manifest for test run
	rootConfig.manifest = file

	for _, resource := range resources {
		var currentResults []Result
		switch f := function.(type) {
		case (func(Resource) []Result):
			currentResults = f(resource)
		case (func(imgFlags, Resource) []Result):
			currentResults = f(image, resource)
		case (func(limitFlags, Resource) []Result):
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

// NewPod returns a simple Pod resource
func NewPod() *PodV1 {
	resources, err := getKubeResourcesManifest("../fixtures/pod_v1.yml")
	if err != nil {
		return nil
	}
	for _, resource := range resources {
		switch t := resource.(type) {
		case *PodV1:
			return t
		}
	}
	return nil
}

func assertEqualYaml(fileToFix string, fileFixed string, auditFunc func(resource Resource) []Result, t *testing.T) {
	assert, fixedResource := FixTestSetup(t, fileToFix, auditFunc)
	fileFixed = filepath.Join(path, fileFixed)
	correctlyFixedResources, err := getKubeResourcesManifest(fileFixed)
	assert.Nil(err)
	assert.Equal(correctlyFixedResources[0], fixedResource)
}
