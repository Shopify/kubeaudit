package cmd

import (
	"io/ioutil"
	"os"

	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

func setContainers(resource Resource, containers []ContainerV1) Resource {
	switch t := resource.(type) {
	case *CronJobV1Beta1:
		t.Spec.JobTemplate.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DaemonSetV1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DaemonSetV1Beta1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DeploymentExtensionsV1Beta1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DeploymentV1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DeploymentV1Beta1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *DeploymentV1Beta2:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *PodV1:
		t.Spec.Containers = containers
		return t.DeepCopyObject()
	case *ReplicationControllerV1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *StatefulSetV1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *StatefulSetV1Beta1:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	}
	return resource
}

func disableDSA(resource Resource) Resource {
	switch t := resource.(type) {
	case *CronJobV1Beta1:
		t.Spec.JobTemplate.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DaemonSetV1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DaemonSetV1Beta1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DeploymentExtensionsV1Beta1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DeploymentV1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DeploymentV1Beta1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *DeploymentV1Beta2:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *PodV1:
		t.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *ReplicationControllerV1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *StatefulSetV1:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *StatefulSetV1Beta1:
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
	case *CronJobV1Beta1:
		t.Spec.JobTemplate.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DaemonSetV1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DaemonSetV1Beta1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DeploymentExtensionsV1Beta1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DeploymentV1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DeploymentV1Beta1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *DeploymentV1Beta2:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *PodV1:
		t.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *ReplicationControllerV1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *StatefulSetV1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *StatefulSetV1Beta1:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	}
	return resource
}

func setPodAnnotations(resource Resource, annotations map[string]string) Resource {
	switch kubeType := resource.(type) {
	case *CronJobV1Beta1:
		kubeType.Spec.JobTemplate.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *DaemonSetV1:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *DaemonSetV1Beta1:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *DeploymentExtensionsV1Beta1:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *DeploymentV1:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *DeploymentV1Beta1:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *DeploymentV1Beta2:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *PodV1:
		kubeType.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *ReplicationControllerV1:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *StatefulSetV1:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	case *StatefulSetV1Beta1:
		kubeType.Spec.Template.ObjectMeta.SetAnnotations(annotations)
		return kubeType.DeepCopyObject()
	}
	return resource
}

func getContainers(resource Resource) (container []ContainerV1) {
	switch kubeType := resource.(type) {
	case *CronJobV1Beta1:
		container = kubeType.Spec.JobTemplate.Spec.Template.Spec.Containers
	case *DaemonSetV1:
		container = kubeType.Spec.Template.Spec.Containers
	case *DaemonSetV1Beta1:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentExtensionsV1Beta1:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentV1:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentV1Beta1:
		container = kubeType.Spec.Template.Spec.Containers
	case *DeploymentV1Beta2:
		container = kubeType.Spec.Template.Spec.Containers
	case *PodV1:
		container = kubeType.Spec.Containers
	case *ReplicationControllerV1:
		container = kubeType.Spec.Template.Spec.Containers
	case *StatefulSetV1:
		container = kubeType.Spec.Template.Spec.Containers
	case *StatefulSetV1Beta1:
		container = kubeType.Spec.Template.Spec.Containers
	}
	return container
}

func getPodAnnotations(resource Resource) (annotations map[string]string) {
	switch kubeType := resource.(type) {
	case *CronJobV1Beta1:
		annotations = kubeType.Spec.JobTemplate.Spec.Template.ObjectMeta.GetAnnotations()
	case *DaemonSetV1:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *DaemonSetV1Beta1:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *DeploymentExtensionsV1Beta1:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *DeploymentV1:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *DeploymentV1Beta1:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *DeploymentV1Beta2:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *PodV1:
		annotations = kubeType.ObjectMeta.GetAnnotations()
	case *ReplicationControllerV1:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *StatefulSetV1:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	case *StatefulSetV1Beta1:
		annotations = kubeType.Spec.Template.ObjectMeta.GetAnnotations()
	}
	return
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
