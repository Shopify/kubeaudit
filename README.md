[![Build Status](https://api.travis-ci.org/Shopify/kubeaudit.svg?branch=master)](https://travis-ci.org/Shopify/kubeaudit/)
![cover.run go](https://cover.run/go/github.com/Shopify/kubeaudit/cmd.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Shopify/kubeaudit)](https://goreportcard.com/report/github.com/Shopify/kubeaudit)
[![GoDoc](https://godoc.org/github.com/Shopify/kubeaudit?status.png)](https://godoc.org/github.com/Shopify/kubeaudit)
# kubeaudit

:cloud: :lock: :muscle:

## What's all this then?

`kubeaudit` is a command line tool written in golang to help you audit your Kubernetes cluster. This tool can audit for the following scenarios:
- [Audit security context](#sc)
- [Audit container image](#image)
- [Audit network policies](#netpol)
- [Audit RBAC policies](#rbac)
- [Audit .yamls](#yamls)

## Installation

#### Download a binary

[Official releases](https://github.com/Shopify/kubeaudit/releases)

#### DIY build

Add kubeaudit and its dependencies by running the following command:

```sh
go get -v github.com/Shopify/kubeaudit
make
```

Upon completion you should find kubeaudit in `$GOPATH/bin/kubeaudit`

## Running tests
```sh
make test
```

## Usage

### General instructions

kubeaudit is driven by [cobra](https://github.com/spf13/cobra) on the command line
```raw
kubeaudit is a program that will help you audit
your Kubernetes clusters. Specify -l to run kubeaudit using ~/.kube/config
otherwise it will attempt to create an in-cluster client.

#patcheswelcome

Usage:
  kubeaudit [command]

Available Commands:
  help        Help about any command
  image       Audit container images
  np          Audit namespace network policies
  rbac        Audit RBAC things
  sc          Audit container security contexts
  version     Print the version number of kubeaudit

Flags:
  -a, --allPods             Audit againsts pods in all the phases (default Running Phase)
  -h, --help                help for kubeaudit
  -j, --json                Enable json logging
  -c, --kubeconfig string   config file (default is $HOME/.kube/config
  -f, --manifest string     yaml configuration to audit
  -l, --local               Local mode, uses ~/.kube/config as configuration
  -v, --verbose             Enable debug (verbose) logging

Use "kubeaudit [command] --help" for more information about a command.
```

<a name="sc" />

### Audit security contexts

It can audit against three different scenarios.

1. General security context which make sure that every Kubernetes pod has a proper security context i.e. privileged linux capabilities are dropped or not:

```sh
% kubeaudit -l sc
ERRO[0004] test/testDeployment                                                       type=deployment
WARN[0004] Capabilities added to test/testStateSet  caps="[IPC_LOCK SYS_RESOURCE]"   type=statefulSet
WARN[0004] No capabilities were dropped! test/testDaemonSet                          type=daemonSet
```

2. Every Kubernetes pod should have a read-only root file system:

```sh
% kubeaudit -l sc rootfs
ERRO[0005] testbuilder/testpod-312-3213                  type=pod
```

3. Every container is running as non-root user:

```sh
% kubeaudit -l sc nonroot
ERRO[0004] test/testPod                                  type=pod
```

<a name="image" />

### Audit container image tags

It checks that every Kubernetes resource is running the specified tag of a given image:

```sh
% kubeaudit -l image -i gcr.io/google_containers/echoserver:1.7
ERRO[0005] test/testReplicationController               type=replicationController
```

<a name="netpol" />

### Audit network policies

It checks that every namespace should have a default deny network policiy installed. See [Kubernetes Network Policies](https://Kubernetes.io/docs/concepts/services-networking/network-policies/) for more information:

```sh
# don't specify -l or -c to run inside the clsuter
% kubeaudit np
WARN[0000] Default allow mode on test/testing           type=netpol
```

<a name="rbac" />

### Audit RBAC policies

It audits against the following scenarios:

- Check for automountServiceAccountToken is nil with no serviceAccountName
- Check for usage of deprecated serviceAccount

```sh
% kubeaudit -l rbac sat
ERRO[0000] automountServiceAccountToken nil (mounted by default) with no serviceAccountName name=alpine namespace=test type=deployment
WARN[0000] deprecated serviceAccount detected (sub for serviceAccountName)  name=nginx namespace=staging serviceAccount=nginx serviceAccountName=nginx type=deployment
```

<a name="yamls" />

### Audit .yamls

Kubeaudit can audit undeployed resources when defined in a yaml as well:

```sh
% kubeaudit -f /path/to/your.yml pick_your_action_from_above
ERRO{0000] ...
```

## Contributing

If you'd like to fix a bug, contribute a feature or just correct a typo, please feel free to do so as long as you follow our [Code of Conduct](https://github.com/Shopify/kubeaudit/blob/master/CODE_OF_CONDUCT.md).

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a PR
