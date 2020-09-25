# Autofix (autofix)

Automatically fixes security issues.

**Note**: `autofix` can only be used in manifest mode.

## General Usage

```
kubeaudit autofix -f [manifest] [flags]
```

## Flags

| Short   | Long       | Description                               | Default                                  |
| :------ | :--------- | :---------------------------------------- | :--------------------------------------- |
| -o      | --outfile  | File to write fixed manifest to           |                                          |
| -k      | --kconfig  | Path to kubeaudit config file             |                                          |

Also see [Global Flags](/README.md#global-flags)

## Examples

Consider this simple manifest file `manifest.yml`:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: myContainer
```

The `autofix` command will make the manifest secure!:
```
kubeaudit autofix -f "manifest.yml"
```

Fixed manifest:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: myContainer
        resources: {}
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
          privileged: false
          readOnlyRootFilesystem: true
          runAsNonRoot: true
      automountServiceAccountToken: false
    metadata:
      annotations:
        container.apparmor.security.beta.kubernetes.io/fakeContainerSC: runtime/default
        seccomp.security.alpha.kubernetes.io/pod: runtime/default
  selector: null
  strategy: {}
metadata:
```

### Example with Multiple Resources

The `autofix` command works on manifest files containing multiple resources:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: myContainer

---

apiVersion: v1
kind: Pod
spec:
  containers:
  - name: myContainer2
    image: polinux/stress
```

Fixed manifest:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: myContainer
        resources: {}
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - AUDIT_WRITE
            - CHOWN
            - DAC_OVERRIDE
            - FOWNER
            - FSETID
            - KILL
            - MKNOD
            - NET_BIND_SERVICE
            - NET_RAW
            - SETFCAP
            - SETGID
            - SETPCAP
            - SETUID
            - SYS_CHROOT
          privileged: false
          readOnlyRootFilesystem: true
          runAsNonRoot: true
      automountServiceAccountToken: false
    metadata:
      annotations:
        container.apparmor.security.beta.kubernetes.io/myContainer: runtime/default
        seccomp.security.alpha.kubernetes.io/pod: runtime/default
  selector: null
  strategy: {}
metadata:

---

apiVersion: v1
kind: Pod
spec:
  containers:
  - name: myContainer2
    image: polinux/stress
    resources: {}
    securityContext:
      allowPrivilegeEscalation: false
      capabilities:
        drop:
        - AUDIT_WRITE
        - CHOWN
        - DAC_OVERRIDE
        - FOWNER
        - FSETID
        - KILL
        - MKNOD
        - NET_BIND_SERVICE
        - NET_RAW
        - SETFCAP
        - SETGID
        - SETPCAP
        - SETUID
        - SYS_CHROOT
      privileged: false
      readOnlyRootFilesystem: true
      runAsNonRoot: true
  automountServiceAccountToken: false
metadata:
  annotations:
    container.apparmor.security.beta.kubernetes.io/myContainer2: runtime/default
    seccomp.security.alpha.kubernetes.io/pod: runtime/default
```

### Example with Comments

The `autofix` command supports comments!
```yaml
# This is a sample Kubernetes config file
#
# Autofix supports comments!

%YAML   1.1
%TAG    !   !foo
%TAG    !yaml!  tag:yaml.org,2002:

---

apiVersion: apps/v1
kind: Deployment
# PodSpec
spec:
  # PodTemplate
  template:
    # ContainerSpec
    spec:
      containers:
      - name: myContainer # this is a sample container
```

Fixed manifest:
```yaml
# This is a sample Kubernetes config file
#
# Autofix supports comments!

%YAML   1.1
%TAG    !   !foo
%TAG    !yaml!  tag:yaml.org,2002:

---

apiVersion: apps/v1
kind: Deployment
# PodSpec
spec:
  # PodTemplate
  template:
    # ContainerSpec
    spec:
      containers:
      - name: myContainer # this is a sample container
        resources: {}
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - AUDIT_WRITE
            - CHOWN
            - DAC_OVERRIDE
            - FOWNER
            - FSETID
            - KILL
            - MKNOD
            - NET_BIND_SERVICE
            - NET_RAW
            - SETFCAP
            - SETGID
            - SETPCAP
            - SETUID
            - SYS_CHROOT
          privileged: false
          readOnlyRootFilesystem: true
          runAsNonRoot: true
      automountServiceAccountToken: false
    metadata:
      annotations:
        container.apparmor.security.beta.kubernetes.io/myContainer: runtime/default
        seccomp.security.alpha.kubernetes.io/pod: runtime/default
  selector: null
  strategy: {}
metadata:

```

### Example with Custom Output File

To write the fixed manifest to a different file, use the `--outfile/-o` flag:
```
kubeaudit autofix -f "manifest.yml" -o "fixed.yaml"
```

### Using Custom Rules with Kubeaudit Config File

To fix a manifest based on custom rules specified on a kubeaudit config file (e.g disable some auditors), use the `-k/--kconfig` flag.

```
kubeaudit autofix -k "/path/to/kubeaudit-config.yml" -f "/path/to/manifest.yml" -o "/path/to/fixed"
```

Also see [Configuration File](/README.md#configuration-file)