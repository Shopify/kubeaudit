package cmd

// Error codes
const (
	_ = iota
	// KubeauditInternalError is an internal error which cannot be fixed by the user.
	KubeauditInternalError
	// ErrorAllowPrivilegeEscalationNIL occurs when AllowPrivilegeEscalation is not
	// set which allows privilege escalation.
	ErrorAllowPrivilegeEscalationNIL
	// ErrorAllowPrivilegeEscalationTrue occurs when AllowPrivilegeEscalation is set to true
	ErrorAllowPrivilegeEscalationTrue
	// ErrorAllowPrivilegeEscalationTrueAllowed occurs when AllowPrivilegeEscalation is
	// allowed to be set to true.
	ErrorAllowPrivilegeEscalationTrueAllowed
	// ErrorAutomountServiceAccountTokenNILAndNoName occurs when automountServiceAccountToken
	// is not set and serviceAccountName is blank.
	ErrorAutomountServiceAccountTokenNILAndNoName
	// ErrorAutomountServiceAccountTokenTrueAllowed occurs when automountServiceAccountToken
	// is allowed to be set to true.
	ErrorAutomountServiceAccountTokenTrueAllowed
	// ErrorAutomountServiceAccountTokenTrueAndNoName occurs when automountServiceAccountToken
	// is set as true and serviceAccountName is blank.
	ErrorAutomountServiceAccountTokenTrueAndNoName
	// ErrorCapabilityAdded occurs when a capability is added that is not allowed
	ErrorCapabilityAdded
	// ErrorCapabilityAllowed occurs when a capability is allowed that is part of the
	// toBeDropped list.
	ErrorCapabilityAllowed
	// ErrorCapabilityNotDropped occurs when a capability should be dropped but it isn't
	ErrorCapabilityNotDropped
	// ErrorImageTagIncorrect occurs when an incorrect image tag is provided.
	ErrorImageTagIncorrect
	// ErrorImageTagMissing occurs when there is no image tag provided.
	ErrorImageTagMissing
	// ErrorMisconfiguredKubeauditAllow occurs when the option to allow a setting is set to
	// true but the option itself is set to false or nil.
	ErrorMisconfiguredKubeauditAllow
	// ErrorPrivilegedNIL occurs when Privileged is not set.
	ErrorPrivilegedNIL
	// ErrorPrivilegedTrue occurs when Privileged is set to true.
	ErrorPrivilegedTrue
	// ErrorPrivilegedTrueAllowed occurs when Privileged is allowed to be set to true.
	ErrorPrivilegedTrueAllowed
	// ErrorReadOnlyRootFilesystemFalse occurs when ReadOnlyRootFilesystem is set to false.
	ErrorReadOnlyRootFilesystemFalse
	// ErrorReadOnlyRootFilesystemFalseAllowed occurs when ReadOnlyRootFilesystem is allowed
	// to be set to false.
	ErrorReadOnlyRootFilesystemFalseAllowed
	// ErrorReadOnlyRootFilesystemNIL occurs when ReadOnlyRootFilesystem is set to nil.
	ErrorReadOnlyRootFilesystemNIL
	// ErrorResourcesLimitsCPUExceeded occurs when the CPU limit is exceeded.
	ErrorResourcesLimitsCPUExceeded
	// ErrorResourcesLimitsCPUNIL occurs when the CPU limit is not set.
	ErrorResourcesLimitsCPUNIL
	// ErrorResourcesLimitsMemoryExceeded occurs when the memory limit is exceeded.
	ErrorResourcesLimitsMemoryExceeded
	// ErrorResourcesLimitsMemoryNIL occurs when the memory limit is not set.
	ErrorResourcesLimitsMemoryNIL
	// ErrorResourcesLimitsNIL occurs when the resource limit is set to nil.
	ErrorResourcesLimitsNIL
	// ErrorRunAsNonRootFalse occurs when RunAsNonRoot is set to false.
	ErrorRunAsNonRootFalse
	// ErrorRunAsNonRootFalseAllowed occurs when RunAsNonRoot is allowed to be set to false.
	ErrorRunAsNonRootFalseAllowed
	// ErrorRunAsNonRootNIL occurs when RunAsNonRoot is not set.
	ErrorRunAsNonRootNIL
	// ErrorServiceAccountTokenDeprecated occurs when serviceAccount is used. ServiceAccount
	// is a deprecated alias for ServiceAccountName.
	ErrorServiceAccountTokenDeprecated
	// InfoImageCorrect occurs when an image tag is correct.
	InfoImageCorrect
)
