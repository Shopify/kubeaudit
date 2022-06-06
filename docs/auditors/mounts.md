# Sensitive Host Path Mounted Auditor (mounts)

Finds containers that have sensitive host paths mounted.

## General Usage

```
kubeaudit mounts [flags]
```

### Flags

| Short   | Long              | Description                                                          | Default                                                                  |
| :------ | :---------------- | :------------------------------------------------------------------- | :----------------------------------------------------------------------- |
| -d      | --denyPathsList   | List of sensitive paths that shouldn't be mounted.                   | [default sensitive host paths list](#Default-sensitive-host-paths-list)  |

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
| /var/lib/kubelet/pki      |  Directory containing the certificate and private key of the kublet                                                               |
| /etc/kubernetes           |  Directory containing Kubernetes related configuration              |
| /etc/kubernetes/manifests |  Directory containing manifest of Kubernetes components                                                         |

## Examples

```
$ kubeaudit mounts -f auditors/mounts/fixtures/proc-mounted.yml

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
      MountName: proc-volume
      MountPath: /host/proc
      MountReadOnly: false
      MountVolume: proc-volume
      MountVolumeHostPath: /proc

```

### Example with Config File

If you don't want kubeaudit to raise errors for all the paths in the default list (`DefaultSensitivePaths`), you can
provide a custom paths list in the config file. See [docs](docs/all.md) for more information. That way kubeaudit will
only raise errors for those specific paths listed in the config file.

`config.yaml`

```yaml
---
enabledAuditors:
  mounts: true
auditors:
  mounts:
    denyPathsList: ["/etc", "/var/run/docker.sock"]
```

`manifest.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
  namespace: example-namespace
spec:
  template:
    spec:
      containers:
        - name: container
          image: scratch
          volumeMounts:
            - mountPath: /host/etc
              name: etc-volume
            - mountPath: /var/run/docker.sock
              name: docker-socket-volume
      volumes:
        - name: etc-volume
          hostPath:
            path: /etc
        - name: docker-socket-volume
          hostPath:
            path: /var/run/docker.sock
```

```shell
$ kubeaudit all --kconfig "config.yaml" -f "manifest.yaml"

---------------- Results for ---------------

  apiVersion: apps/v1beta2
  kind: Deployment
  metadata:
    name: deployment
    namespace: example-namespace

--------------------------------------------

-- [error] SensitivePathsMounted
   Message: Sensitive path mounted as volume: etc-volume (hostPath: /etc). It should be removed from the container's mounts list.
   Metadata:
      Container: container
      MountName: etc-volume
      MountPath: /host/etc
      MountReadOnly: false
      MountVolume: etc-volume
      MountVolumeHostPath: /etc

-- [error] SensitivePathsMounted
   Message: Sensitive path mounted as volume: docker-socket-volume (hostPath: /var/run/docker.sock). It should be removed from the container's mounts list.
   Metadata:
      MountReadOnly: false
      MountVolume: docker-socket-volume
      MountVolumeHostPath: /var/run/docker.sock
      Container: container
      MountName: docker-socket-volume
      MountPath: /var/run/docker.sock
```

### Example with Custom Paths List

A custom paths list can be provided as a comma separated value list of paths using the `--denyPathsList` flag. These are
the host paths you'd like to have kubeaudit raise an error when they are mounted in a container.

`manifest.yaml` (example manifest)

```yaml
volumes:
  - name: etc-volume
    hostPath:
      path: /etc
  - name: docker-socket-volume
    hostPath:
      path: /var/run/docker.sock
```

```shell
$ kubeaudit mounts --denyPathsList "/etc,/var/run/docker.sock" -f "manifest.yaml"
---------------- Results for ---------------

  apiVersion: apps/v1beta2
  kind: Deployment
  metadata:
    name: deployment
    namespace: example-namespace

--------------------------------------------

-- [error] SensitivePathsMounted
   Message: Sensitive path mounted as volume: etc-volume (hostPath: /etc). It should be removed from the container's mounts list.
   Metadata:
      Container: container
      MountName: etc-volume
      MountPath: /host/etc
      MountReadOnly: false
      MountVolume: etc-volume
      MountVolumeHostPath: /etc

-- [error] SensitivePathsMounted
   Message: Sensitive path mounted as volume: docker-socket-volume (hostPath: /var/run/docker.sock). It should be removed from the container's mounts list.
   Metadata:
      Container: container
      MountName: docker-socket-volume
      MountPath: /var/run/docker.sock
      MountReadOnly: false
      MountVolume: docker-socket-volume
      MountVolumeHostPath: /var/run/docker.sock
```

## Explanation

Mounting some sensitive host paths (like `/etc`, `/proc`, or `/var/run/docker.sock`) may allow a container to access
sensitive information from the host like credentials or to spy on other workloads' activity.

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

The override identifier has the format `allow-host-path-mount-[mount name]` which allows for each mount to be
individually overridden.

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
