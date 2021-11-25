package k8s

func IsNamespaceV1(resource Resource) bool {
	_, ok := resource.(*NamespaceV1)
	return ok
}

func IsPodV1(resource Resource) bool {
	_, ok := resource.(*PodV1)
	return ok
}
