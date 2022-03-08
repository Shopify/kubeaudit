# Seccomp Auditor (seccomp)

Finds containers running without Seccomp.

## General Usage

```
kubeaudit seccomp [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit seccomp -f "auditors/seccomp/fixtures/seccomp-annotation-missing.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod
    namespace: seccomp-annotation-missing

--------------------------------------------

-- [error] SeccompAnnotationMissing
   Message: Seccomp annotation is missing. The annotation seccomp.security.alpha.kubernetes.io/pod: runtime/default should be added.
   Metadata:
      MissingAnnotation: seccomp.security.alpha.kubernetes.io/pod
```

## Explanation

Seccomp (Secure computing mode) is a Linux kernel feature.

Seccomp is enabled by adding a pod-level annotation. The annotation can be either a pod annotation, which enables seccomp for all containers within that pod, or a container annotation, which enables seccomp only for that container.

The pod annotation has the following format:
```
seccomp.security.alpha.kubernetes.io/pod: [seccomp profile]
```

The container annotation has the following format:
```
container.seccomp.security.alpha.kubernetes.io/[container name]: [seccomp profile]
```

Ideally the pod annotation should be used.

The value of the annotation (the `seccomp profile`) can be set to either the default profile (`runtime/default`) or a custom profile (`localhost/[profile name]`).

Example of a resource which passes the `seccomp` audit:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      annotations:
        seccomp.security.alpha.kubernetes.io/pod: runtime/default
    spec:
      containers:
      - name: myContainer
```

To learn more about Seccomp, see https://en.wikipedia.org/wiki/Seccomp

To learn more about Seccomp in Kubernetes, see https://kubernetes.io/docs/tutorials/security/seccomp/

## Override Errors

Overrides are not currently supported for `seccomp`.
