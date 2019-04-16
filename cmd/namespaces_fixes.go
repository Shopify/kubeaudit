package cmd

func fixNamespace(result *Result, resource Resource) Resource {
	switch kubeType := resource.(type) {
	case *PodV1:
		if labelExists, _ := getPodOverrideLabelReason(result, "allow-namespace-host-network"); !labelExists {
			kubeType.Spec.HostNetwork = false
		}
		if labelExists, _ := getPodOverrideLabelReason(result, "allow-namespace-host-PID"); !labelExists {
			kubeType.Spec.HostPID = false
		}
		if labelExists, _ := getPodOverrideLabelReason(result, "allow-namespace-host-IPC"); !labelExists {
			kubeType.Spec.HostIPC = false
		}
		return kubeType.DeepCopyObject()
	}
	return resource
}
