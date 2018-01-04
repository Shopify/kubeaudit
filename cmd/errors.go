package cmd

const (
	_ = iota
	KubeauditInternalError
	ErrorAllowPrivilegeEscalationNIL
	ErrorAllowPrivilegeEscalationTrue
	ErrorAllowPrivilegeEscalationTrueAllowed
	ErrorAutomountServiceAccountTokenNILAndNoName
	ErrorAutomountServiceAccountTokenTrueAllowed
	ErrorAutomountServiceAccountTokenTrueAndNoName
	ErrorCapabilityAdded
	ErrorCapabilityAllowed
	ErrorCapabilityNotDropped
	ErrorImageTagIncorrect
	ErrorImageTagMissing
	ErrorMisconfiguredKubeauditAllow
	ErrorPrivilegedNIL
	ErrorPrivilegedTrue
	ErrorPrivilegedTrueAllowed
	ErrorReadOnlyRootFilesystemFalse
	ErrorReadOnlyRootFilesystemFalseAllowed
	ErrorReadOnlyRootFilesystemNIL
	ErrorResourcesLimitsCpuExceeded
	ErrorResourcesLimitsCpuNIL
	ErrorResourcesLimitsMemoryExceeded
	ErrorResourcesLimitsMemoryNIL
	ErrorResourcesLimitsNIL
	ErrorRunAsNonRootFalse
	ErrorRunAsNonRootFalseAllowed
	ErrorRunAsNonRootNIL
	ErrorServiceAccountTokenDeprecated
	ErrorServiceAccountTokenNoName
	InfoImageCorrect
)
