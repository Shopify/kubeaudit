# AppArmor Auditor (apparmor)

Finds containers that do not have AppArmor enabled.

## General Usage

```
kubeaudit apparmor [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit apparmor -f "auditors/apparmor/fixtures/apparmor_annotation_missing_v1.yml"
ERRO[0000] AppArmor annotation missing. The annotation 'container.apparmor.security.beta.kubernetes.io/AAcontainer' should be added.  AuditResultName=AppArmorAnnotationMissing Container=AAcontainer MissingAnnotation=container.apparmor.security.beta.kubernetes.io/AAcontainer
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
