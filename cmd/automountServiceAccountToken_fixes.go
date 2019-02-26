package cmd

func fixServiceAccountToken(result *Result, resource Resource) Resource {
	return setASAT(resource, false)
}

func fixDeprecatedServiceAccount(resource Resource) Resource {
	return disableDSA(resource)
}
