[![Build Status](https://github.com/Shopify/kubeaudit/actions/workflows/ci.yml/badge.svg)](https://github.com/Shopify/kubeaudit/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Shopify/kubeaudit)](https://goreportcard.com/report/github.com/Shopify/kubeaudit)
[![GoDoc](https://godoc.org/github.com/Shopify/kubeaudit?status.png)](https://godoc.org/github.com/Shopify/kubeaudit)

> Kubeaudit no longer supports APIs deprecated as of [Kubernetes v.1.16 release](https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/). So, it is now a requirement for clusters to run Kubernetes >=1.16


# kubeaudit :cloud: :lock: :muscle:

`kubeaudit` is a command line tool and a Go package to audit Kubernetes clusters for various
different security concerns, such as:
* run as non-root
* use a read-only root filesystem
* drop scary capabilities, don't add new ones
* don't run privileged
* and more!

**tldr. `kubeaudit` makes sure you deploy secure containers!**

## Package
To use kubeaudit as a Go package, see the [package docs](https://pkg.go.dev/github.com/Shopify/kubeaudit).

The rest of this README will focus on how to use kubeaudit as a command line tool.

## Command Line Interface (CLI)

* [Installation](#installation)
* [Quick Start](#quick-start)
* [Audit Results](#audit-results)
* [Commands](#commands)
* [Configuration File](#configuration-file)
* [Override Errors](#override-errors)
* [Contributing](#contributing)

## Installation

### Brew

```
brew install kubeaudit
```

### Download a binary

Kubeaudit has official releases that are blessed and stable:
[Official releases](https://github.com/Shopify/kubeaudit/releases)

### DIY build

Master may have newer features than the stable releases. If you need a newer
feature not yet included in a release, make sure you're using Go 1.17+ and run
the following:

```sh
go get -v github.com/Shopify/kubeaudit
```

Start using `kubeaudit` with the [Quick Start](#quick-start) or view all the [supported commands](#commands).

### Kubectl Plugin

Prerequisite: kubectl v1.12.0 or later

With kubectl v1.12.0 introducing [easy pluggability](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) of external functions, kubeaudit can be invoked as `kubectl audit` by

- running `make plugin` and having `$GOPATH/bin` available in your path.

or

- renaming the binary to `kubectl-audit` and having it available in your path.

### Docker

We also release a [Docker image](https://hub.docker.com/r/shopify/kubeaudit): `shopify/kubeaudit`. To run kubeaudit as a job in your cluster see [Running kubeaudit in a cluster](docs/cluster.md).

## Quick Start

kubeaudit has three modes:

1. Manifest mode
1. Local mode
1. Cluster mode

### Manifest Mode

If a Kubernetes manifest file is provided using the `-f/--manifest` flag, kubeaudit will audit the manifest file.

Example command:
```
kubeaudit all -f "/path/to/manifest.yml"
```

Example output:
```
$ kubeaudit all -f "internal/test/fixtures/all_resources/deployment-apps-v1.yml"

---------------- Results for ---------------

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: deployment
    namespace: deployment-apps-v1

--------------------------------------------

-- [error] AppArmorAnnotationMissing
   Message: AppArmor annotation missing. The annotation 'container.apparmor.security.beta.kubernetes.io/container' should be added.
   Metadata:
      Container: container
      MissingAnnotation: container.apparmor.security.beta.kubernetes.io/container

-- [error] AutomountServiceAccountTokenTrueAndDefaultSA
   Message: Default service account with token mounted. automountServiceAccountToken should be set to 'false' or a non-default service account should be used.

-- [error] CapabilityShouldDropAll
   Message: Capability not set to ALL. Ideally, you should drop ALL capabilities and add the specific ones you need to the add list.
   Metadata:
      Container: container
      Capability: AUDIT_WRITE
...
```

If no errors with a given minimum severity are found, the following is returned:

```shell
All checks completed. 0 high-risk vulnerabilities found
```

#### Autofix

Manifest mode also supports autofixing all security issues using the `autofix` command:

```
kubeaudit autofix -f "/path/to/manifest.yml"
```

To write the fixed manifest to a new file instead of modifying the source file, use the `-o/--output` flag.

```
kubeaudit autofix -f "/path/to/manifest.yml" -o "/path/to/fixed"
```

To fix a manifest based on custom rules specified on a kubeaudit config file, use the `-k/--kconfig` flag.

```
kubeaudit autofix -k "/path/to/kubeaudit-config.yml" -f "/path/to/manifest.yml" -o "/path/to/fixed"
```

### Cluster Mode

Kubeaudit can detect if it is running within a container in a cluster. If so, it will try to audit all Kubernetes resources in that cluster:
```
kubeaudit all
```

### Local Mode

Kubeaudit will try to connect to a cluster using the local kubeconfig file (`$HOME/.kube/config`). A different kubeconfig location can be specified using the `-c/--kubeconfig` flag.
```
kubeaudit all -c "/path/to/config"
```

For more information on kubernetes config files, see https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/

## Audit Results

Kubeaudit produces results with three levels of severity:

`Error`: A security issue or invalid kubernetes configuration
`Warning`: A best practice recommendation
`Info`: Informational, no action required. This includes results that are [overridden](#override-errors)

The minimum severity level can be set using the `--minSeverity/-m` flag.

By default kubeaudit will output results in a human-readable way. If the output is intended to be further processed, it can be set to output JSON using the `--format json` flag. To output results as logs (the previous default) use `--format logrus`.

If there are results of severity level `error`, kubeaudit will exit with exit code 2. This can be changed using the `--exitcode/-e` flag.

For all the ways kubeaudit can be customized, see [Global Flags](#global-flags).

## Commands

| Command   | Description                                                               | Documentation           |
| :-------- | :------------------------------------------------------------------------ | :---------------------- |
| `all`     | Runs all available auditors, or those specified using a kubeaudit config. | [docs](docs/all.md)     |
| `autofix` | Automatically fixes security issues.                                      | [docs](docs/autofix.md) |
| `version` | Prints the current kubeaudit version.                                     |                         |

### Auditors

Auditors can also be run individually.

| Command        | Description                                                                                                    | Documentation                         |
| :------------- | :------------------------------------------------------------------------------------------------------------- | :------------------------------------ |
| `apparmor`     | Finds containers running without AppArmor.                                                                     | [docs](docs/auditors/apparmor.md)     |
| `asat`         | Finds pods using an automatically mounted default service account                                              | [docs](docs/auditors/asat.md)         |
| `capabilities` | Finds containers that do not drop the recommended capabilities or add new ones.                                | [docs](docs/auditors/capabilities.md) |
| `hostns`       | Finds containers that have HostPID, HostIPC or HostNetwork enabled.                                            | [docs](docs/auditors/hostns.md)       |
| `image`        | Finds containers which do not use the desired version of an image (via the tag) or use an image without a tag. | [docs](docs/auditors/image.md)        |
| `limits`       | Finds containers which exceed the specified CPU and memory limits or do not specify any.                       | [docs](docs/auditors/limits.md)       |
| `mounts`       | Finds containers that have sensitive host paths mounted.                                                       | [docs](docs/auditors/mounts.md)       |
| `netpols`      | Finds namespaces that do not have a default-deny network policy.                                               | [docs](docs/auditors/netpols.md)      |
| `nonroot`      | Finds containers running as root.                                                                              | [docs](docs/auditors/nonroot.md)      |
| `privesc`      | Finds containers that allow privilege escalation.                                                              | [docs](docs/auditors/privesc.md)      |
| `privileged`   | Finds containers running as privileged.                                                                        | [docs](docs/auditors/privileged.md)   |
| `rootfs`       | Finds containers which do not have a read-only filesystem.                                                     | [docs](docs/auditors/rootfs.md)       |
| `seccomp`      | Finds containers running without Seccomp.                                                                      | [docs](docs/auditors/seccomp.md)      |

### Global Flags

| Short | Long               | Description                                                                                                                                            |
| :---- | :----------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------- |
|       | --format           | The output format to use (one of "pretty", "logrus", "json") (default is "pretty")                                                                     |
| -c    | --kubeconfig       | Path to local Kubernetes config file. Only used in local mode (default is `$HOME/.kube/config`)                                                        |
| -f    | --manifest         | Path to the yaml configuration to audit. Only used in manifest mode. You may use `-` to read from stdin.                                               |
| -n    | --namespace        | Only audit resources in the specified namespace. Not currently supported in manifest mode.                                                             |
| -g    | --includegenerated | Include generated resources in scan  (such as Pods generated by deployments). If you would like kubeaudit to produce results for generated resources (for example if you have custom resources or want to catch orphaned resources where the owner resource no longer exists) you can use this flag. |
| -m    | --minseverity      | Set the lowest severity level to report (one of "error", "warning", "info") (default "info")                                                           |
| -e    | --exitcode         | Exit code to use if there are results with severity of "error". Conventionally, 0 is used for success and all non-zero codes for an error. (default 2) |

## Configuration File

The kubeaudit config can be used for two things:

1. Enabling only some auditors
1. Specifying configuration for auditors

Any configuration that can be specified using flags for the individual auditors can be represented using the config.

The config has the following format:

```yaml
enabledAuditors:
  # Auditors are enabled by default if they are not explicitly set to "false"
  apparmor: false
  asat: false
  capabilities: true
  hostns: true
  image: true
  limits: true
  mounts: true
  netpols: true
  nonroot: true
  privesc: true
  privileged: true
  rootfs: true
  seccomp: true
auditors:
  capabilities:
    # add capabilities needed to the add list, so kubeaudit won't report errors
    allowAddList: ['AUDIT_WRITE', 'CHOWN']
  image:
    # If no image is specified and the 'image' auditor is enabled, WARN results
    # will be generated for containers which use an image without a tag
    image: 'myimage:mytag'
  limits:
    # If no limits are specified and the 'limits' auditor is enabled, WARN results
    # will be generated for containers which have no cpu or memory limits specified
    cpu: '750m'
    memory: '500m'
```

For more details about each auditor, including a description of the auditor-specific configuration in the config, see the [Auditor Docs](#auditors).

**Note**: The kubeaudit config is not the same as the kubeconfig file specified with the `-c/--kubeconfig` flag, which refers to the Kubernetes config file (see [Local Mode](/README.md#local-mode)). Also note that only the `all` and `autofix` commands support using a kubeaudit config. It will not work with other commands.

**Note**: If flags are used in combination with the config file, flags will take precedence.

## Override Errors

Security issues can be ignored for specific containers or pods by adding override labels. This means the auditor will produce `info` results instead of `error` results and the audit result name will have `Allowed` appended to it. The labels are documented in each auditor's documentation, but the general format for auditors that support overrides is as follows:

An override label consists of a `key` and a `value`.

The `key` is a combination of the override type (container or pod) and an `override identifier` which is unique to each auditor (see the [docs](#auditors) for the specific auditor). The `key` can take one of two forms depending on the override type:

1. **Container overrides**, which override the auditor for that specific container, are formatted as follows:

```yaml
container.audit.kubernetes.io/[container name].[override identifier]
```

2. **Pod overrides**, which override the auditor for all containers within the pod, are formatted as follows:

```yaml
audit.kubernetes.io/pod.[override identifier]
```

If the `value` is set to a non-empty string, it will be displayed in the `info` result as the `OverrideReason`:

```
$ kubeaudit asat -f "auditors/asat/fixtures/service-account-token-true-allowed.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: ReplicationController
  metadata:
    name: replicationcontroller
    namespace: service-account-token-true-allowed

--------------------------------------------

-- [info] AutomountServiceAccountTokenTrueAndDefaultSAAllowed
   Message: Audit result overridden: Default service account with token mounted. automountServiceAccountToken should be set to 'false' or a non-default service account should be used.
   Metadata:
      OverrideReason: SomeReason
```

As per Kubernetes spec, `value` must be 63 characters or less and must be empty or begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.

Multiple override labels (for multiple auditors) can be added to the same resource.

See the specific [auditor docs](#auditors) for the auditor you wish to override for examples.

To learn more about labels, see https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/

## Contributing

If you'd like to fix a bug, contribute a feature or just correct a typo, please feel free to do so as long as you follow our [Code of Conduct](https://github.com/Shopify/kubeaudit/blob/master/CODE_OF_CONDUCT.md).

1. Create your own fork!
1. Get the source: `go get github.com/Shopify/kubeaudit`
1. Go to the source: `cd $GOPATH/src/github.com/Shopify/kubeaudit`
1. Add your forked repo as a fork: `git remote add fork https://github.com/you-are-awesome/kubeaudit`
1. Create your feature branch: `git checkout -b awesome-new-feature`
1. Install [Kind](https://kind.sigs.k8s.io/#installation-and-usage)
1. Run the tests to see everything is working as expected: `make test` (to run tests without Kind: `USE_KIND=false make test`)
1. Commit your changes: `git commit -am 'Adds awesome feature'`
1. Push to the branch: `git push fork`
1. Sign the [Contributor License Agreement](https://cla.shopify.com/)
1. Submit a PR (All PR must be labeled with :bug: (Bug fix), :sparkles: (New feature), :book: (Documentation update), or :warning: (Breaking changes) )
1. ???
1. Profit

Note that if you didn't sign the CLA before opening your PR, you can re-run the check by adding a comment to the PR that says "I've signed the CLA!"!
