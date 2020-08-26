# readOnlyRootFilesystem Auditor (rootfs)

Finds containers which do not have a read-only filesystem.

## General Usage

```
kubeaudit rootfs [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit rootfs -f "auditors/rootfs/fixtures/read-only-root-filesystem-nil.yml"

---------------- Results for ---------------

  apiVersion: apps/v1
  kind: StatefulSet
  metadata:
    name: statefulset
    namespace: read-only-root-filesystem-nil

--------------------------------------------

-- [error] ReadOnlyRootFilesystemNil
   Message: readOnlyRootFilesystem is not set in container SecurityContext. It should be set to 'true'.
   Metadata:
      Container: container
```

## Explanation

If a container does not need to write files, it should be run with a read-only filesystem.

To run a container with a read-only filesystem, `readOnlyRootFilesystem` should be set to `true` in the container SecurityContext. The field defaults to `false` so it must be explicitly set to `true` to pass the `rootfs` audit:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: myContainer
        securityContext:
          readOnlyRootFilesystem: true
```

If a container needs to write files, an override label needs to be used so kubeaudit knows it is intentional. See [Override Errors](#override-errors).

For more information on pod and container security contexts see https://kubernetes.io/docs/tasks/configure-pod-container/security-context/

## Override Errors

First, see the [Introduction to Override Errors](/README.md#override-errors).

Override identifier: `allow-read-only-root-filesystem-false`

Container overrides have the form:
```yaml
container.audit.kubernetes.io/[container name].allow-read-only-root-filesystem-false: ""
```

Pod overrides have the form:
```yaml
audit.kubernetes.io/pod.allow-read-only-root-filesystem-false: ""
```

Example of resource with `rootfs` overridden for a specific container:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        container.audit.kubernetes.io/myContainer.allow-read-only-root-filesystem-false: ""
    spec:
      containers:
      - name: myContainer
        securityContext:
          readOnlyRootFilesystem: false
```

Example of resource with `rootfs` overridden for a whole pod:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        audit.kubernetes.io/pod.allow-read-only-root-filesystem-false: ""
    spec:
      containers:
      - name: myContainer
        securityContext:
          readOnlyRootFilesystem: false
```
