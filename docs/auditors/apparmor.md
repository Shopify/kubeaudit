# AppArmor Auditor (apparmor)

Finds containers that do not have AppArmor enabled.

## General Usage

```
kubeaudit apparmor [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit apparmor -f "auditors/apparmor/fixtures/apparmor-annotation-missing.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod
    namespace: apparmor-annotation-missing

--------------------------------------------

-- [error] AppArmorAnnotationMissing
   Message: AppArmor annotation missing. The annotation 'container.apparmor.security.beta.kubernetes.io/container' should be added.
   Metadata:
      MissingAnnotation: container.apparmor.security.beta.kubernetes.io/container
      Container: container
```

If an apparmor annotation refers to a container which doesn't exist, `kubectl apply` will fail. Kubeaudit produces an error for this case:

```
$ kubeaudit apparmor -f "auditors/apparmor/fixtures/apparmor-invalid-annotation.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod
    namespace: apparmor-enabled

--------------------------------------------

-- [error] AppArmorInvalidAnnotation
   Message: AppArmor annotation key refers to a container that doesn't exist. Remove the annotation 'container.apparmor.security.beta.kubernetes.io/container2: runtime/default'.
   Metadata:
      Container: container2
      Annotation: container.apparmor.security.beta.kubernetes.io/container2: runtime/default
```

## Explanation

AppArmor is a Mandatory Access Control (MAC) system used by Linux.

AppArmor is enabled by adding `container.apparmor.security.beta.kubernetes.io/[container name]` as a pod-level annotation and setting its value to either `runtime/default` or a profile (`localhost/[profile name]`).

Example of a resource which passes the `apparmor` audit:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      annotations:
        container.apparmor.security.beta.kubernetes.io/myContainer: runtime/default
    spec:
      containers:
      - name: myContainer
```

To learn more about AppArmor, see https://wiki.ubuntu.com/AppArmor

To learn more about AppArmor in Kubernetes, see https://kubernetes.io/docs/tutorials/clusters/apparmor/#securing-a-pod

## Override Errors

Overrides are not currently supported for `apparmor`.
