package cmd

import (
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

func setContainers(resource k8sRuntime.Object, containers []Container) k8sRuntime.Object {
	switch t := resource.(type) {
	case *DaemonSet:
		t.Spec.Template.Spec.Containers = containers
		return t.DeepCopyObject()
	case *Deployment:
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

func disableDSA(resource k8sRuntime.Object) k8sRuntime.Object {
	switch t := resource.(type) {
	case *DaemonSet:
		t.Spec.Template.Spec.DeprecatedServiceAccount = ""
		return t.DeepCopyObject()
	case *Deployment:
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

func setASAT(resource k8sRuntime.Object, b bool) k8sRuntime.Object {
	var boolean *bool
	if b {
		boolean = newTrue()
	} else {
		boolean = newFalse()
	}
	switch t := resource.(type) {
	case *DaemonSet:
		t.Spec.Template.Spec.AutomountServiceAccountToken = boolean
		return t.DeepCopyObject()
	case *Deployment:
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

func getContainers(resource k8sRuntime.Object) (container []Container) {
	switch kubeType := resource.(type) {
	case *DaemonSet:
		container = kubeType.Spec.Template.Spec.Containers
	case *Deployment:
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
