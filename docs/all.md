# All Auditor (all)

Runs all available auditors, or those specified using a kubeaudit config.

## General Usage

```
kubeaudit all [flags]
```

## Flags

| Short   | Long       | Description                               | Default                                  |
| :------ | :--------- | :---------------------------------------- | :--------------------------------------- |
| -k      | --kconfig  | Path to kubeaudit config                  |                                          |

Also see [Global Flags](/README.md#global-flags)

### Kubeaudit Config

The kubeaudit config can be used for two things:
1. Enabling only some auditors
1. Specifying configuration for auditors

Any configuration that can be specified using flags for the individual auditors can be represented using the config.

The config has the following format:

```yaml
enabledAuditors:
    # Auditors are enabled by default if they are not explicitly set to "false"
    apparmor: false
    asat: false
    capabilities: true
    hostns: true
    image: true
    limits: true
    mountds: true
    netpols: true
    nonroot: true
    privesc: true
    privileged: true
    rootfs: true
    seccomp: true
auditors:
    capabilities:
        # If no capabilities are specified and the 'capabilities' auditor is enabled,
        # a list of recommended capabilities to drop is used
        drop: ["AUDIT_WRITE", "CHOWN"]
    image:
        # If no image is specified and the 'image' auditor is enabled, WARN results
        # will be generated for containers which use an image without a tag
        image: "myimage:mytag"
    limits:
        # If no limits are specified and the 'limits' auditor is enabled, WARN results
        # will be generated for containers which have no cpu or memory limits specified
        cpu: "750m"
        memory: "500m"
```

For more details about each auditor, including a description of the auditor-specific configuration in the config, see the [Auditor Docs](/README.md#auditors).

**Note**: The kubeaudit config is not the same as the kubeconfig file specified with the `-c/--kubeconfig` flag, which refers to the Kubernetes config file (see [Local Mode](/README.md#local-mode)). Also note that only the `all` command supports using a kubeaudit config. It will not work with other commands.

## Examples

```
$ kubeaudit all -f "auditors/all/fixtures/audit_all_v1.yml"
ERRO[0000] AppArmor annotation missing. The annotation 'container.apparmor.security.beta.kubernetes.io/fakeContainerSC' should be added.  AuditResultName=AppArmorAnnotationMissing Container=fakeContainerSC MissingAnnotation=container.apparmor.security.beta.kubernetes.io/fakeContainerSC
ERRO[0000] Default service account with token mounted. automountServiceAccountToken should be set to 'false' or a non-default service account should be used.  AuditResultName=AutomountServiceAccountTokenTrueAndDefaultSA
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=AUDIT_WRITE Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=CHOWN Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=DAC_OVERRIDE Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=FOWNER Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=FSETID Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=KILL Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=MKNOD Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=NET_BIND_SERVICE Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=NET_RAW Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=SETFCAP Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=SETGID Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=SETPCAP Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=SETUID Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=SYS_CHROOT Container=fakeContainerSC
ERRO[0000] hostNetwork is set to 'true' in PodSpec. It should be set to 'false'.  AuditResultName=NamespaceHostNetworkTrue PodHost=
ERRO[0000] hostIPC is set to 'true' in PodSpec. It should be set to 'false'.  AuditResultName=NamespaceHostIPCTrue PodHost=
ERRO[0000] hostPID is set to 'true' in PodSpec. It should be set to 'false'.  AuditResultName=NamespaceHostPIDTrue PodHost=
WARN[0000] Image tag is missing.                         AuditResultName=ImageTagMissing Container=fakeContainerSC
WARN[0000] Resource limits not set.                      AuditResultName=LimitsNotSet Container=fakeContainerSC
ERRO[0000] runAsNonRoot is not set in container SecurityContext nor the PodSecurityContext. It should be set to 'true' in at least one of the two.  AuditResultName=RunAsNonRootPSCNilCSCNil Container=fakeContainerSC
ERRO[0000] allowPrivilegeEscalation not set which allows privilege escalation. It should be set to 'false'.  AuditResultName=AllowPrivilegeEscalationNil Container=fakeContainerSC
WARN[0000] privileged is not set in container SecurityContext. Privileged defaults to 'false' but it should be explicitly set to 'false'.  AuditResultName=PrivilegedNil Container=fakeContainerSC
ERRO[0000] readOnlyRootFilesystem is not set in container SecurityContext. It should be set to 'true'.  AuditResultName=ReadOnlyRootFilesystemNil Container=fakeContainerSC
ERRO[0000] Seccomp annotation is missing. The annotation seccomp.security.alpha.kubernetes.io/pod: runtime/default should be added.  AuditResultName=SeccompAnnotationMissing MissingAnnotation=seccomp.security.alpha.kubernetes.io/pod
```

### Example with Kubeaudit Config

Consider the following kubeaudit config `config.yaml`
```yaml
enabledAuditors:
    # Auditors are enabled by default if they are not explicitly set to "false"
    hostns: false
    image: false
    limits: false
auditors:
    capabilities:
        drop: ["AUDIT_WRITE", "CHOWN"]
```

The config can be passed to the `all` command using the `-k/--kconfig` flag:
```
$ kubeaudit all -k "config.yaml" -f "auditors/all/fixtures/audit_all_v1.yml"
ERRO[0000] allowPrivilegeEscalation not set which allows privilege escalation. It should be set to 'false'.  AuditResultName=AllowPrivilegeEscalationNil Container=fakeContainerSC
WARN[0000] privileged is not set in container SecurityContext. Privileged defaults to 'false' but it should be explicitly set to 'false'.  AuditResultName=PrivilegedNil Container=fakeContainerSC
ERRO[0000] Default service account with token mounted. automountServiceAccountToken should be set to 'false' or a non-default service account should be used.  AuditResultName=AutomountServiceAccountTokenTrueAndDefaultSA
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=AUDIT_WRITE Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=CHOWN Container=fakeContainerSC
ERRO[0000] runAsNonRoot is not set in container SecurityContext nor the PodSecurityContext. It should be set to 'true' in at least one of the two.  AuditResultName=RunAsNonRootPSCNilCSCNil Container=fakeContainerSC
ERRO[0000] AppArmor annotation missing. The annotation 'container.apparmor.security.beta.kubernetes.io/fakeContainerSC' should be added.  AuditResultName=AppArmorAnnotationMissing Container=fakeContainerSC MissingAnnotation=container.apparmor.security.beta.kubernetes.io/fakeContainerSC
ERRO[0000] readOnlyRootFilesystem is not set in container SecurityContext. It should be set to 'true'.  AuditResultName=ReadOnlyRootFilesystemNil Container=fakeContainerSC
ERRO[0000] Seccomp annotation is missing. The annotation seccomp.security.alpha.kubernetes.io/pod: runtime/default should be added.  AuditResultName=SeccompAnnotationMissing MissingAnnotation=seccomp.security.alpha.kubernetes.io/pod
ERRO[0000] Namespace is missing a default deny ingress and egress NetworkPolicy.  AuditResultName=MissingDefaultDenyIngressAndEgressNetworkPolicy Namespace=fakeDeploymentSC
```

### Example with Flags

The behaviour of the `all` command can also be customized by using flags. The `all` command supports all flags supported by invididual auditors (see the individual [auditor docs](/README.md#auditors) for all the flags). For example, the `caps` auditor supports specifying capabilities to drop with the `--drop/-d` flag so this flag can be used with the `all` command:
```
kubeaudit all -f "auditors/all/fixtures/audit_all_v1.yml" --drop "AUDIT_WRITE"
```

### Example with Kubeaudit Config and Flags

Passing flags in addition to the config will override the corresponding fields from the config. For example, if the capabilities to drop are specified with the `--drop/-d` flag:
```
kubeaudit all -f "auditors/all/fixtures/audit_all_v1.yml" --drop "AUDIT_WRITE"
```

And they are also specified in the Kubeaudit config file:
```yaml
auditors:
    capabilities:
        drop: ["CHOWN", "MKNOD]
```

The capabilities specified by the flag will take precedence over those specified in the config file resulting in only `AUDIT_WRITE` being dropped.

