package cmd

// Error codes
const (
	_ = iota
	// KubeauditInternalError is an internal error which cannot be fixed by the user.
	KubeauditInternalError
	// ErrorAllowPrivilegeEscalationNil occurs when AllowPrivilegeEscalation is not set which allows privilege
	// escalation.
	ErrorAllowPrivilegeEscalationNil
	// ErrorAllowPrivilegeEscalationTrue occurs when AllowPrivilegeEscalation is set to true
	ErrorAllowPrivilegeEscalationTrue
	// ErrorAllowPrivilegeEscalationTrueAllowed occurs when AllowPrivilegeEscalation is allowed to be set to true.
	ErrorAllowPrivilegeEscalationTrueAllowed
	// ErrorAutomountServiceAccountTokenNilAndNoName occurs when automountServiceAccountToken is not set and
	// serviceAccountName is blank.
	ErrorAutomountServiceAccountTokenNilAndNoName
	// ErrorAutomountServiceAccountTokenTrueAllowed occurs when automountServiceAccountToken is allowed to be set
	// to true.
	ErrorAutomountServiceAccountTokenTrueAllowed
	// ErrorAutomountServiceAccountTokenTrueAndNoName occurs when automountServiceAccountToken is set as true and
	// serviceAccountName is blank.
	ErrorAutomountServiceAccountTokenTrueAndNoName
	// ErrorCapabilityAdded occurs when a capability is added that is not allowed
	ErrorCapabilityAdded
	// ErrorCapabilityAllowed occurs when a capability is allowed that is part of the toBeDropped list.
	ErrorCapabilityAllowed
	// ErrorCapabilityNotDropped occurs when a capability should be dropped but it isn't
	ErrorCapabilityNotDropped
	// ErrorImageTagIncorrect occurs when an incorrect image tag is provided.
	ErrorImageTagIncorrect
	// ErrorImageTagMissing occurs when there is no image tag provided.
	ErrorImageTagMissing
	// ErrorMisconfiguredKubeauditAllow occurs when the option to allow a setting is set to true but the option
	// itself is set to false or nil.
	ErrorMisconfiguredKubeauditAllow
	// ErrorPrivilegedNil occurs when Privileged is not set.
	ErrorPrivilegedNil
	// ErrorPrivilegedTrue occurs when Privileged is set to true.
	ErrorPrivilegedTrue
	// ErrorPrivilegedTrueAllowed occurs when Privileged is allowed to be set to true.
	ErrorPrivilegedTrueAllowed
	// ErrorReadOnlyRootFilesystemFalse occurs when ReadOnlyRootFilesystem is set to false.
	ErrorReadOnlyRootFilesystemFalse
	// ErrorReadOnlyRootFilesystemFalseAllowed occurs when ReadOnlyRootFilesystem is allowed to be set to false.
	ErrorReadOnlyRootFilesystemFalseAllowed
	// ErrorReadOnlyRootFilesystemNil occurs when ReadOnlyRootFilesystem is set to nil.
	ErrorReadOnlyRootFilesystemNil
	// ErrorResourcesLimitsCPUExceeded occurs when the CPU limit is exceeded.
	ErrorResourcesLimitsCPUExceeded
	// ErrorResourcesLimitsCPUNil occurs when the CPU limit is not set.
	ErrorResourcesLimitsCPUNil
	// ErrorResourcesLimitsMemoryExceeded occurs when the memory limit is exceeded.
	ErrorResourcesLimitsMemoryExceeded
	// ErrorResourcesLimitsMemoryNil occurs when the memory limit is not set.
	ErrorResourcesLimitsMemoryNil
	// ErrorResourcesLimitsNil occurs when the resource limit is set to nil.
	ErrorResourcesLimitsNil
	// ErrorRunAsNonRootPSCTrueCSCFalse occurs when RunAsNonRoot is set to false in the ContainerSecurityContext and to true/false in PodSecurityContext.
	ErrorRunAsNonRootPSCTrueFalseCSCFalse
	// ErrorRunAsNonRootPSCFalseCSCNil occurs when RunAsNonRoot is Nil in the ContainerSecurityContext and to false in Pod ecurityContext.
	ErrorRunAsNonRootPSCFalseCSCNil
	// ErrorRunAsNonRootFalseAllowed occurs when RunAsNonRoot is allowed to be set to false.
	ErrorRunAsNonRootFalseAllowed
	// ErrorRunAsNonRootNil occurs when RunAsNonRoot is not set in either PodSecurityContext or ContainerSecurityContext.
	ErrorRunAsNonRootPSCNilCSCNil
	// ErrorServiceAccountTokenDeprecated occurs when serviceAccount is used. ServiceAccount is a deprecated alias
	// for ServiceAccountName.
	ErrorServiceAccountTokenDeprecated
	// ErrorAppArmorDisabled occurs when the AppArmor annotation is set to a bad value.
	ErrorAppArmorDisabled
	// ErrorAppArmorAnnotationMissing occurs when there is no annotation enabling AppArmor on the pod.
	ErrorAppArmorAnnotationMissing
	// ErrorSeccompDisabledPod occurs when the Seccomp annotation is set to a bad value.
	ErrorSeccompDisabledPod
	// ErrorSeccompDisabled occurs when the Seccomp annotation is set to a bad value.
	ErrorSeccompDisabled
	// ErrorSeccompAnnotationMissing occurs when there is no annotation enabling Seccomp on the pod.
	ErrorSeccompAnnotationMissing
	// ErrorSeccompDeprecatedPod occurs when the Seccomp annotation is set to a deprecated value.
	ErrorSeccompDeprecatedPod
	// ErrorSeccompDeprecated occurs when the Seccomp annotation is set to a deprecated value.
	ErrorSeccompDeprecated
	// InfoImageCorrect occurs when an image tag is correct.
	InfoImageCorrect
	// ErrorMissingDefaultDenyIngressAndEgressNetworkPolicy missing a default deny egress and default deny egress NetworkPolicy but it's set to be allowed
	ErrorMissingDefaultDenyIngressAndEgressNetworkPolicy
	// ErrorMissingDefaultDenyIngressAndEgressNetworkPolicyAllowed occurs when missing a default deny egress and default deny egress NetworkPolicy but it's set to be allowed
	ErrorMissingDefaultDenyIngressAndEgressNetworkPolicyAllowed
	// ErrorMissingDefaultDenyEgressNetworkPolicy occurs when a namespace is missing a default deny egress NetworkPolicy
	ErrorMissingDefaultDenyEgressNetworkPolicy
	// ErrorMissingDefaultDenyEgressNetworkPolicyAllowed occurs when a namespace is missing a default deny egress NetworkPolicy but it's allowed
	ErrorMissingDefaultDenyEgressNetworkPolicyAllowed
	// ErrorMissingDefaultDenyEgressNetworkPolicy occurs when a namespace is missing a default deny ingress NetworkPolicy
	ErrorMissingDefaultDenyIngressNetworkPolicy
	// ErrorMissingDefaultDenyIngressNetworkPolicyAllowed  occurs when a namespace is missing a default deny ingress NetworkPolicy but it's allowed
	ErrorMissingDefaultDenyIngressNetworkPolicyAllowed
	//  ErrorNamespaceHostIPCTrue occurs when a hostIPC is set to true in PodSpec
	ErrorNamespaceHostIPCTrue
	//  ErrorNamespaceHostIPCTrueAllowed occurs when a hostIPC is set to true in PodSpec but it's allowed
	ErrorNamespaceHostIPCAllowed
	//  ErrorNamespaceHostIPCTrue occurs when a hostNetwork is set to true in PodSpec
	ErrorNamespaceHostNetworkTrue
	//  ErrorNamespaceHostIPCTrueAllowed occurs when a hostNetwork is set to true in PodSpec but it's allowed
	ErrorNamespaceHostNetworkAllowed
	//  ErrorNamespaceHostIPCTrue occurs when a hostPID is set to true in PodSpec
	ErrorNamespaceHostPIDTrue
	//  ErrorNamespaceHostIPCTrue occurs when a hostPID is set to true in PodSpec but it's allowed
	ErrorNamespaceHostPIDAllowed
	// InfoDefaultDenyNetworkPolicyExists occurs when a namespace has a default deny NetworkPolicy
	InfoDefaultDenyNetworkPolicyExists
	// WarningAllowAllIngressNetworkPolicyExists occurs when a namespace has an allow all ingress NetworkPolicy
	WarningAllowAllIngressNetworkPolicyExists
	// WarningAllowAllEgressNetworkPolicyExists occurs when a namespace has an allow all egress NetworkPolicy
	WarningAllowAllEgressNetworkPolicyExists
)
