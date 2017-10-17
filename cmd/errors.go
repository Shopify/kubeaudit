package cmd

const (
	KubeAuditInfo = iota
	ErrorImageTagMissing
	ErrorImageTagIncorrect
	ErrorSecurityContextNIL
	ErrorReadOnlyRootFilesystemNIL
	ErrorReadOnlyRootFilesystemFalse
	ErrorRunAsNonRootNIL
	ErrorRunAsNonRootFalse
	ErrorCapabilitiesNIL
	ErrorCapabilitiesAddedOrNotDropped
)
