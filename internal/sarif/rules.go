package sarif

import (
	"github.com/Shopify/kubeaudit/auditors/apparmor"
	"github.com/Shopify/kubeaudit/auditors/asat"
	"github.com/Shopify/kubeaudit/auditors/capabilities"
	"github.com/Shopify/kubeaudit/auditors/deprecatedapis"
	"github.com/Shopify/kubeaudit/auditors/hostns"
	"github.com/Shopify/kubeaudit/auditors/image"
	"github.com/Shopify/kubeaudit/auditors/limits"
	"github.com/Shopify/kubeaudit/auditors/mounts"
	"github.com/Shopify/kubeaudit/auditors/netpols"
	"github.com/Shopify/kubeaudit/auditors/nonroot"
	"github.com/Shopify/kubeaudit/auditors/privesc"
	"github.com/Shopify/kubeaudit/auditors/privileged"
	"github.com/Shopify/kubeaudit/auditors/rootfs"
	"github.com/Shopify/kubeaudit/auditors/seccomp"
)

var allAuditors = map[string]string{
	apparmor.Name:       "Finds containers that do not have AppArmor enabled",
	asat.Name:           "Finds containers where the deprecated SA field is used or with a mounted default SA",
	capabilities.Name:   "Finds containers that do not drop the recommended capabilities or add new ones",
	deprecatedapis.Name: "Finds any resource defined with a deprecated API version.",
	hostns.Name:         "Finds containers that have HostPID, HostIPC or HostNetwork enabled",
	image.Name:          "Finds containers which do not use the desired version of an image (via the tag) or use an image without a tag",
	limits.Name:         "Finds containers which exceed the specified CPU and memory limits or do not specify any",
	mounts.Name:         "Finds containers that have sensitive host paths mounted",
	netpols.Name:        "Finds namespaces that do not have a default-deny network policy",
	nonroot.Name:        "Finds containers allowed to run as root",
	privesc.Name:        "Finds containers that allow privilege escalation",
	privileged.Name:     "Finds containers running as privileged",
	rootfs.Name:         "Finds containers which do not have a read-only filesystem",
	seccomp.Name:        "Finds containers running without seccomp",
}

var violationsToRules = map[string]string{
	apparmor.AppArmorAnnotationMissing:                      apparmor.Name,
	apparmor.AppArmorDisabled:                               apparmor.Name,
	apparmor.AppArmorInvalidAnnotation:                      apparmor.Name,
	asat.AutomountServiceAccountTokenDeprecated:             asat.Name,
	asat.AutomountServiceAccountTokenTrueAndDefaultSA:       asat.Name,
	capabilities.CapabilityAdded:                            capabilities.Name,
	capabilities.CapabilityOrSecurityContextMissing:         capabilities.Name,
	capabilities.CapabilityShouldDropAll:                    capabilities.Name,
	deprecatedapis.DeprecatedAPIUsed:                        deprecatedapis.Name,
	hostns.NamespaceHostIPCTrue:                             hostns.Name,
	hostns.NamespaceHostNetworkTrue:                         hostns.Name,
	hostns.NamespaceHostPIDTrue:                             hostns.Name,
	image.ImageCorrect:                                      image.Name,
	image.ImageTagIncorrect:                                 image.Name,
	image.ImageTagMissing:                                   image.Name,
	limits.LimitsCPUExceeded:                                limits.Name,
	limits.LimitsCPUNotSet:                                  limits.Name,
	limits.LimitsMemoryExceeded:                             limits.Name,
	limits.LimitsMemoryNotSet:                               limits.Name,
	limits.LimitsNotSet:                                     limits.Name,
	mounts.SensitivePathsMounted:                            mounts.Name,
	netpols.MissingDefaultDenyIngressAndEgressNetworkPolicy: netpols.Name,
	netpols.MissingDefaultDenyIngressNetworkPolicy:          netpols.Name,
	netpols.MissingDefaultDenyEgressNetworkPolicy:           netpols.Name,
	netpols.AllowAllIngressNetworkPolicyExists:              netpols.Name,
	netpols.AllowAllEgressNetworkPolicyExists:               netpols.Name,
	nonroot.RunAsUserCSCRoot:                                nonroot.Name,
	nonroot.RunAsUserPSCRoot:                                nonroot.Name,
	nonroot.RunAsNonRootCSCFalse:                            nonroot.Name,
	nonroot.RunAsNonRootPSCNilCSCNil:                        nonroot.Name,
	nonroot.RunAsNonRootPSCFalseCSCNil:                      nonroot.Name,
	privesc.AllowPrivilegeEscalationNil:                     privesc.Name,
	privesc.AllowPrivilegeEscalationTrue:                    privesc.Name,
	privileged.PrivilegedTrue:                               privileged.Name,
	privileged.PrivilegedNil:                                privileged.Name,
	rootfs.ReadOnlyRootFilesystemFalse:                      rootfs.Name,
	rootfs.ReadOnlyRootFilesystemNil:                        rootfs.Name,
	seccomp.SeccompAnnotationMissing:                        seccomp.Name,
	seccomp.SeccompDeprecatedPod:                            seccomp.Name,
	seccomp.SeccompDisabledPod:                              seccomp.Name,
	seccomp.SeccompDeprecatedContainer:                      seccomp.Name,
	seccomp.SeccompDisabledContainer:                        seccomp.Name,
}
