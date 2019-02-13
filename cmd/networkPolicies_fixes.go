package cmd

func fixNetworkPolicy(resource Resource, occurrence Occurrence) Resource {
	res := &NetworkPolicyV1{}
	obj := res.DeepCopyObject()
	nsName := getNamespaceName(resource)
	if occurrence.id == ErrorMissingDefaultDenyIngressNetworkPolicy {
		obj = setNetworkPolicyFields(obj, nsName, "Ingress")
	}
	if occurrence.id == ErrorMissingDefaultDenyEgressNetworkPolicy {
		obj = setNetworkPolicyFields(obj, nsName, "Egress")
	}
	return obj
}
