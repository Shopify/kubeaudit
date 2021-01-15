# All Auditor (all)

Runs all available auditors, or those specified using a kubeaudit config.

## General Usage

```
kubeaudit all [flags]
```

## Flags

| Short | Long      | Description              | Default |
| :---- | :-------- | :----------------------- | :------ |
| -k    | --kconfig | Path to kubeaudit config |         |

Also see [Global Flags](/README.md#global-flags)

### Kubeaudit Config

A kubeaudit config file can be used instead of flags.

```
kubeaudit all -k "/path/to/kubeaudit-config.yml" -f "/path/to/manifest.yml"
```

Also see [Configuration File](/README.md#configuration-file)

## Examples

```
$ kubeaudit all -f "internal/test/fixtures/all_resources/deployment-apps-v1.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Namespace
  metadata:
    name: deployment-apps-v1

--------------------------------------------

-- [error] MissingDefaultDenyIngressAndEgressNetworkPolicy
   Message: Namespace is missing a default deny ingress and egress NetworkPolicy.
   Metadata:
      Namespace: deployment-apps-v1


---------------- Results for ---------------

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: deployment
    namespace: deployment-apps-v1

--------------------------------------------

-- [error] AppArmorAnnotationMissing
   Message: AppArmor annotation missing. The annotation 'container.apparmor.security.beta.kubernetes.io/container' should be added.
   Metadata:
      Container: container
      MissingAnnotation: container.apparmor.security.beta.kubernetes.io/container

-- [error] AutomountServiceAccountTokenTrueAndDefaultSA
   Message: Default service account with token mounted. automountServiceAccountToken should be set to 'false' or a non-default service account should be used.

-- [error] CapabilityOrSecurityContextMissing
   Message: Security Context not set. The Security Context should be specified and all Capabilities should be dropped by setting the Drop list to ALL.
   Metadata:
      Container: container

-- [error] NamespaceHostNetworkTrue
   Message: hostNetwork is set to 'true' in PodSpec. It should be set to 'false'.
   Metadata:
      PodHost:

-- [error] NamespaceHostIPCTrue
   Message: hostIPC is set to 'true' in PodSpec. It should be set to 'false'.
   Metadata:
      PodHost:

-- [error] NamespaceHostPIDTrue
   Message: hostPID is set to 'true' in PodSpec. It should be set to 'false'.
   Metadata:
      PodHost:

-- [warning] ImageTagMissing
   Message: Image tag is missing.
   Metadata:
      Container: container

-- [warning] LimitsNotSet
   Message: Resource limits not set.
   Metadata:
      Container: container

-- [error] RunAsNonRootPSCNilCSCNil
   Message: runAsNonRoot is not set in container SecurityContext nor the PodSecurityContext. It should be set to 'true' in at least one of the two.
   Metadata:
      Container: container

-- [error] AllowPrivilegeEscalationNil
   Message: allowPrivilegeEscalation not set which allows privilege escalation. It should be set to 'false'.
   Metadata:
      Container: container

-- [warning] PrivilegedNil
   Message: privileged is not set in container SecurityContext. Privileged defaults to 'false' but it should be explicitly set to 'false'.
   Metadata:
      Container: container

-- [error] ReadOnlyRootFilesystemNil
   Message: readOnlyRootFilesystem is not set in container SecurityContext. It should be set to 'true'.
   Metadata:
      Container: container

-- [error] SeccompAnnotationMissing
   Message: Seccomp annotation is missing. The annotation seccomp.security.alpha.kubernetes.io/pod: runtime/default should be added.
   Metadata:
      MissingAnnotation: seccomp.security.alpha.kubernetes.io/pod
```

### Example with Kubeaudit Config

Consider the following kubeaudit config `config.yaml`

```yaml
enabledAuditors:
  # Auditors are enabled by default if they are not explicitly set to "false"
  hostns: false
  image: false
auditors:
  capabilities:
  add:
    - AUDIT_WRITE
    - CHOWN
```

The config can be passed to the `all` command using the `-k/--kconfig` flag:

```
$ kubeaudit all -k "config.yaml" -f "auditors/all/fixtures/audit_all_v1.yml"
```

### Example with Flags

The behaviour of the `all` command can also be customized by using flags. The `all` command supports all flags supported by individual auditors (see the individual [auditor docs](/README.md#auditors) for all the flags).

For example, we can use the `--memory` flag (supported by the `limits` auditor):

```
kubeaudit all -f "manifest.yml" --memory 200
```

Here, if the memory specified is higher than 200, `kubeaudit` will report that the memory limit was exceeded.
