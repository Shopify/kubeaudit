# Capabilities Auditor (capabilities)

Finds containers that do not drop the recommended capabilities or add new ones.

## General Usage

```
kubeaudit capabilities [flags]
```

### Flags

[Global Flags](/README.md#global-flags)

#### Default drop list

- ALL

All is equivalent to dropping all the following capabilities:

| Capability Name   | Description                                                                               |
| :---------------  |  :--------------------------------------------------------------------------------------  |
| AUDIT_WRITE       |  Write records to kernel auditing log                                                     |
| DAC_OVERRIDE      |  Bypass file read, write, and execute permission checks                                   |
| FOWNER            |  Bypass permission checks on operations that normally require the file system UID of the process to match the UID of the file |
| FSETID            |  Don’t clear set-user-ID and set-group-ID permission bits when a file is modified         |
| KILL              |  Bypass permission checks for sending signals                                             |
| MKNOD             |  Create special files using mknod(2)                                                      |
| NET_BIND_SERVICE  |  Bind a socket to internet domain privileged ports (port numbers less than 1024)          |
| NET_RAW           |  Use RAW and PACKET sockets                                                               |
| SETFCAP           |  Set file capabilities                                                                    |
| SETGID            |  Make arbitrary manipulations of process GIDs and supplementary GID list                  |
| SETPCAP           |  Modify process capabilities.                                                             |
| SETUID            |  Make arbitrary manipulations of process UIDs                                             |
| SYS_CHROOT        |  Use chroot(2), change root directory                                                     |

**Note**: if using http://man7.org/linux/man-pages/man7/capabilities.7.html as a reference for capability names, drop the `CAP_` prefix.

## Examples

```
$ kubeaudit capabilities -f "auditors/capabilities/fixtures/capabilities-nil.yml"

---------------- Results for ---------------

  apiVersion: apps/v1beta2
  kind: Deployment
  metadata:
    name: deployment
    namespace: capabilities-nil

--------------------------------------------

-- [error] CapabilityShouldDropAll
   Message: Security Context not set. Ideally, the Security Context should be specified. All capacities should be dropped by setting drop to ALL.
   Metadata:
      Container: container
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
container.audit.kubernetes.io/[container name].[override identifier]: ""
```

Pod overrides have the form:
```yaml
audit.kubernetes.io/pod.[override identifier]: ""
```

Example of a resource with `AUDIT_WRITE` and `DAC_OVERRIDE` capabilities overridden for a specific container:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        container.audit.kubernetes.io/myContainer.allow-capability-audit-write: ""
        container.audit.kubernetes.io/myContainer.allow-capability-dac-override: ""
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
        audit.kubernetes.io/pod.allow-capability-audit-write: ""
        audit.kubernetes.io/pod.allow-capability-dac-override: ""
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
