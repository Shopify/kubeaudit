package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

func TestFixSecurityContextNIL(t *testing.T) {
	assert := assert.New(t)
	pod := NewPod()
	pod.Spec.AutomountServiceAccountToken = newFalse()
	obj := pod.DeepCopyObject()
	results := runAllAudits(obj)
	var empty []Result
	assert.NotEqual(empty, results, "Results were empty")
	obj = fixSecurityContextNIL(obj)
	results = runAllAudits(obj)
	assert.Equal(empty, results, "Results not empty")
}

func TestFixPrivilegeEscalation(t *testing.T) {
	assert := assert.New(t)
	pod := NewPod()
	resource := pod.DeepCopyObject()
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext = &SecurityContext{
			Privileged: newTrue(),
		}
		containers = append(containers, container)
	}
	resource = setContainers(resource, containers)
	var empty []Result
	var resultsBefore []Result
	for _, result := range getResults([]k8sRuntime.Object{resource}, auditPrivileged) {
		resultsBefore = append(resultsBefore, result)
	}
	assert.NotEqual(empty, resultsBefore, "Results were empty")
	var resultsAfter []Result
	resource = fixPrivilegeEscalation(resource)
	for _, result := range getResults([]k8sRuntime.Object{resource}, auditPrivileged) {
		resultsAfter = append(resultsAfter, result)
	}
	assert.Equal(empty, resultsAfter, "Results not empty")
}

//
//			ReadOnlyRootFilesystem:   newTrue(),
//
func TestFixAllowPrivilegeEscalation(t *testing.T) {
	assert := assert.New(t)
	pod := NewPod()
	resource := pod.DeepCopyObject()
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext = &SecurityContext{
			AllowPrivilegeEscalation: newTrue(),
		}
		containers = append(containers, container)
	}
	resource = setContainers(resource, containers)
	var empty []Result
	var resultsBefore []Result
	for _, result := range getResults([]k8sRuntime.Object{resource}, auditAllowPrivilegeEscalation) {
		resultsBefore = append(resultsBefore, result)
	}
	assert.NotEqual(empty, resultsBefore, "Results were empty")
	var resultsAfter []Result
	resource = fixAllowPrivilegeEscalation(resource)
	for _, result := range getResults([]k8sRuntime.Object{resource}, auditAllowPrivilegeEscalation) {
		resultsAfter = append(resultsAfter, result)
	}
	assert.Equal(empty, resultsAfter, "Results not empty")
}

func TestFixReadOnlyRootFilesystem(t *testing.T) {
	assert := assert.New(t)
	pod := NewPod()
	resource := pod.DeepCopyObject()
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext = &SecurityContext{
			ReadOnlyRootFilesystem: newFalse(),
		}
		containers = append(containers, container)
	}
	resource = setContainers(resource, containers)
	var empty []Result
	var resultsBefore []Result
	for _, result := range getResults([]k8sRuntime.Object{resource}, auditReadOnlyRootFS) {
		resultsBefore = append(resultsBefore, result)
	}
	assert.NotEqual(empty, resultsBefore, "Results were empty")
	var resultsAfter []Result
	resource = fixReadOnlyRootFilesystem(resource)
	for _, result := range getResults([]k8sRuntime.Object{resource}, auditReadOnlyRootFS) {
		resultsAfter = append(resultsAfter, result)
	}
	assert.Equal(empty, resultsAfter, "Results not empty")
}

func TestFixRunAsNonRoot(t *testing.T) {
	assert := assert.New(t)
	pod := NewPod()
	resource := pod.DeepCopyObject()
	var containers []Container
	for _, container := range getContainers(resource) {
		container.SecurityContext = &SecurityContext{
			RunAsNonRoot: newFalse(),
		}
		containers = append(containers, container)
	}
	resource = setContainers(resource, containers)
	var empty []Result
	var resultsBefore []Result
	for _, result := range getResults([]k8sRuntime.Object{resource}, auditRunAsNonRoot) {
		resultsBefore = append(resultsBefore, result)
	}
	assert.NotEqual(empty, resultsBefore, "Results were empty")
	var resultsAfter []Result
	resource = fixRunAsNonRoot(resource)
	for _, result := range getResults([]k8sRuntime.Object{resource}, auditRunAsNonRoot) {
		resultsAfter = append(resultsAfter, result)
	}
	assert.Equal(empty, resultsAfter, "Results not empty")
}

func TestFixServiceAccountToken(t *testing.T) {
}

func TestFixDeprecatedServiceAccount(t *testing.T) {
}
