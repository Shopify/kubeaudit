package k8stypes

func IsNamespaceV1(resource Resource) bool {
	_, ok := resource.(*NamespaceV1)
	return ok
}

func IsNetworkPolicyV1(resource Resource) bool {
	_, ok := resource.(*NetworkPolicyV1)
	return ok
}

func IsServiceAccountV1(resource Resource) bool {
	_, ok := resource.(*ServiceAccountV1)
	return ok
}
