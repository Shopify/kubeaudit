# Sensitive Host Path Mounted Auditor (mounts)

Finds containers that have sensitive host paths mounted. 

## General Usage

```
kubeaudit mounts [flags]
```

### Flags

| Short   | Long      | Description                                                          | Default                                                                  |
| :------ | :-------- | :------------------------------------------------------------------- | :----------------------------------------------------------------------- |
| -s      | --paths   | List of sensitive paths that shouldn't be mounted.                   | [default sensitive host paths list](#Default-sensitive-host-paths-list)  |

Also see [Global Flags](/README.md#global-flags)

#### Default sensitive host paths list

| Host path                 | Description                                                                              |
| :------------------------ | :--------------------------------------------------------------------------------------- |
| /proc                     |  Pseudo-filesystem which provides an interface to kernel data structures                                                  |
| /var/run/docker.sock      |  Unix socket used to communicate with Docker daemon                                   |
| /                         |  Filesystem's root |
| /etc                      |  Directory that usually contains all system related configurations files         |
| /root                     |  Home directory of the `root` user                                             |
| /var/run/crio/crio.sock   |  Unix socket used to communicate with the CRI-O Container Engine                                                 |
| /home/admin               |  Home directory of the `admin` user        |
| /var/lib/kubelet          |  Directory for Kublet-related configuration                                                             |
| /var/lib/kubelet/pki      |  Directory containing the certificate and private key of the Kublet                                                               |
| /etc/kubernetes           |  Directory containing Kubernetes related configuration              |
| /etc/kubernetes/manifests |  Directory containing manifest of Kubernetes components                                                         |


## Examples

```
$ kubeaudit mounts -f /Users/j.courtial/dev/go/src/github.com/kubeaudit/auditors/mounts/fixtures/proc-mounted.yml

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod
    namespace: proc-mounted

--------------------------------------------

-- [error] SensitivePathsMounted
   Message: Sensitive path mounted as volume: proc-volume (/proc -> /host/proc, readOnly: false). It should be removed from the container's mounts list.
   Metadata:
      Container: container
      Mount: proc-volume

```

## Explanation

Mounting some sensitive host paths (like `/etc`, `/proc`, or `/var/run/docker.sock`) may allow a container to access sensitive information from the host like credentials or to spy on other workloads' activity.

These sensitive paths should not be mounted.

Example of a resource which **fails** the `mounts` audit:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: container
          image: scratch
          volumeMounts:
            - mountPath: /host/proc
              name: proc-volume
      volumes:
        - name: proc-volume
          hostPath:
            path: /proc
```

## Override Errors

First, see the [Introduction to Override Errors](/README.md#override-errors).

The override identifier has the format `allow-host-path-mount--[mount name]` which allows for each mount to be individually overridden.

Example of resource with `mounts` overridden for a specific container:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template: #PodTemplateSpec
    metadata:
      labels:
        container.audit.kubernetes.io/container2.allow-host-path-mount-proc-volume: "SomeReason"
    spec: #PodSpec
      containers:
        - name: container1
          image: scratch
        - name: container2
          image: scratch
          volumeMounts:
            - mountPath: /host/proc
              name: proc-volume
      volumes:
        - name: proc-volume
          hostPath:
            path: /proc
```

Example of resource with `mounts` overridden for a whole pod:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template: #PodTemplateSpec
    metadata:
      labels:
        audit.kubernetes.io/pod.allow-host-path-mount-proc-volume: "SomeReason"
    spec: #PodSpec
      containers:
        - name: container1
          image: scratch
          volumeMounts:
            - mountPath: /host/proc
              name: proc-volume
        - name: container2
          image: scratch
          volumeMounts:
            - mountPath: /host/proc
              name: proc-volume
      volumes:
        - name: proc-volume
          hostPath:
            path: /proc
```
