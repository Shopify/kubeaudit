package cmd

func fixNetworkPolicy(resource Resource, occurrence Occurrence) Resource {
	obj := &NetworkPolicyV1{}
	nsName := getNamespaceName(resource)
	obj.ObjectMeta.Namespace = nsName
	obj.Spec.PolicyTypes = nil
	if occurrence.id == ErrorMissingDefaultDenyIngressNetworkPolicy {
		obj.Spec.PolicyTypes = append(obj.Spec.PolicyTypes, "Ingress")
	}
	if occurrence.id == ErrorMissingDefaultDenyEgressNetworkPolicy {
		obj.Spec.PolicyTypes = append(obj.Spec.PolicyTypes, "Egress")
	}
	return obj
}
