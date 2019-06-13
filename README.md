[![Build Status](https://api.travis-ci.org/Shopify/kubeaudit.svg?branch=master)](https://travis-ci.org/Shopify/kubeaudit/)
[![codecov](https://codecov.io/gh/Shopify/kubeaudit/branch/master/graph/badge.svg)](https://codecov.io/gh/Shopify/kubeaudit)
[![Go Report Card](https://goreportcard.com/badge/github.com/Shopify/kubeaudit)](https://goreportcard.com/report/github.com/Shopify/kubeaudit)
[![GoDoc](https://godoc.org/github.com/Shopify/kubeaudit?status.png)](https://godoc.org/github.com/Shopify/kubeaudit)
# kubeaudit :cloud: :lock: :muscle:

`kubeaudit` is a command line tool to audit Kubernetes clusters for various
different security concerns: run the container as a non-root user, use a read
only root filesystem, drop scary capabilities, don't add new ones, don't run
privileged, ... You get the gist of it and more on that later. Just know:

## `kubeaudit` makes sure you deploy secure containers!

- [Installation](#installation)
- [General instructions](#general)
- [Autofix](#autofix)
- [Audits](#audits)
- [Override Labels](#labels)
- [Audit Configuration](#audit-configuration)
- [Contribute!](#contribute)

<a name="installation" />

## Installation

#### Download a binary

Kubeaudit has official releases that are blessed and stable here:
[Official releases](https://github.com/Shopify/kubeaudit/releases)

#### DIY build

Master will have newer features than the stable releases. If you need a newer
feature not yet included in a release you can do the following to get
kubeaudit:

```sh
go get -v github.com/Shopify/kubeaudit
make
make install
```

Now you can just call `kubeaudit` with one of commands from [here](#audits)

#### Kubectl Plugin

Prerequisite: kubectl v1.12.0 or later

With kubectl v1.12.0 introducing [easy pluggability](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) of external functions, kubeaudit can be invoked as `kubectl audit` just by
- running `make plugin` and having $GOPATH/bin available in your path.

or

- renaming the binary to `kubectl-audit` and having it available in your path.



<a name="general" />

## General instructions

`kubeaudit` has three different modes for its audits:
1. Cluster mode
  If kubeaudit detects that it's running in a container, `kubeaudit cmd` will
  attempt to audit the cluster it's running in.
1. Local config mode
  If kubeaudit is not running in a container, `kubeaudit cmd` will audit the
  resources specified by your local kubeconfig (`$HOME/.kube/config`) file.
  You can force kubeaudit to use a specific local config file with the switch
  `-c/--kubeconfig /config/path`
1. Manifest mode
  If you wish to audit a manifest file, use the command
  `kubeaudit -f/--manifest /path/to/manifest.yml`

`kubeaudit` supports two different output types:
1. just running `kubeaudit` will log human readable output
1. if run with `-j/--json` it will log output json formatted so that its output
   can be used by other programs easily

`kubeaudit` supports using manual audit configuration provided by the user, use the command
`kubeaudit -f/--manifest /path/to/manifest.yml -k/--auditConfig /path/to/config.yml`
For more details on audit config check out [Audit Configuration](#audit-configuration).

`kubeaudit` has four different log levels `INFO, WARN, ERROR` controlled by
`-v/--verbose LEVEL` and for those who counted and want to work on `kubeaudit`
`DEBUG`
1. by default the debug level is set to `ERROR` and will log `INFO`, `WARN` and
   `ERROR`
1. if you only care about `ERROR` set it to `ERROR`
1. if you care about `ERROR` and `WARN` set it to `WARN`

But wait! Which version am I actually running? `kubeaudit version` will tell you.

I need help! Run `kubeaudit help` every audit has its own help so you can run
`kubeaudit help sc`

Last but not least before we look at the audits: `kubeaudit -a/--allPods`
audits against pods in all the phases (default Running Phase)

<a name="autofix" />

## Autofix

As humans we are lazy and `kubeaudit` knows that so it comes with the functionality to `autofix` workload manifests. Point it at your workload manifests and it will automagically fix everything so that manifests are as secure as it gets.

`kubeaudit autofix -f path/to/manifest.yml`

The manifest might end up a little too secure for the work it is supposed to do. If that is the case check out [labels](#labels) to opt out of certain checks.

<a name="audits" />

## Audits

`kubeaudit` has multiple checks:
- [Audit all](#all)
- [Audit security context](#sc)
  - [Audit readOnlyRootFilesystem](#rootfs)
  - [Audit runAsNonRoot](#root)
  - [Audit allowPrivilegeEscalation](#allowpe)
  - [Audit privileged](#priv)
  - [Audit capabilities](#caps)
- [Audit image](#image)
- [Audit Service Accounts](#sat)
- [Audit network policies](#netpol)
- [Audit resources](#resources)
- [Audit mounting Docker Socket](#dockersock)
- [Audit AppArmor](#apparmor)
- [Audit Seccomp](#seccomp)
- [Audit namespaces](#namespaces)

<a name="all" />

### Audit all
Runs all the above checks.

```sh
kubeaudit all
ERRO[0000] RunAsNonRoot is not set, which results in root user being allowed!
ERRO[0000] Default serviceAccount with token mounted. Please set automountServiceAccountToken to false
WARN[0000] Privileged defaults to false, which results in non privileged, which is okay.
ERRO[0000] Capability not dropped     CapName=AUDIT_WRITE
```

<a name="sc" />

#### Audit security contexts

The security context holds a couple of different security related
configurations. For convenience, `kubeaudit` will always log the following
information when it creates a log:
```sh
kubeaudit command
LOG[0000] KubeType=deployment Name=THEdeployment Namespace=deploymentNS
```
And for brevity, the information will not be shown in the commands below.

Currently, `kubeaudit` is able to check for the following fields in the security context:

<a name="rootfs" />

#### Audit readOnlyRootFilesystem

`kubeaudit` will detect whether `readOnlyRootFilesystem` is either not set `nil` or explicitly set to `false`

```sh
kubeaudit rootfs
ERRO[0000] ReadOnlyRootFilesystem not set which results in a writable rootFS, please set to true
ERRO[0000] ReadOnlyRootFilesystem set to false, please set to true
```

<a name="root" />

#### Audit runAsNonRoot

`kubeaudit` will detect whether the container is to be run as root:

```sh
kubeaudit nonroot
ERRO[0000] RunAsNonRoot is set to false (root user allowed), please set to true!
ERRO[0000] RunAsNonRoot is not set, which results in root user being allowed!
```

<a name="allowpe" />

#### Audit allowPrivilegeEscalation

`kubeaudit` will detect whether `allowPrivilegeEscalation` is either set to `nil` or explicitly set to `false`

```sh
kubeaudit allowpe
ERRO[0000] AllowPrivilegeEscalation set to true, please set to false
ERRO[0000] AllowPrivilegeEscalation not set which allows privilege escalation, please set to false
```

<a name="priv" />

#### Audit privileged

`kubeaudit` will detect whether the container is to be run privileged:

```sh
kubeaudit priv
ERRO[0000] Privileged set to true! Please change it to false!
```

Since we want to make sure everything is intentionally configured correctly `kubeaudit` warns about `privileged` not being set:

```sh
kubeaudit priv
WARN[0000] Privileged defaults to false, which results in non privileged, which is okay.
```

<a name="caps" />

#### Audit capabilities

Docker comes with a couple of capabilities that shouldn't be needed and
therefore should be dropped. `kubeaudit` will also complain about added capabilities.

If the capabilities field doesn't exist within the security context:

```sh
kubeaudit caps
ERRO[0000] Capabilities field not defined!
```

When capabilities were added:

```sh
kubeaudit caps
ERRO[0000] Capability added  CapName=NET_ADMIN
```

[`config/caps`](https://github.com/Shopify/kubeaudit/blob/master/config/capabilities-drop-list.yml)
holds a list of capabilities that we recommend be dropped, change it if you
want to keep some of the capabilities otherwise `kubeaudit` will complain about
them not being dropped:

```sh
kubeaudit caps
ERRO[0000] Capability not dropped  CapName=AUDIT_WRITE
```

<a name="image" />

### Audit container image tags

`kubeaudit` can check for image names and image tags:

1. If the image tag is incorrect an ERROR will issued
   ```sh
   kubeaudit image -i gcr.io/google_containers/echoserver:1.7
   ERRO[0000] Image tag was incorrect
   ```

1. If the image doesn't have a tag but an image of the name was found a WARNING
   will be created:
   ```sh
   kubeaudit image -i gcr.io/google_containers/echoserver:1.7
   WARN[0000] Image tag was missing
   ```

1. If the image was found with correct tag `kubeaudit` notifies with an INFO message:
   ```sh
   kubeaudit image -i gcr.io/google_containers/echoserver:1.7
   INFO[0000] Image tag was correct
   ```

<a name="sat" />

### Audit Service Accounts

It audits against the following scenarios:

1. A default serviceAccount mounted with a token:
   ```sh
   kubeaudit sat
   ERRO[0000] Default serviceAccount with token mounted. Please set AutomountServiceAccountToken to false
   ```

1. A deprecated service account:
   ```sh
   kubeaudit sat
   WARN[0000] serviceAccount is a deprecated alias for ServiceAccountName, use that one instead  DSA=DeprecatedServiceAccount
   ```

<a name="netpol" />

### Audit network policies

It checks that every namespace should have a default deny network policy
installed. See [Kubernetes Network Policies](https://Kubernetes.io/docs/concepts/services-networking/network-policies/)
for more information:

```sh
kubeaudit np
WARN[0000] Default allow mode on test/testing
```

<a name="resources" />

### Audit resources limits

It checks that every resource has a CPU and memory limit. See [Kubernetes Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/)
for more information:

```sh
kubeaudit limits
WARN[0000] CPU limit not set, please set it!
WARN[0000] Memory limit not set, please set it!
```

With the `--cpu` and `--memory` parameters, `kubeaudit` can check the limits not to be exceeded.

```sh
kubeaudit limits --cpu 500m --memory 125Mi
WARN[0000] CPU limit exceeded, it is set to 1 but it must not exceed 500m. Please adjust it! !
WARN[0000] Memory limit exceeded, it is set to 512Mi but it must not exceed 125Mi. Please adjust it!
```

<a name="dockersock" />

## Audit Mounting Docker Socket

It checks that no container in the pod mounts `/var/run/docker.sock`, as this can be a [very dangerous practice](https://dev.to/petermbenjamin/docker-security-best-practices-45ih). 
If a container does this, it will be indicated as such:

```
containers:
      - image: <image name>
        name: <container name>
        volumeMounts:
        - mountPath: /var/run/docker.sock
          name: <volume name>
volumes:
      - name: <volume name>
        hostPath:
          path: /var/run/docker.sock
```

If `/var/run/docker.sock` is being mounted by a container:

```sh
kubeaudit mountds
WARN[0000] /var/run/docker.sock is being mounted, please avoid this practice. Container=myContainer KubeType=pod Name=myPod Namespace=myNamespace
```

<a name="apparmor" />

## Audit AppArmor

It checks that AppArmor is enabled for all containers by making sure the following annotation exists on the pod.
There must be an annotation for each container in the pod:

```
container.apparmor.security.beta.kubernetes.io/<container name>: <profile>
```

where profile can be "runtime/default" or start with "localhost/" to be considered valid.

If the AppArmor annotation is missing:

```sh
kubeaudit apparmor
ERRO[0000] AppArmor annotation missing. Container=myContainer KubeType=pod Name=myPod Namespace=myNamespace
```

When AppArmor annotations are misconfigured:

```sh
kubeaudit apparmor
ERRO[0000] AppArmor disabled. Annotation=container.apparmor.security.beta.kubernetes.io/myContainer
  Container=myContainer KubeType=pod Name=myPod Namespace=myNamespace Reason=badval
```

<a name="seccomp" />

## Audit Seccomp

It checks that Seccomp is enabled for all containers by making sure one or both of the following annotations exists
on the pod. If no pod annotation is used, then there must be an annotation for each container. Container annotations
override the pod annotation:

```
# pod annotation
seccomp.security.alpha.kubernetes.io/pod: <profile>

# container annotation
container.seccomp.security.alpha.kubernetes.io/<container name>: <profile>
```

where profile can be "runtime/default" or start with "localhost/" to be considered valid. "docker/default" is
deprecated and will show a warning. It should be replaced with "runtime/default".

If the Seccomp annotation is missing:

```sh
kubeaudit seccomp
ERRO[0000] Seccomp annotation missing. Container=myContainer KubeType=pod Name=myPod Namespace=myNamespace
```

When Seccomp annotations are misconfigured for a container:

```sh
kubeaudit seccomp
ERRO[0000] Seccomp disabled for container. Annotation=container.seccomp.security.alpha.kubernetes.io/myContainer
  Container=myContainer KubeType=pod Name=myPod Namespace=myNamespace Reason=badval
```

When Seccomp annotations are misconfigured for a pod:
```sh
kubeaudit seccomp
ERRO[0000] Seccomp disabled for pod. Annotation=seccomp.security.alpha.kubernetes.io/pod Container= KubeType=pod
  Name=myPod Namespace=myNamespace Reason=unconfined
```

<a name="namespaces" />

## Audit namespaces

`kubeaudit` will detect whether `hostNetwork`,`hostIPC` or `hostPID` is either set to `true` in `podSpec` for `Pod` workloads 

```sh
kubeaudit namespaces
ERRO[0000] hostNetwork is set to true  in podSpec, please set to false!
ERRO[0000] hostIPC is set to true  in podSpec, please set to false!
ERRO[0000] hostPID is set to true  in podSpec, please set to false!
```

<a name="labels" />

## Override Labels

Override labels give you the ability to have `kubeaudit` allow certain audits to fail.
For example, if you want `kubeaudit` to ignore the fact that `AllowPrivilegeEscalation` was set to `true`, you can add the following label:

```sh
spec:
  template:
    metadata:
      labels:
        apps: YourAppNameHere
        container.audit.kubernetes.io/<container-name>/allow-privilege-escalation: "YourReasonForOverrideHere"
```

Any label with a non-nil reason string will prevent `kubeaudit` from throwing the corresponding error and issue a warning instead.
Reasons matching `"true"` (not case sensitive) will be displayed as `Unspecified`.

`kubeaudit` can skip certain audits by applying override labels to containers. If you want skip an audit for a specific container inside a pod, you can add an container override label. For example, if you use `kubeaudit` to ignore `allow-run-as-root` check for container "MyContainer1", you can add the following label:

```sh
spec:
  template:
    metadata:
      labels:
        apps: YourAppNameHere
        container.audit.kubernetes.io/MyContainer1/allow-run-as-root: "YourReasonForOverrideHere"
```

Similarly, you can have `kubeaudit` to skip a specific audit for all containers inside the pod by adding a pod override label. For example, if you use `kubeaudit` to ignore `allow-run-as-root` check for all containers inside the pod, you can add the following label:

```sh
spec:
  template:
    metadata:
      labels:
        apps: YourAppNameHere
        audit.kubernetes.io/pod/allow-run-as-root: "YourReasonForOverrideHere"
```

`kubeaudit` can also skip a specific audit for all network policies associated with a namespace resource
by adding a namespace override label. For example, if you use `kubeaudit` to ignore the `allow-non-default-deny-egress-network-policy` check for the namespace `namespaceName1` you can add the following label to the namespace:

```sh
metadata:
  name: namespaceName1
  labels:
    audit.kubernetes.io/namespaceName1/allow-non-default-deny-egress-network-policy: "YourReasonForOverrideHere"
```

`kubeaudit` supports many labels on pod, namespace or container level:
- [audit.kubernetes.io/pod/allow-privilege-escalation](#allowpe_label)
- [container.audit.kubernetes.io/\<container-name\>/allow-privilege-escalation](#allowpe_label)
- [audit.kubernetes.io/pod/allow-privileged](#priv_label)
- [container.audit.kubernetes.io/\<container-name\>/allow-privileged](#priv_label)
- [audit.kubernetes.io/pod/allow-capability](#caps_label)
- [container.audit.kubernetes.io/\<container-name\>/allow-capability](#caps_label)
- [audit.kubernetes.io/pod/allow-run-as-root](#nonroot_label)
- [container.audit.kubernetes.io/\<container-name\>/allow-run-as-root](#nonroot_label)
- [audit.kubernetes.io/pod/allow-automount-service-account-token](#sat_label)
- [audit.kubernetes.io/pod/allow-read-only-root-filesystem-false](#rootfs_label)
- [container.audit.kubernetes.io/\<container-name\>/allow-read-only-root-filesystem-false](#rootfs_label)
- [audit.kubernetes.io/\<namespace-name\>/allow-non-default-deny-egress-network-policy](#egress_label)
- [audit.kubernetes.io/\<namespace-name\>/allow-non-default-deny-ingress-network-policy](#ingress_label)
- [audit.kubernetes.io/pod/allow-namespace-host-network](#namespacenetwork_label)
- [audit.kubernetes.io/pod/allow-namespace-host-IPC](#namespaceipc_label)
- [audit.kubernetes.io/pod/allow-namespace-host-PID](#namespacepid_label)

<a name="allowpe_label"/>

### container.audit.kubernetes.io/\<container-name\>/allow-privilege-escalation

Allow `allowPrivilegeEscalation` to be set to `true` to a specific container.

### audit.kubernetes.io/pod/allow-privilege-escalation

Allows `allowPrivilegeEscalation` to be set to `true` to all the containers in a pod.

```sh
kubeaudit.allow.privilegeEscalation: "Superuser privileges needed"

WARN[0000] Allowed setting AllowPrivilegeEscalation to true  Reason="Superuser privileges needed"
```

<a name="priv_label"/>

### container.audit.kubernetes.io/\<container-name\>/allow-privileged

Allow `privileged` to be set to `true` to a specific container.

### audit.kubernetes.io/pod/allow-privileged

Allows `privileged` to be set to `true` to all the containers in a pod.

```sh
kubeaudit.allow.privileged: "Privileged execution required"

WARN[0000] Allowed setting privileged to true                Reason="Privileged execution required"
```

<a name="caps_label"/>

### container.audit.kubernetes.io/\<container-name\>/allow-capability

Allows adding a capability or keeping one that would otherwise be dropped to a specific container.

### audit.kubernetes.io/pod/allow-capability

Allows adding a capability or keeping one that would otherwise be dropped to all the containers in a pod.

```sh
kubeaudit.allow.capability.chown: "true"

WARN[0000] Capability allowed                                CapName=CHOWN Reason=Unspecified
```

<a name="nonroot_label"/>

### container.audit.kubernetes.io/\<container-name\>/allow-run-as-root

Allows setting `runAsNonRoot` to `false` to a specific container.

### audit.kubernetes.io/pod/allow-run-as-root

Allows setting `runAsNonRoot` to `false` to all the containers in a pod.

```sh
kubeaudit.allow.runAsRoot: "Root privileges needed"

WARN[0000] Allowed setting RunAsNonRoot to false             Reason="Root privileges needed"
```

<a name="sat_label"/>

### audit.kubernetes.io/pod/allow-automount-service-account-token

Allows setting `automountServiceAccountToken` to `true` to a pod.

```sh
kubeaudit.allow.autmountServiceAccountToken: "True"

WARN[0000] Allowed setting automountServiceAccountToken to true  Reason=Unspecified
```

<a name="rootfs_label"/>

### container.audit.kubernetes.io/\<container-name\>/allow-read-only-root-filesystem-false

Allows setting `readOnlyRootFilesystem` to `false` to a specific container.

### audit.kubernetes.io/pod/allow-read-only-root-filesystem-false

Allows setting `readOnlyRootFilesystem` to `false` to all containers in a pod.

```sh
kubeaudit.allow.readOnlyRootFilesystemFalse: "Write permissions needed"

WARN[0000] Allowed setting readOnlyRootFilesystem to false Reason="Write permissions needed"
```

<a name="egress_label"/>

### audit.kubernetes.io/\<namespace-name\>/allow-non-default-deny-egress-network-policy

Allows absense of `default-deny` egress network policy for that specific namespace.

<a name="ingress_label"/>

### audit.kubernetes.io/\<namespace-name\>/allow-non-default-deny-ingress-network-policy

Allows absense of `default-deny` ingress network policy for that specific namespace.

```sh
audit.kubernetes.io/default/allow-non-default-deny-egress-network-policy: "Egress is allowed"

WARN[0000] Allowed Namespace without a default deny egress NetworkPolicy  KubeType=namespace Name=default Reason="Egress is allowed"
```

<a name="namespacenetwork_label"/>

### audit.kubernetes.io/pod/allow-namespace-host-network

```sh
audit.kubernetes.io/pod/allow-namespace-host-network: "hostNetwork is allowed"

WARN[0000] Allowed setting hostNetwork to true           KubeType=pod Name=Pod Namespace=PodNamespace Reason="hostNetwork is allowed"
```

<a name="namespaceipc_label"/>

### audit.kubernetes.io/pod/allow-namespace-host-IPC

```sh
audit.kubernetes.io/pod/allow-namespace-host-IPC: "hostIPC is allowed"

WARN[0000] Allowed setting hostIPC to true               KubeType=pod Name=Pod Namespace=PodNamespace Reason="hostIPC is allowed"
```

<a name="namespacepid_label"/>

### audit.kubernetes.io/pod/allow-namespace-host-PID

```sh
audit.kubernetes.io/pod/allow-namespace-host-PID: "hostPID is allowed"

WARN[0000] Allowed setting hostPID to true               KubeType=pod Name=Pod Namespace=PodNamespace Reason="hostPID is allowed"
```

<a name="contribute" />

## Drop capabilities list

Allows configuring the audit against drop capabilities. Sane defaults are as follows:

```
# SANE DEFAULTS:
capabilitiesToBeDropped:
  # https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities
  - SETPCAP #Modify process capabilities.
  - MKNOD #Create special files using mknod(2).
  - AUDIT_WRITE #Write records to kernel auditing log.
  - CHOWN #Make arbitrary changes to file UIDs and GIDs (see chown(2)).
  - NET_RAW #Use RAW and PACKET sockets.
  - DAC_OVERRIDE #Bypass file read, write, and execute permission checks.
  - FOWNER #Bypass permission checks on operations that normally require the file system UID of the process to match the UID of the file.
  - FSETID #Don’t clear set-user-ID and set-group-ID permission bits when a file is modified.
  - KILL #Bypass permission checks for sending signals.
  - SETGID #Make arbitrary manipulations of process GIDs and supplementary GID list.
  - SETUID #Make arbitrary manipulations of process UIDs.
  - NET_BIND_SERVICE #Bind a socket to internet domain privileged ports (port numbers less than 1024).
  - SYS_CHROOT #Use chroot(2), change root directory.
  - SETFCAP #Set file capabilities.
```

This can be overridden by using `-k` flag and providing your own defaults in the yaml format as shown below.

<a name="audit-configuration" />

## Audit Configuration

Allows configuring your own audit settings for kubeaudit. By default following configuration is used:

```
apiVersion: v1
kind: kubeauditConfig
audit: true  # Set to false if you want kubeaudit to not audit your k8s manifests
spec:
  capabilities: # List of all supported capabilities
    NET_ADMIN: drop         # Set to `keep` to keep capability
    SETPCAP: drop           # Set to `keep` to keep capability
    MKNOD: drop             # Set to `keep` to keep capability
    AUDIT_WRITE: drop       # Set to `keep` to keep capability
    CHOWN: drop             # Set to `keep` to keep capability
    NET_RAW: drop           # Set to `keep` to keep capability
    DAC_OVERRIDE: drop      # Set to `keep` to keep capability
    FOWNER: drop            # Set to `keep` to keep capability
    FSETID: drop            # Set to `keep` to keep capability
    KILL: drop              # Set to `keep` to keep capability
    SETGID: drop            # Set to `keep` to keep capability
    SETUID: drop            # Set to `keep` to keep capability
    NET_BIND_SERVICE: drop  # Set to `keep` to keep capability
    SYS_CHROOT: drop        # Set to `keep` to keep capability
    SETFCAP: drop           # Set to `keep` to keep capability
  overrides: # List of all supported overrides
    privilege-escalation: deny                      # Set to `allow` to skip auditing potential vulnerability
    privileged: deny                                # Set to `allow` to skip auditing potential vulnerability
    run-as-root: deny                               # Set to `allow` to skip auditing potential vulnerability
    automount-service-account-token: deny           # Set to `allow` to skip auditing potential vulnerability
    read-only-root-filesystem-false: deny           # Set to `allow` to skip auditing potential vulnerability
    non-default-deny-ingress-network-policy: deny   # Set to `allow` to skip auditing potential vulnerability
    non-default-deny-egress-network-policy: deny    # Set to `allow` to skip auditing potential vulnerability
    namespace-host-network: deny                    # Set to `allow` to skip auditing potential vulnerability
    namespace-host-IPC: deny                        # Set to `allow` to skip auditing potential vulnerability
    namespace-host-PID: deny                        # Set to `allow` to skip auditing potential vulnerability
```

<a name="contribute" />

## Contributing

If you'd like to fix a bug, contribute a feature or just correct a typo, please feel free to do so as long as you follow our [Code of Conduct](https://github.com/Shopify/kubeaudit/blob/master/CODE_OF_CONDUCT.md).

1. Create your own fork!
1. Get the source: `go get github.com/Shopify/kubeaudit`
1. Go to the source: `cd $GOPATH/src/github.com/Shopify/kubeaudit`
1. Add your forked repo as a fork: `git remote add fork https://github.com/you-are-awesome/kubeaudit`
1. Create your feature branch: `git checkout -b awesome-new-feature`
1. Run the tests to see everything is working as expected: `make test`
1. Commit your changes: `git commit -am 'Adds awesome feature'`
1. Push to the branch: `git push fork`
1. Submit a PR (All PR must be labeled with :bug: (Bug fix), :sparkles: (New feature), :book: (Documentation update), or :warning: (Breaking changes) )
1. ???
1. Profit
