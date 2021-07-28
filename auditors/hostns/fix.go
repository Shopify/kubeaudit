package hostns

import "github.com/Shopify/kubeaudit/pkg/k8s"

type fixHostNetworkTrue struct {
	podSpec *k8s.PodSpecV1
}

func (f *fixHostNetworkTrue) Plan() string {
	return "Set hostNetwork to 'false' in PodSpec"
}

func (f *fixHostNetworkTrue) Apply(resource k8s.Resource) []k8s.Resource {
	f.podSpec.HostNetwork = false
	return nil
}

type fixHostIPCTrue struct {
	podSpec *k8s.PodSpecV1
}

func (f *fixHostIPCTrue) Plan() string {
	return "Set hostIPC to 'false' in PodSpec"
}

func (f *fixHostIPCTrue) Apply(resource k8s.Resource) []k8s.Resource {
	f.podSpec.HostIPC = false
	return nil
}

type fixHostPIDTrue struct {
	podSpec *k8s.PodSpecV1
}

func (f *fixHostPIDTrue) Plan() string {
	return "Set hostPID to 'false' in PodSpec"
}

func (f *fixHostPIDTrue) Apply(resource k8s.Resource) []k8s.Resource {
	f.podSpec.HostPID = false
	return nil
}
