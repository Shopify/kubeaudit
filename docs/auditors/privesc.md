# Privilege Escalation Allowed Auditor (privesc)

Finds containers that allow privilege escalation.

## General Usage

```
kubeaudit privesc [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit privesc -f "auditors/privesc/fixtures/allow-privilege-escalation-nil.yml"

---------------- Results for ---------------

  apiVersion: apps/v1
  kind: StatefulSet
  metadata:
    name: statefulset
    namespace: allow-privilege-escalation-nil

--------------------------------------------

-- [error] AllowPrivilegeEscalationNil
   Message: allowPrivilegeEscalation not set which allows privilege escalation. It should be set to 'false'.
   Metadata:
      Container: container
```

## Explanation

`allowPrivilegeEscalation` controls whether a process can gain more privileges than its parent process.

Privilege escalation is disabled by setting `allowPrivilegeEscalation` to `false` in the container SecurityContext. The field defaults to `true` so it must be explicitly set to `false` to pass the `privesc` audit:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: myContainer
        securityContext:
          allowPrivilegeEscalation: false
```

For more information on pod and container security contexts see https://kubernetes.io/docs/tasks/configure-pod-container/security-context/

## Override Errors

First, see the [Introduction to Override Errors](/README.md#override-errors).

Override identifier: `allow-privilege-escalation`

Container overrides have the form:
```yaml
container.audit.kubernetes.io/[container name].allow-privilege-escalation: ""
```

Pod overrides have the form:
```yaml
audit.kubernetes.io/pod.allow-privilege-escalation: ""
```

Example of resource with `privesc` overridden for a specific container:
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
          allowPrivilegeEscalation: true
```

Example of resource with `privesc` overridden for a whole pod:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        audit.kubernetes.io/pod.allow-privilege-escalation: ""
    spec:
      containers:
      - name: myContainer
        securityContext:
          allowPrivilegeEscalation: true
```
