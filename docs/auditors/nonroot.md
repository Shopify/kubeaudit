# runAsNonRoot Auditor (nonroot)

Finds containers running as root.

## General Usage

```
kubeaudit nonroot [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit nonroot -f "auditors/nonroot/fixtures/run_as_non_root_nil_v1.yml"
ERRO[0000] runAsNonRoot is not set in container SecurityContext nor the PodSecurityContext. It should be set to 'true' in at least one of the two.  AuditResultName=RunAsNonRootPSCNilCSCNil Container=fakeContainerRANR
```

## Explanation

Containers should be run as a non-root user with the minimum required permissions (principle of least privilege).

This can be done by setting `runAsNonRoot` to `true` in either the PodSecurityContext or container SecurityContext. If `runAsNonRoot` is unset in the Container SecurityContext, it will inherit the value of the Pod SecurityContext. If `runAsNonRoot` is unset in the Pod SecurityContext, it defaults to `false` which means it must be explicitly set to `true` in either the Container SecurityContext or the Pod SecurityContext for the `nonroot` audit to pass.

Note that the Container SecurityContext takes precedence over the Pod SecurityContext so setting `runAsNonRoot` to `false` in the Container SecurityContext will always fail the `nonroot` audit unless an [override](#override-errors) is used.

Ideally, `runAsNonRoot` should be set to `true` in the PodSecurityContext:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template: #PodTemplateSpec
    spec: #PodSpec
      securityContext: #PodSecurityContext
        runAsNonRoot: true
      containers:
      - name: myContainer
```

If a container needs to run as root, it should be enabled for that container only in the container's SecurityContext. This will require an override label so kubeaudit knows it is intentional. See [Override Errors](#override-errors).

For more information on pod and container security contexts see https://kubernetes.io/docs/tasks/configure-pod-container/security-context/

## Override Errors

First, see the [Introduction to Override Errors](/README.md#override-errors).

Override identifer: `allow-run-as-root`

Container overrides have the form:
```yaml
container.audit.kubernetes.io/[container name].allow-run-as-root: ""
```

Pod overrides have the form:
```yaml
audit.kubernetes.io/pod.allow-run-as-root: ""
```

Example of resource with `nonroot` overridden for a specific container:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template: #PodTemplateSpec
    metadata:
      labels:
        container.audit.kubernetes.io/myContainer.allow-run-as-root: ""
    spec: #PodSpec
      securityContext: #PodSecurityContext
        runAsNonRoot: true
      containers:
      - name: myContainer
        securityContext: #SecurityContext
          runAsNonRoot: false
      - name: myContainer2
```

Example of resource with `nonroot` overridden for a whole pod:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template: #PodTemplateSpec
    metadata:
      labels:
        audit.kubernetes.io/pod.allow-run-as-root: ""
    spec: #PodSpec
      securityContext: #PodSecurityContext
        runAsNonRoot: true
      containers:
      - name: myContainer
        securityContext: #SecurityContext
          runAsNonRoot: false
      - name: myContainer2
        securityContext: #SecurityContext
          runAsNonRoot: false
```
