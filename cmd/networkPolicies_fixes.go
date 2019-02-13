package cmd

import "fmt"

func fixNetworkPolicy(resource Resource, occurrence Occurrence) Resource {
	var obj Resource
	fmt.Println("I was HERe")
	fmt.Println("I was HERe")
	fmt.Println("I was HERe")
	nsName := getNamespaceName(resource)
	if occurrence.id == ErrorMissingDefaultDenyIngressNetworkPolicy {
		obj = setNetworkPolicyFields(obj, nsName, "Ingress")
	}
	if occurrence.id == ErrorMissingDefaultDenyEgressNetworkPolicy {
		obj = setNetworkPolicyFields(obj, nsName, "Egress")
	}
	return obj
}
