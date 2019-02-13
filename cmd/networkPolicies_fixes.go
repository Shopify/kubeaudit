package cmd

func fixNetworkPolicy(resource Resource, occurrence Occurrence) Resource {
	var obj Resource
	nsName := getNamespaceName(resource)
	if occurrence.id == ErrorMissingDefaultDenyIngressNetworkPolicy {
		obj = setNetworkPolicyFields(nsName, []string{"Ingress"})
	}
	if occurrence.id == ErrorMissingDefaultDenyEgressNetworkPolicy {
		obj = setNetworkPolicyFields(nsName, []string{"Egress"})
	}
	if occurrence.id == ErrorMissingDefaultDenyIngressAndEgressNetworkPolicy {
		obj = setNetworkPolicyFields(nsName, []string{"Ingress", "Egress"})
	}
	return obj
}
