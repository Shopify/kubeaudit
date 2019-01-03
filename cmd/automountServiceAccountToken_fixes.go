package cmd

func fixServiceAccountToken(resource Resource) Resource {
	return setASAT(resource, false)
}

func fixDeprecatedServiceAccount(resource Resource) Resource {
	return disableDSA(resource)
}
