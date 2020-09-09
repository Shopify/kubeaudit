package k8stypes

func IsNamespaceV1(resource Resource) bool {
	_, ok := resource.(*NamespaceV1)
	return ok
}
