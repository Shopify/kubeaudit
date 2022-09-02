# Seccomp Auditor (seccomp)

Finds containers running without Seccomp.

## General Usage

```
kubeaudit seccomp [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit seccomp -f "auditors/seccomp/fixtures/seccomp-profile-missing.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod
    namespace: seccomp-profile-missing

--------------------------------------------

-- [error] SeccompProfileMissing
   Message: Pod Seccomp profile is missing. Seccomp profile should be added to the pod SecurityContext.
```

## Explanation

Seccomp (Secure computing mode) is a Linux kernel feature.

Seccomp is enabled by adding a seccomp profile to the security context. The seccomp profile can be either added to a pod security context, which enables seccomp for all containers within that pod, or a security context, which enables seccomp only for that container.

The seccomp profile added to a pod security context has the following format:
```
spec:
  securityContext:
    seccompProfile:
      type: [seccomp profile]
```

The seccomp profile added to a container security context has the following format:
```
spec:
  containers:
    - name: [container name]
      image: [container image]
      securityContext:
        seccompProfile:
          type: [seccomp profile]
```

Ideally, the pod security context should be used.

The value of the seccomp profile type can be set to either the default profile (`RuntimeDefault`) or a custom profile (`Localhost`). For `Localhost` type `localhostProfile: [profile file]` should be added.

Example of a resource which passes the `seccomp` audit:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      securityContext:
        seccompProfile:
          type: RuntimeDefault
      containers:
      - name: myContainer
```

To learn more about Seccomp, see https://en.wikipedia.org/wiki/Seccomp

To learn more about Seccomp in Kubernetes, see https://kubernetes.io/docs/tutorials/security/seccomp/

## Override Errors

Overrides are not currently supported for `seccomp`.
