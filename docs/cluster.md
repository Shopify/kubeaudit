# Running kubeaudit in a Cluster

Kubeaudit can be run in a Kubernetes cluster by using the [official Docker image](https://hub.docker.com/r/shopify/kubeaudit): `shopify/kubeaudit`.

## Without RBAC

Example Job configuration:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubeaudit
  namespace: default

---

apiVersion: batch/v1
kind: Job
metadata:
  name: kubeaudit
  namespace: default
spec:
  template:
    metadata:
      annotations:
        container.apparmor.security.beta.kubernetes.io/kubeaudit: runtime/default
        seccomp.security.alpha.kubernetes.io/pod: runtime/default
    spec:
      serviceAccountName: kubeaudit
      restartPolicy: OnFailure
      containers:
        - name: kubeaudit
          image: shopify/kubeaudit:v0.11
          args: ["all", "--exitcode", "0"]
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop: ["all"]
            privileged: false
            readOnlyRootFilesystem: true
            runAsNonRoot: true
```

## With RBAC

If RBAC is enabled on the cluster:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubeaudit
  namespace: default

---

kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: kubeaudit
rules:
  - apiGroups: [""]
    resources:
      - pods
      - podtemplates
      - replicationcontrollers
      - namespaces
      - serviceaccounts
    verbs: ["list"]
  - apiGroups: ["apps"]
    resources:
      - daemonsets
      - statefulsets
      - deployments
    verbs: ["list"]
  - apiGroups: ["batch"]
    resources:
      - cronjobs
    verbs: ["list"]
  - apiGroups: ["networking.k8s.io"]
    resources:
      - networkpolicies
    verbs: ["list"]

---

kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: kubeaudit
subjects:
  - kind: ServiceAccount
    name: kubeaudit
    namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kubeaudit

---

apiVersion: batch/v1
kind: Job
metadata:
  name: kubeaudit
  namespace: default
spec:
  template:
    metadata:
      annotations:
        container.apparmor.security.beta.kubernetes.io/kubeaudit: runtime/default
        seccomp.security.alpha.kubernetes.io/pod: runtime/default
    spec:
      serviceAccountName: kubeaudit
      restartPolicy: OnFailure
      containers:
        - name: kubeaudit
          image: shopify/kubeaudit:v0.11
          args: ["all", "--exitcode", "0"]
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop: ["all"]
            privileged: false
            readOnlyRootFilesystem: true
            runAsNonRoot: true
```

## With RBAC and a Specific Namespace

If you are running kubeaudit on a specific namespace and don't want to grant it cluster wide access, the binding can be made into a namespaced binding, but note that kubeaudit will still need to be able to list namespaces at the cluster level (as namespace resources don't have a namespaced scope).

In the following example, the `kubeaudit` Job is created in the `kubeaudit` namespace and is assigned a ServiceAccount which can list namespaces at a cluster scope but can only list the other resources for the provided namespace.

**Important**: Replace the two instances of `<TARGET_NAMESPACE>` with the namespace you want kubeaudit to audit:

```yaml
# Optionally, run kubeaudit in its own namespace
apiVersion: v1
kind: Namespace
metadata:
  name: kubeaudit

---

# Don't allow internet traffic in or out of the kubeaudit namespace
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: kubeaudit
spec:
  policyTypes:
  - Ingress
  - Egress

---

apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubeaudit
  namespace: kubeaudit

---

kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: kubeaudit-namespaces
rules:
  - apiGroups: [""]
    resources:
      - namespaces
    verbs: ["list"]

---

kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: kubeaudit
rules:
  - apiGroups: [""]
    resources:
      - pods
      - podtemplates
      - replicationcontrollers
      - serviceaccounts
    verbs: ["list"]
  - apiGroups: ["apps"]
    resources:
      - daemonsets
      - statefulsets
      - deployments
    verbs: ["list"]
  - apiGroups: ["batch"]
    resources:
      - cronjobs
    verbs: ["list"]
  - apiGroups: ["networking.k8s.io"]
    resources:
      - networkpolicies
    verbs: ["list"]

---

kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: kubeaudit-namespaces
subjects:
  - kind: ServiceAccount
    name: kubeaudit
    namespace: kubeaudit
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kubeaudit-namespaces

---

kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: kubeaudit
  namespace: <TARGET_NAMESPACE>
subjects:
  - kind: ServiceAccount
    name: kubeaudit
    namespace: kubeaudit
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kubeaudit

---

apiVersion: batch/v1
kind: Job
metadata:
  name: kubeaudit
  namespace: kubeaudit
spec:
  template:
    metadata:
      annotations:
        container.apparmor.security.beta.kubernetes.io/kubeaudit: runtime/default
        seccomp.security.alpha.kubernetes.io/pod: runtime/default
    spec:
      serviceAccountName: kubeaudit
      restartPolicy: OnFailure
      containers:
        - name: kubeaudit
          image: shopify/kubeaudit:v0.11
          args: ["all", "--exitcode", "0", "--namespace", "<TARGET_NAMESPACE>"]
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop: ["all"]
            privileged: false
            readOnlyRootFilesystem: true
            runAsNonRoot: true
```
