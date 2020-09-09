# automountServiceAccountToken Auditor (asat)

Finds containers that meet either of the following conditions:
1. The deprecated `serviceAccount` field is used 
1. The default service account is automatically mounted

## General Usage

```
kubeaudit asat [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples
```
kubeaudit asat -f "auditors/asat/fixtures/service-account-token-true-and-no-name.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: ReplicationController
  metadata:
    name: replicationcontroller
    namespace: service-account-token-true-and-no-name

--------------------------------------------

-- [error] AutomountServiceAccountTokenTrueAndDefaultSA
   Message: Default service account with token mounted. automountServiceAccountToken should be set to 'false' on either the ServiceAccount or on the PodSpec or a non-default service account should be used.
```

## Explanation

`serviceAccount` is a deprecated field. `serviceAccountName` should be used instead.

Example of a resource which fails the `asat` check because it uses `serviceAccount`:
```yaml
apiVersion: v1
kind: Deployment
spec:
  template:
    spec:
      serviceAccount: ThisFieldIsDeprecated
      containers:
      - name: myContainer
```

Automounting a default service account would allow any compromised pod to run API commands against the cluster. Either automounting should be disabled or a non-default service account with sane permissions should be used.

To make sure a non-default service account is used, `serviceAccountName` must be set to a value other than `default`.

To make sure a service account is not automatically mounted, `automountServiceAccountToken` must be explicitly set to `false` (it defaults to `true`) on either the ServiceAccount (for kubernetes 1.6+) or on the PodSpec.

Example of disabling `automountServiceAccountToken` on the default ServiceAccount (kubernetes 1.6+):
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: default
automountServiceAccountToken: false
```

Example of disabling `automountServiceAccountToken` on the PodSpec of a Deployment:
```yaml
apiVersion: v1
kind: Deployment
spec:
  template:
    spec:
      automountServiceAccountToken: false
      containers:
      - name: myContainer
```

Note that if `automountServiceAccountToken` is set on the PodSpec, this will take precedence over `automountServiceAccountToken` set on the ServiceAccount, so you should never set `automountServiceAccountToken: true` in the PodSpec when using the default ServiceAccount.

Example of using a non-default service account:
```yaml
apiVersion: v1
kind: Deployment
spec:
  template:
    spec:
      serviceAccountName: customServiceAccount
      containers:
      - name: myContainer
```

## Override Errors

First, see the [Introduction to Override Errors](/README.md#override-errors).

Override identifier: `allow-automount-service-account-token`

Only pod overrides are supported:
```yaml
audit.kubernetes.io/pod.allow-automount-service-account-token: ""
```

Example of a resource with `asat` results overridden:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      labels:
        audit.kubernetes.io/pod.allow-automount-service-account-token: ""
    spec:
      automountServiceAccountToken: true
      containers:
      - name: myContainer
```
