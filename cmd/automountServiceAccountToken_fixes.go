package cmd

import k8sRuntime "k8s.io/apimachinery/pkg/runtime"

func fixServiceAccountToken(resource k8sRuntime.Object) k8sRuntime.Object {
	return setASAT(resource, false)
}

func fixDeprecatedServiceAccount(resource k8sRuntime.Object) k8sRuntime.Object {
	return disableDSA(resource)
}
