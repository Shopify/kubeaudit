# Privileged Auditor (privileged)

Finds containers running as privileged.

## General Usage

```
kubeaudit privileged [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit privileged -f "auditors/privileged/fixtures/privileged-true.yml"

---------------- Results for ---------------

  apiVersion: apps/v1
  kind: DaemonSet
  metadata:
    name: daemonset
    namespace: privileged-true

--------------------------------------------

-- [error] PrivilegedTrue
   Message: privileged is set to 'true' in container SecurityContext. It should be set to 'false'.
   Metadata:
      Container: container
```

## Explanation

Running a container as privileged gives all capabilities to the container, and it also lifts all the limitations enforced by the device cgroup controller. In other words, the container can then do almost everything that the host can do. This option exists to allow special use-cases, like running Docker within Docker, but should not be used in most cases.

To prevent a container from running as privileged, `privileged` should be set to `false` in the container SecurityContext. The field defaults to `false` so omitting the field is sufficient to pass the `privileged` audit:

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: myContainer
        securityContext:
          privileged: false
```

For more information on pod and container security contexts see https://kubernetes.io/docs/tasks/configure-pod-container/security-context/

## Override Errors

First, see the [Introduction to Override Errors](/README.md#override-errors).

Override identifier: `allow-privileged`

Container overrides have the form:
```yaml
container.audit.kubernetes.io/[container name].allow-privileged: ""
```

Pod overrides have the form:
```yaml
audit.kubernetes.io/pod.allow-privileged: ""
```

Example of resource with `privileged` overridden for a specific container:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        container.audit.kubernetes.io/myContainer.allow-privilege-escalation: ""
    spec:
      containers:
      - name: myContainer
        securityContext:
          privileged: true
```

Example of resource with `privileged` overridden for a whole pod:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        audit.kubernetes.io/pod.allow-privileged: ""
    spec:
      containers:
      - name: myContainer
        securityContext:
          privileged: true
```
