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
- [Audits](#audits)
- [Override Labels](#labels)
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

<a name="general" />

## General instructions

`kubeaudit` has three different modes for its audits:
1. `kubeaudit cmd` will attempt to create an in-cluster client and audit.
1. `kubeaudit -l/--local cmd` will use your kubeconfig (`~/.kube/config` or if
   you need different path use `-c /config/path`
1. `kubeaudit -f/--manifest /path/to/manifest.yml` will audit the manifest

`kubeaudit` supports to different output types:
1. just running `kubeaudit` will log human readable output
1. if run with `-j/--json` it will log output json formatted so that its output
   can be used by other programs easily

`kubeaudit` has 4 different log levels `INFO, WARN, ERROR` controlled by
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

<a name="audits" />

## Audits

`kubeaudit` has multiple checks:
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

<a name="sc" />

#### Audit security contexts

The security context holds a couple of different security related
configurations. For convenience, `kubeaudit` will always log the following
information when it creates a log:
```sh
kubeaudit -l command
LOG[0000] KubeType=deployment Name=THEdeployment Namespace=deploymentNS
```
And for brevity, the information will not be shown in the commands below.

Currently, `kubeaudit` is able to check for the following fields in the security context:

<a name="rootfs" />

#### Audit readOnlyRootFilesystem

`kubeaudit` will detect whether `readOnlyRootFilesystem` is either not set `nil` or explicitly set to `false`

```sh
kubeaudit -l rootfs
ERRO[0000] ReadOnlyRootFilesystem not set which results in a writable rootFS, please set to true
ERRO[0000] ReadOnlyRootFilesystem set to false, please set to true
```

<a name="root" />

#### Audit runAsNonRoot

`kubeaudit` will detect whether the container is to be run as root:

```sh
kubeaudit -l nonroot
ERRO[0000] RunAsNonRoot is set to false (root user allowed), please set to true!
ERRO[0000] RunAsNonRoot is not set, which results in root user being allowed!
```

<a name="allowpe" />

#### Audit allowPrivilegeEscalation

`kubeaudit` will detect whether `allowPrivilegeEscalation` is either set to `nil` or explicitly set to `false`

```sh
kubeaudit -l allowpe
ERRO[0000] AllowPrivilegeEscalation set to true, please set to false
ERRO[0000] AllowPrivilegeEscalation not set which allows privilege escalation, please set to false
```

<a name="priv" />

#### Audit privileged

`kubeaudit` will detect whether the container is to be run privileged:

```sh
kubeaudit -l priv
ERRO[0000] Privileged set to true! Please change it to false!
```

Since we want to make sure everything is intentionally configured correctly `kubeaudit` warns about `privileged` not being set:

```sh
kubeaudit -l priv
WARN[0000] Privileged defaults to false, which results in non privileged, which is okay.
```

<a name="caps" />

#### Audit capabilities

Docker comes with a couple of capabilities that shouldn't be needed and
therefore should be dropped. `kubeaudit` will also complain about added capabilities.

If the capabilities field doesn't exist within the security context:

```sh
kubeaudiit -l caps
ERRO[0000] Capabilities field not defined!
```

When capabilities were added:

```sh
kubeaudiit -l caps
ERRO[0000] Capability added  CapName=NET_ADMIN
```

[`config/caps`](https://github.com/Shopify/kubeaudit/blob/master/config/capabilities-drop-list.yml)
holds a list of capabilities that we recommend be dropped, change it if you
want to keep some of the capabilities otherwise `kubeaudit` will complain about
them not being dropped:

```sh
kubeaudiit -l caps
ERRO[0000] Capability not dropped  CapName=AUDIT_WRITE
```

<a name="image" />

### Audit container image tags

`kubeaudit` can check for image names and image tags:

1. If the image tag is incorrect an ERROR will issued
```sh
kubeaudit -l image -i gcr.io/google_containers/echoserver:1.7
ERRO[0000] Image tag was incorrect
```

1. If the image doesn't have a tag but an image of the name was found a WARNING
   will be created:
```sh
kubeaudit -l image -i gcr.io/google_containers/echoserver:1.7
WARN[0000] Image tag was missing
```

1. If the image was found with correct tag `kubeaudit` notifies with an INFO message:
```sh
kubeaudit -l image -i gcr.io/google_containers/echoserver:1.7
INFO[0000] Image tag was correct
```

<a name="sat" />

### Audit Service Accounts

It audits against the following scenarios:

1.  A default serviceAccount mounted with a token:
```sh
kubeaudit -l sat
ERRO[0000] Default serviceAccount with token mounted. Please set AutomountServiceAccountToken to false
```

1.  A deprecated service account:
```sh
kubeaudit -l sat
WARN[0000] serviceAccount is a deprecated alias for ServiceAccountName, use that one instead  DSA=DeprecatedServiceAccount
```

<a name="netpol" />

### Audit network policies

It checks that every namespace should have a default deny network policy
installed. See [Kubernetes Network Policies](https://Kubernetes.io/docs/concepts/services-networking/network-policies/)
for more information:

```sh
# don't specify -l or -c to run inside the cluster
kubeaudit np
WARN[0000] Default allow mode on test/testing
```

<a name="resources" />

### Audit resources limits

It checks that every resource has a CPU and memory limit. See [Kubernetes Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/)
for more information:

```sh
kubeaudit -l limits
WARN[0000] CPU limit not set, please set it!
WARN[0000] Memory limit not set, please set it!
```

With the `--cpu` and `--memory` parameters, `kubeaudit` can check the limits not to be exceeded.

```sh
kubeaudit -l limits --cpu 500m --memory 125Mi
WARN[0000] CPU limit exceeded, it is set to 1 but it must not exceed 500m. Please adjust it! !
WARN[0000] Memory limit exceeded, it is set to 512Mi but it must not exceed 125Mi. Please adjust it!
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
        kubeaudit.allow.privilegeEscalation: "YourReasonForOverrideHere"
```

Any label with a non-nil reason string will prevent `kubeaudit` from throwing the corresponding error and issue a warning instead.
Reasons matching `"true"` (not case sensitive) will be displayed as `Unspecified`.

`kubeaudit` supports many labels:
- [kubeaudit.allow.privilegeEscalation](#allowpe_label)
- [kubeaudit.allow.privileged](#priv_label)
- [kubeaudit.allow.capability](#caps_label)
- [kubeaudit.allow.runAsRoot](#nonroot_label)
- [kubeaudit.allow.automountServiceAccountToken](#sat_label)
- [kubeaudit.allow.readOnlyRootFilesystemFalse](#rootfs_label)

<a name="allowpe_label"/>

### kubeaudit.allow.privilegeEscalation

Allows `allowPrivilegeEscalation` to be set to `true`.

```sh
kubeaudit.allow.privilegeEscalation: "Superuser privileges needed"

WARN[0000] Allowed setting AllowPrivilegeEscalation to true  Reason="Superuser privileges needed"
```

<a name="priv_label"/>

### kubeaudit.allow.privileged

Allows `privileged` to be set to `true`.

```sh
kubeaudit.allow.privileged: "Privileged execution required"

WARN[0000] Allowed setting privileged to true                Reason="Privileged execution required"
```

<a name="caps_label"/>

### kubeaudit.allow.capability

Allows adding a capability or keeping one that would otherwise be dropped.

```sh
kubeaudit.allow.capability.chown: "true"

WARN[0000] Capability allowed                                CapName=CHOWN Reason=Unspecified
```

<a name="nonroot_label"/>

### kubeaudit.allow.runAsRoot

Allows setting `runAsNonRoot` to `false`.

```sh
kubeaudit.allow.runAsRoot: "Root privileges needed"

WARN[0000] Allowed setting RunAsNonRoot to false             Reason="Root privileges needed"
```

<a name="sat_label"/>

### kubeaudit.automountServiceAccountToken

Allows setting `automountServiceAccountToken` to `true`.

```sh
kubeaudit.allow.autmountServiceAccountToken: "True"

WARN[0000] Allowed setting automountServiceAccountToken to true  Reason=Unspecified
```

<a name="rootfs_label"/>

### kubeaudit.readOnlyRootFilesystemFalse

Allows setting `readOnlyRootFilesystem` to `false`

```sh
kubeaudit.allow.readOnlyRootFilesystemFalse: "Write permissions needed"

WARN[0000] Allowed setting readOnlyRootFilesystem to false Reason="Write permissions needed"
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
1. Submit a PR
1. ???
1. Profit
