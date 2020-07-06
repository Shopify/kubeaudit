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
kubeaudit asat -f "auditors/asat/fixtures/service_account_token_true_and_no_name_v1.yml"
ERRO[0000] Default service account with token mounted. automountServiceAccountToken should be set to 'false' or a non-default service account should be used.  AuditResultName=AutomountServiceAccountTokenTrueAndDefaultSA
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

To make sure a service account is not automatically mounted, `automountServiceAccountToken` must be explicitly set to `false` (it defaults to `true`).

Example of a resource which passes the `asat` audit:
```yaml
apiVersion: v1
kind: Deployment
spec:
  template:
    spec:
      serviceAccountName: myServiceAccount
      automountServiceAccountToken: false
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
