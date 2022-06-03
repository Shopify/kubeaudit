# Capabilities Auditor (capabilities)

Finds containers that do not drop the recommended capabilities or add new ones.

## General Usage

```
kubeaudit capabilities [flags]
```

### Flags

| Flag  | Description                                                                         |
| :---- | :---------------------------------------------------------------------------------- |
| --allow-add-list | Comma separated list of added capabilities that can be ignored by kubeaudit reports |

Also see [Global Flags](/README.md#global-flags)

## Examples

```shell
$ kubeaudit capabilities -f "auditors/capabilities/fixtures/capabilities-nil.yml"

---------------- Results for ---------------

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: deployment
    namespace: capabilities-nil

--------------------------------------------

-- [error] CapabilityOrSecurityContextMissing
   Message: Security Context not set. The Security Context should be specified and all Capabilities should be dropped by setting the Drop list to ALL.
   Metadata:
      Container: container
```

### Example with Config File

A custom Add list can be provided in the config file. See [docs](docs/all.md) for more information. These are the capabilities you'd like to add and not have kubeaudit raise an error. In this example, kubeaudit will only error for "CHOWN" because it wasn't added to the add list in the config.

`config.yaml`

```yaml
---
auditors:
  capabilities:
    # add capabilities needed to the add list, so kubeaudit won't report errors
    allowAddList: ['KILL', 'MKNOD']
```

`manifest.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
  namespace: example-namespace
spec:
  template:
    spec:
      containers:
        - name: container1
          image: scratch
          securityContext:
            capabilities:
              add:
                - CHOWN
                - KILL
                - MKNOD
              drop:
                - ALL
```

```shell
$ kubeaudit all --kconfig "config.yaml" -f "manifest.yaml"

---------------- Results for ---------------

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: deployment
    namespace: capabilities-some-allowed-multi-containers-some-labels

--------------------------------------------

-- [error] CapabilityAdded
   Message: Capability "CHOWN" added. It should be removed from the capability add list. If you need this capability, add an override label such as'container.audit.kubernetes.io/container1.allow-capability-chown: SomeReason'.
   Metadata:
      Container: container1
```

**Note**: if using http://man7.org/linux/man-pages/man7/capabilities.7.html as a reference for capability names, drop the `CAP_` prefix.

### Example with Custom Add List

A custom add list can be provided as a comma separated value list of capabilities using the `--allow-add-list` flag. These are the capabilities you'd like to add and not have kubeaudit raise an error:

`manifest.yaml` (example manifest)

```yaml
capabilities:
  add:
    - CHOWN
    - KILL
    - MKNOD
    - NET_ADMIN
```

Here we're only adding 3 capabilities to the add list to be ignored. Since we didn't add `NET_ADMIN` to the list, kubeaudit will raise an error for this one.

```shell
  $ kubeaudit capabilities --allow-add-list "CHOWN,KILL,MKNOD" -f "manifest.yaml"
  ---------------- Results for ---------------

  apiVersion: apps/v1beta2
  kind: Deployment
  metadata:
    name: deployment
    namespace: example-namespace

--------------------------------------------

-- [error] CapabilityAdded
   Message: Capability "NET_ADMIN" added. It should be removed from the capability add list. If you need this capability, add an override label such as 'container.audit.kubernetes.io/container1.allow-capability-net-admin: SomeReason'.
   Metadata:
      Container: container1
      Capabiliy: NET_ADMIN

exit status 2
```

## Explanation

Capabilities (specifically, Linux capabilities), are used for permission management in Linux. Some capabilities are enabled by default.

Ideally, all capabilities should be dropped:

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: myContainer
          securityContext:
            capabilities:
              drop:
                - ALL
```

If capabiltiies are required, only those required capabilities should be added:

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: myContainer
          securityContext:
            capabilities:
              drop:
                - all
              add:
                - AUDIT_WRITE
```

In this case, an override label needs to be added to tell kubeaudit that the capability was added on purpose. See [Override Errors](#override-errors).

To learn more about capabilities, see http://man7.org/linux/man-pages/man7/capabilities.7.html

To learn more about capabilities in Kubernetes, see https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-capabilities-for-a-container

## Override Errors

First, see the [Introduction to Override Errors](/README.md#override-errors).

The override identifier has the format `allow-capability-[capability]` which allows for each capability to be individually overridden. To turn a capability name into an override identifier do the following:

1. Lowercase the capability name
1. Replace underscores (`_`) with dashes (`-`)
1. Prepend `allow-capability-`

For example, the override identifier for the `AUDIT_WRITE` capability would be `allow-capability-audit-write`.

Container overrides have the form:

```yaml
container.audit.kubernetes.io/[container name].[override identifier]: ''
```

Pod overrides have the form:

```yaml
audit.kubernetes.io/pod.[override identifier]: ''
```

Example of a resource with `AUDIT_WRITE` and `DAC_OVERRIDE` capabilities overridden for a specific container:

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        container.audit.kubernetes.io/myContainer.allow-capability-audit-write: ''
        container.audit.kubernetes.io/myContainer.allow-capability-dac-override: ''
    spec:
      containers:
        - name: myContainer
          securityContext:
            capabilities:
              drop:
                - ALL
              add:
                - AUDIT_WRITE
                - DAC_OVERRIDE
```

Example of a resource with `AUDIT_WRITE` and `DAC_OVERRIDE` capabilities overridden for a whole pod:

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        audit.kubernetes.io/pod.allow-capability-audit-write: ''
        audit.kubernetes.io/pod.allow-capability-dac-override: ''
    spec:
      containers:
        - name: myContainer
          securityContext:
            capabilities:
              drop:
                - ALL
              add:
                - AUDIT_WRITE
                - DAC_OVERRIDE
```
