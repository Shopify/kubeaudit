# Docker Socket Mounted Auditor (mountds)

DEPRECATED. Please use `mounts` instead. This command will be removed in a future minor version.

Finds containers that have the docker socket mounted.

## General Usage

```
kubeaudit mountds [flags]
```

See [Global Flags](/README.md#global-flags)

## Examples

```
$ kubeaudit mountds -f "auditors/mountds/fixtures/docker-sock-mounted.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod
    namespace: docker-sock-mounted

--------------------------------------------

-- [warning] DockerSocketMounted
   Message: Docker socket is mounted. '/var/run/docker.sock' should be removed from the container's volume mount list.
   Metadata:
      Container: container
```

## Explanation

The `/var/run/docker.sock` file is the Unix socket the Docker daemon listens on by default. Mounting this file as a volume allows containers to communicate with the Docker daemon.

The docker socket should not be mounted to prevent compromised containers from controlling the Docker daemon.

Example of a resource which **fails** the `mountds` audit:
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: myContainer
        volumeMounts:
        - mountPath: /var/run/docker.sock
          name: docker-sock-volume
      volumes:
      - name: docker-sock-volume
        hostPath:
          path: /var/run/docker.sock
```

## Override Errors

Overrides are not currently supported for `mountds`.
