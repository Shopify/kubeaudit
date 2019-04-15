package cmd

func fixNamespace(result *Result, resource Resource) Resource {
	switch kubeType := resource.(type) {
	case *PodV1:
		if labelExists, _ := getPodOverrideLabelReason(result, "allow-namespace-host-network"); !labelExists {
			kubeType.Spec.HostNetwork = false
			resource = kubeType.DeepCopyObject()
		}
		if labelExists, _ := getPodOverrideLabelReason(result, "allow-namespace-host-PID"); !labelExists {
			kubeType.Spec.HostPID = false
			resource = kubeType.DeepCopyObject()
		}
		if labelExists, _ := getPodOverrideLabelReason(result, "allow-namespace-host-IPC"); !labelExists {
			kubeType.Spec.HostIPC = false
			resource = kubeType.DeepCopyObject()
		}
	}
	return resource
}
