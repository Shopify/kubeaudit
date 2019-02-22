package cmd

func fixServiceAccountToken(result *Result, resource Resource) Resource {
	if labelExists, _ := getPodOverrideLabelReason(result, "allow-automount-service-account-token"); labelExists {
		return resource
	}
	return setASAT(resource, false)
}

func fixDeprecatedServiceAccount(resource Resource) Resource {
	return disableDSA(resource)
}
