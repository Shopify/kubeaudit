# Capabilities Auditor (capabilities)

Finds containers that do not drop the recommended capabilities or add new ones.

## General Usage

```
kubeaudit capabilities [flags]
```

### Flags

| Short   | Long      | Description                                                          | Default                                  |
| :------ | :-------- | :------------------------------------------------------------------- | :--------------------------------------- |
| -d      | --drop    | Space separated list of capabilities that should be dropped.         | [default drop list](#default-drop-list)  |

Also see [Global Flags](/README.md#global-flags)

#### Default drop list

| Capability Name  | Description                                                                              |
| :--------------- | :--------------------------------------------------------------------------------------- |
| AUDIT_WRITE      | Write records to kernel auditing log                                                     |
| DAC_OVERRIDE     | Bypass file read, write, and execute permission checks                                   |
| FOWNER           | Bypass permission checks on operations that normally require the file system UID of the process to match the UID of the file |
| FSETID           | Don’t clear set-user-ID and set-group-ID permission bits when a file is modified         |
| KILL             | Bypass permission checks for sending signals                                             |
| MKNOD             |  Create special files using mknod(2)                                                    |
| NET_BIND_SERVICE  |  Bind a socket to internet domain privileged ports (port numbers less than 1024)        |
| NET_RAW           |  Use RAW and PACKET sockets                                                             |
| SETFCAP           |  Set file capabilities                                                                  |
| SETGID            |  Make arbitrary manipulations of process GIDs and supplementary GID list                |
| SETPCAP           |  Modify process capabilities.                                                           |
| SETUID            |  Make arbitrary manipulations of process UIDs                                           |
| SYS_CHROOT        |  Use chroot(2), change root directory                                                   |

## Examples

```
$ kubeaudit capabilities -f "auditors/capabilities/fixtures/capabilities_nil_v1beta2.yml"
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
```

### Example with Custom Drop List

A custom drop list can be provided as a space-separated list of capabilities using the `-d/--drop` flag:

```
$ kubeaudit capabilities --drop "MAC_ADMIN AUDIT_WRITE" -f "auditors/capabilities/fixtures/capabilities_nil_v1beta2.yml"
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=MAC_ADMIN Container=fakeContainerSC
ERRO[0000] Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.  AuditResultName=CapabilityNotDropped Capability=AUDIT_WRITE Container=fakeContainerSC
```

**Note**: if using http://man7.org/linux/man-pages/man7/capabilities.7.html as a reference for capability names, drop the `CAP_` prefix.

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
            - all
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
            - all
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
            - all
            add:
            - AUDIT_WRITE
            - DAC_OVERRIDE
```
