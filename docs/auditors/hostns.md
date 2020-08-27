# Host Namespaces Auditor (hostns)

Finds containers that have HostPID, HostIPC or HostNetwork enabled.

## General Usage

```
kubeaudit hostns [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit hostns -f "auditors/hostns/fixtures/namespaces-all-true.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod
    namespace: namespaces-all-true

--------------------------------------------

-- [error] NamespaceHostNetworkTrue
   Message: hostNetwork is set to 'true' in PodSpec. It should be set to 'false'.

-- [error] NamespaceHostIPCTrue
   Message: hostIPC is set to 'true' in PodSpec. It should be set to 'false'.

-- [error] NamespaceHostPIDTrue
   Message: hostPID is set to 'true' in PodSpec. It should be set to 'false'.
```

## Explanation

**HostPID** - Controls whether the pod containers can share the host process ID namespace. Note that when paired with ptrace this can be used to escalate privileges outside of the container (ptrace is forbidden by default).

**HostIPC** - Controls whether the pod containers can share the host IPC namespace.

**HostNetwork** - Controls whether the pod may use the node network namespace. Doing so gives the pod access to the loopback device, services listening on localhost, and could be used to snoop on network activity of other pods on the same node.

All host namespaces should be disabled unless they are needed. They default to `false` so removing them is sufficient to pass the `hostns` audit, though they can also be explicitly set to `false` if desired.

Example of a resource which **fails** the `hostns` audit:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      hostPID: true
      hostIPC: true
      hostNetwork: true
      containers:
      - name: myContainer
```

For more information on host namespaces, see https://kubernetes.io/docs/concepts/policy/pod-security-policy/#host-namespaces

## Override Errors

First, see the [Introduction to Override Errors](/README.md#override-errors).

Each host namespace field can be individually overridden using their respective override identifiers:
| Host Namespace | Override Identifier |
| :------------- | :--------------------- |
| HostPID | `allow-namespace-host-PID` |
| HostIPC | `allow-namespace-host-IPC` |
| HostNetwork | `allow-namespace-host-network` |

Container overrides have the form:
```yaml
container.audit.kubernetes.io/[container name].[override identifier]: ""
```

Pod overrides have the form:
```yaml
audit.kubernetes.io/pod.[override identifier]: ""
```

Example of a resource with `HostPID` overridden for a specific container:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        container.audit.kubernetes.io/myContainer.allow-namespace-host-PID: ""
    spec:
      hostPID: true
      containers:
      - name: myContainer
```

Example of a resource with `HostPID` overridden for a whole pod:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        audit.kubernetes.io/pod.allow-namespace-host-PID: ""
    spec:
      hostPID: true
      containers:
      - name: myContainer
```
