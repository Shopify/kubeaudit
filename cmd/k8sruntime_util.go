package cmd

import (
	"io/ioutil"
	"os"

	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

func setContainers(resource k8sRuntime.Object, containers []Container) k8sRuntime.Object {
	switch t := resource.(type) {
	case *CronJob:
		t.Spec.JobTemplate.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DaemonSet:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DaemonSetV1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DeploymentV1Beta1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DeploymentV1Beta2:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DeploymentV1:
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
	case *StatefulSetV1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	}
	return resource
}

func disableDSA(resource k8sRuntime.Object) k8sRuntime.Object {
	switch t := resource.(type) {
	case *CronJob:
		t.Spec.JobTemplate.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DaemonSet:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DaemonSetV1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DeploymentV1Beta1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DeploymentV1Beta2:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DeploymentV1:
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
	case *StatefulSetV1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	}
	return resource
}

func setASAT(resource k8sRuntime.Object, b bool) k8sRuntime.Object {
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
	case *DaemonSetV1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DeploymentV1Beta1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DeploymentV1Beta2:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DeploymentV1:
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
	case *StatefulSetV1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	}
	return resource
}

func setPodAnnotations(resource k8sRuntime.Object, annotations map[string]string) k8sRuntime.Object {
	switch kubeType := resource.(type) {
	case *CronJob:
		kubeType.Spec.JobTemplate.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *DaemonSet:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *DeploymentV1Beta1:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *DeploymentV1Beta2:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *DeploymentExtensionsV1Beta1:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *Pod:
		kubeType.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *ReplicationController:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *StatefulSet:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	}
	return resource
}

func getContainers(resource k8sRuntime.Object) (container []Container) {
	switch kubeType := resource.(type) {
	case *CronJob:
		container = kubeType.Spec.JobTemplate.Spec.Template.Spec.Containers
	case *DaemonSet:
		container = kubeType.Spec.Template.Spec.Containers
	case *DaemonSetV1:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentV1Beta1:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentV1Beta2:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentV1:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentExtensionsV1Beta1:
		container = kubeType.Spec.Template.Spec.Containers
	case *Pod:
		container = kubeType.Spec.Containers
	case *ReplicationController:
		container = kubeType.Spec.Template.Spec.Containers
	case *StatefulSet:
		container = kubeType.Spec.Template.Spec.Containers
	case *StatefulSetV1:
		container = kubeType.Spec.Template.Spec.Containers
	}
	return container
}

func getPodAnnotations(resource k8sRuntime.Object) (annotations map[string]string) {
	switch kubeType := resource.(type) {
	case *CronJob:
		annotations = kubeType.Spec.JobTemplate.Spec.Template.ObjectMeta.GetAnnotations()
	case *DaemonSet:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *DeploymentV1Beta1:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *DeploymentV1Beta2:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *DeploymentExtensionsV1Beta1:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *Pod:
		annotations = kubeType.ObjectMeta.GetAnnotations()
	case *ReplicationController:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *StatefulSet:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	}
	return
}

// WriteToFile writes and then appends incoming resource
func WriteToFile(decode k8sRuntime.Object, filename string, toAppend bool) error {
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
