package cmd

import (
	"io/ioutil"
	"os"

	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

func setContainers(resource Resource, containers []Container) Resource {
	switch t := resource.(type) {
	case *CronJob:
		t.Spec.JobTemplate.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DaemonSet:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DeploymentV1Beta1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DeploymentV1Beta2:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DeploymentExtensionsV1Beta1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *Pod:
		t.Spec.Containers = containers
		return t.DeepCopyObject()
	case *ReplicationController:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *StatefulSet:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	}
	return resource
}

func disableDSA(resource Resource) Resource {
	switch t := resource.(type) {
	case *CronJob:
		t.Spec.JobTemplate.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DaemonSet:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DeploymentV1Beta1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DeploymentV1Beta2:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DeploymentExtensionsV1Beta1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *Pod:
		t.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *ReplicationController:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *StatefulSet:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	}
	return resource
}

func setASAT(resource Resource, b bool) Resource {
	var boolean *bool
	if b {
		boolean = newTrue()
	} else {
		boolean = newFalse()
	}
	switch t := resource.(type) {
	case *CronJob:
		t.Spec.JobTemplate.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DaemonSet:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DeploymentV1Beta1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DeploymentV1Beta2:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DeploymentExtensionsV1Beta1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *Pod:
		t.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *ReplicationController:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *StatefulSet:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	}
	return resource
}

func getContainers(resource Resource) (container []Container) {
	switch kubeType := resource.(type) {
	case *CronJob:
		container = kubeType.Spec.JobTemplate.Spec.Template.Spec.Containers
	case *DaemonSet:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentV1Beta1:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentV1Beta2:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentExtensionsV1Beta1:
		container = kubeType.Spec.Template.Spec.Containers
	case *Pod:
		container = kubeType.Spec.Containers
	case *ReplicationController:
		container = kubeType.Spec.Template.Spec.Containers
	case *StatefulSet:
		container = kubeType.Spec.Template.Spec.Containers
	}
	return container
}

// WriteToFile writes and then appends incoming resource
func WriteToFile(decode Resource, filename string, toAppend bool) error {
	info, _ := k8sRuntime.SerializerInfoForMediaType(scheme.Codecs.SupportedMediaTypes(), "application/yaml")
	groupVersion := schema.GroupVersion{Group: decode.GetObjectKind().GroupVersionKind().Group, Version: decode.GetObjectKind().GroupVersionKind().Version}
	encoder := scheme.Codecs.EncoderForVersion(info.Serializer, groupVersion)
	yaml, err := k8sRuntime.Encode(encoder, decode)
	if err != nil {
		return err
	}
	if !toAppend {
		err = ioutil.WriteFile(filename, yaml, 0644)
		if err != nil {
			return err
		}
	} else {
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(yaml)
		if err != nil {
			return err
		}
	}
	return nil
}
