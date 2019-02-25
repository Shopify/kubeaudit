package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
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

// FixTestSetupMultipleResources allows kubeaudit to be used programmatically instead of via the shell for multiple Resources. It is intended to be used for testing.
func FixTestSetupMultipleResources(t *testing.T, file string, auditFunction func(Resource) []Result) (*assert.Assertions, []Resource) {
	var fixedResources []Resource
	assert := assert.New(t)
	file = filepath.Join(path, file)
	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)
	for _, resource := range resources {
		results := getResults([]Resource{resource}, auditFunction)
		for _, result := range results {
			resource = fixPotentialSecurityIssue(resource, result)
		}
		fixedResources = append(fixedResources, resource)
	}
	return assert, fixedResources
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

func runFakeResourceAuditTest(t *testing.T, function interface{}, fakeResources []Resource) {
	for _, resource := range fakeResources {
		switch f := function.(type) {
		case (func(Resource) []Result):
			hook := logTest.NewGlobal()
			_ = f(resource)
			assert.Equal(t, log.ErrorLevel, hook.LastEntry().Level)
			hook.Reset()
		default:
			name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			log.Fatal("Invalid audit function provided: ", name)
		}
	}

}

func runAuditTestInNamespace(t *testing.T, namespace string, file string, function interface{}, errCodes []int) {
	rootConfig.namespace = namespace
	runAuditTest(t, file, function, errCodes)
	rootConfig.namespace = apiv1.NamespaceAll
}

// NewUnsupportedResource returns a fake unsupported resource for testing purposes
func NewUnsupportedResource() Resource {
	var unsupportedResource UnsupportedType
	return unsupportedResource.DeepCopyObject()
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

// WriteToTmpFile writes a single resource to a tmpfile, you are responsible
// for deleting the file afterwards, that's why the function returns the file
// name.
func WriteToTmpFile(decode Resource) (string, error) {
	info, _ := k8sRuntime.SerializerInfoForMediaType(scheme.Codecs.SupportedMediaTypes(), "application/yaml")
	groupVersion := schema.GroupVersion{Group: decode.GetObjectKind().GroupVersionKind().Group, Version: decode.GetObjectKind().GroupVersionKind().Version}
	encoder := scheme.Codecs.EncoderForVersion(info.Serializer, groupVersion)
	yaml, err := k8sRuntime.Encode(encoder, decode)
	if err != nil {
		return "", err
	}
	tmpfile, err := ioutil.TempFile("", "kubeaudit-test-yaml")
	if err != nil {
		return "", err
	}
	_, err = tmpfile.Write(yaml)
	if err != nil {
		return "", err
	}
	err = tmpfile.Close()
	if err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}

func compareTextFiles(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		return false
	}

	f2, err := os.Open(file2)
	if err != nil {
		return false
	}

	s1 := bufio.NewScanner(f1)
	s2 := bufio.NewScanner(f2)

	for s1.Scan() {
		s2.Scan()
		text1 := s1.Text()
		text2 := s2.Text()
		if text1 != text2 {
			fmt.Printf("Files don't match here:\n%v\n%v\n\n", text1, text2)
			return false
		}
	}
	f1stat, err := f1.Stat()
	if err != nil {
		return false
	}

	f2stat, err := f2.Stat()
	if err != nil {
		return false
	}

	if f1stat.Size() != f2stat.Size() {
		fmt.Printf("File sizes don't match")
		return false
	}
	return true
}

func assertEqualWorkloads(assert *assert.Assertions, resource1, resource2 []Resource) {
	tmpfile1, err := WriteToTmpFile(resource1[0])
	assert.Nil(err)
	defer os.Remove(tmpfile1)
	tmpfile2, err := WriteToTmpFile(resource2[0])
	assert.Nil(err)
	defer os.Remove(tmpfile2)
	assert.True(compareTextFiles(tmpfile1, tmpfile2))
}
