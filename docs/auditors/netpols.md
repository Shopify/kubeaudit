# Default Deny NetworkPolicies for Namespaces Auditor (netpols)

Finds namespaces that do not have a default-deny network policy.

## General Usage

```
kubeaudit netpols [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit netpols -f "auditors/netpols/fixtures/namespace_missing_default_deny_netpol.yml"
ERRO[0000] Namespace is missing a default deny ingress and egress NetworkPolicy.  AuditResultName=MissingDefaultDenyIngressAndEgressNetworkPolicy Namespace=default
```

## Explanation

Just like with firewall rules, the best practice is to deny all internet traffic by default and explicitly allow expected traffic (that is, whitelist expected traffic rather than blacklist unexpected traffic).

This can be done by creating a Network Policy for each namespace which denies all ingress (incoming) and egress (outgoing) traffic. This Network Policy should have an empty pod selector:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: default
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

To allow traffic to a pod, an additional Network Policy can be created which selects that pod.

For more information on network policies, see https://kubernetes.io/docs/concepts/services-networking/network-policies/

## Override Errors

First, see the [Introduction to Override Errors](/README.md#override-errors).

The `netpols` auditor uses a unique override label type not used by any other auditor because the label applies to a namespace (rather than a container or pod):
```
audit.kubernetes.io/namespace.[override identifier]: ""
```

Deny-all ingress and egress network policies can be individually overridden using their respective override identifiers:

| Traffic Type   | Override Identifier                              |
| :------------- | :----------------------------------------------- |
| Ingress        | `allow-non-default-deny-ingress-network-policy`  |
| Egress         | `allow-non-default-deny-egress-network-policy`   |

The override label is placed directly on the Namespace resource:
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: "default"
  labels:
    audit.kubernetes.io/namespace.allow-non-default-deny-ingress-network-policy: ""
```

### Override Example

Consider this Network Policy which denies all egress traffic in the `my-namespace` namespace, but allows all ingress traffic:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: my-namespace
spec:
  podSelector: {}
  policyTypes:
  - Egress
```

The `netpols` auditor will produce an error because there is no `deny-all` Network Policy for ingress traffic:
```
ERRO[0000] All ingress traffic should be blocked by default for namespace my-namespace.  AuditResultName=MissingDefaultDenyIngressNetworkPolicy Namespace=my-namespace
```

This error can be overridden by adding the `audit.kubernetes.io/namespace.allow-non-default-deny-ingress-network-policy: ""` label to the corresponding Namespace resource:
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: "my-namespace"
  labels:
    audit.kubernetes.io/namespace.allow-non-default-deny-ingress-network-policy: ""
```

The auditor will now produce a warning instead of an error:
```
WARN[0000] All ingress traffic should be blocked by default for namespace my-namespace.  AuditResultName=MissingDefaultDenyIngressNetworkPolicyAllowed Namespace=my-namespace
```
