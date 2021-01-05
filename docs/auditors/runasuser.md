# runAsUser Auditor (runasuser)

Finds containers not overriding the default user with a non-root user.

## General Usage

```
kubeaudit runasuser [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit runasuser -f "auditors/runasuser/fixtures/run-as-user-nil.yml"

---------------- Results for ---------------

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: deployment
    namespace: run-as-user-nil

--------------------------------------------

-- [error] RunAsUserPSCNilCSCNil
   Message: runAsUser is not set in container SecurityContext nor the PodSecurityContext. It should be set to > 0 in at least one of the two.
   Metadata:
      Container: container

## Explanation

Containers should be run as a non-root user with the minimum required permissions (principle of least privilege). By default, Kubernetes will run the container with the user ID used in the container image. Most images use the default user which is root. Therefore, it's good practices to override the container's user ID.

This can be done by setting `runAsUser` to a non-root UID (any UID > 0) in either the PodSecurityContext or container SecurityContext. If `runAsUser` is unset in the Container SecurityContext, it will inherit the value of the Pod SecurityContext. If `runAsUser` is unset in the Pod SecurityContext, it defaults to the image user ID which is usually root which means it must be explicitly set to a non-root UID in either the Container SecurityContext or the Pod SecurityContext for the `runasuser` audit to pass.

Note that the Container SecurityContext takes precedence over the Pod SecurityContext so setting `runAsUser` to the root ID 0 in the Container SecurityContext will always fail the `runasuser` audit unless an [override](#override-errors) is used.

Ideally, `runAsUser` should be set to a non-root UID in the PodSecurityContext:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template: #PodTemplateSpec
    spec: #PodSpec
      securityContext: #PodSecurityContext
        runAsUser: 10000
      containers:
      - name: myContainer
```

If a container needs to run with the user ID set in the image, it should be enabled for that container only in the container's SecurityContext. This will require an override label so kubeaudit knows it is intentional. See [Override Errors](#override-errors).

For more information on pod and container security contexts see https://kubernetes.io/docs/tasks/configure-pod-container/security-context/

## Override Errors

First, see the [Introduction to Override Errors](/README.md#override-errors).

Override identifier: `allow-not-overridden-non-root-user`

Container overrides have the form:
```yaml
container.audit.kubernetes.io/[container name].allow-not-overridden-non-root-user: ""
```

Pod overrides have the form:
```yaml
audit.kubernetes.io/pod.allow-not-overridden-non-root-user: ""
```

Example of resource with `runasuser` overridden for a specific container:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template: #PodTemplateSpec
    metadata:
      labels:
        container.audit.kubernetes.io/myContainer.allow-not-overridden-non-root-user: ""
    spec: #PodSpec
      securityContext: #PodSecurityContext
        runAsUser: 10000
      containers:
      - name: myContainer
        securityContext: #SecurityContext
          runAsUser: 0
      - name: myContainer2
```

Example of resource with `runasuser` overridden for a whole pod:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template: #PodTemplateSpec
    metadata:
      labels:
        audit.kubernetes.io/pod.allow-not-overridden-non-root-user: ""
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
