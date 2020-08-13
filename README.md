[![Build Status](https://api.travis-ci.org/Shopify/kubeaudit.svg?branch=master)](https://travis-ci.org/Shopify/kubeaudit/)
[![codecov](https://codecov.io/gh/Shopify/kubeaudit/branch/master/graph/badge.svg)](https://codecov.io/gh/Shopify/kubeaudit)
[![Go Report Card](https://goreportcard.com/badge/github.com/Shopify/kubeaudit)](https://goreportcard.com/report/github.com/Shopify/kubeaudit)
[![GoDoc](https://godoc.org/github.com/Shopify/kubeaudit?status.png)](https://godoc.org/github.com/Shopify/kubeaudit)

> Kubeaudit can now be used as both a command line tool (CLI) and as a Go package!

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
* [Commands](#commands)
* [Configuration File](#configuration-file)
* [Override Errors](#override-errors)
* [CI/CD Usage](#cicd-usage)
* [Contributing](#contributing)

## Installation
 
### Download a binary

Kubeaudit has official releases that are blessed and stable:
[Official releases](https://github.com/Shopify/kubeaudit/releases)

### DIY build

Master will have newer features than the stable releases. If you need a newer
feature not yet included in a release you can do the following to get
kubeaudit:

**For go 1.12 and higher:**
```sh
GO111MODULE=on go get -v github.com/Shopify/kubeaudit
```

**For older versions of go:**
```sh
git clone https://github.com/Shopify/kubeaudit.git
cd kubeaudit
make
make install
```

Start using `kubeaudit` with the [Quick Start](#quick-start) or view all the [supported  commands](#commands).

### Kubectl Plugin

Prerequisite: kubectl v1.12.0 or later

With kubectl v1.12.0 introducing [easy pluggability](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) of external functions, kubeaudit can be invoked as `kubectl audit` by

- running `make plugin` and having `$GOPATH/bin` available in your path.

or

- renaming the binary to `kubectl-audit` and having it available in your path.

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

--------- Results for ---------------------

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

-- [error] CapabilityNotDropped
   Message: Capability not dropped. Ideally, the capability drop list should include the single capability 'ALL' which drops all capabilities.
   Metadata:
      Container: container
      Capability: AUDIT_WRITE
...
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

### Local Mode

If a kubeconfig file is provided using the `-c/--kubeconfig` flag, kubeaudit will audit the resources specified in the kubeconfig file. If no kubeconfig file is specified, `$HOME/.kube/config` is used by default:

```
kubeaudit all -c "/path/to/config"
```

For more information on kubernetes config files, see https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/

### Cluster Mode

Kubeaudit can detect if it is running within a container in a cluster. If so, it will try to audit all Kubernetes resources in that cluster:
```
kubeaudit all
```

## Commands

| Command          | Description                                                  | Documentation                     |
| :--------------- | :----------------------------------------------------------- | :-------------------------------- |
| `all`            | Runs all available auditors, or those specified using a kubeaudit config. | [docs](docs/all.md)  |
| `autofix`        | Automatically fixes security issues.                         | [docs](docs/autofix.md)           |

### Auditors

Auditors can also be run individually.

| Command          | Description                                                              | Documentation                     |
| :--------------- | :----------------------------------------------------------------------- | :-------------------------------- |
| `apparmor`       | Finds containers running without AppArmor.                               | [docs](docs/auditors/apparmor.md) |
| `asat`           | Finds pods using an automatically mounted default service account        | [docs](docs/auditors/asat.md) |
| `capabilities`   | Finds containers that do not drop the recommended capabilities or add new ones. | [docs](docs/auditors/capabilities.md) |
| `hostns`         | Finds containers that have HostPID, HostIPC or HostNetwork enabled.      | [docs](docs/auditors/hostns.md) |
| `image`          | Finds containers which do not use the desired version of an image (via the tag) or use an image without a tag. | [docs](docs/auditors/image.md) |
| `limits`         | Finds containers which exceed the specified CPU and memory limits or do not specify any. | [docs](docs/auditors/limits.md) |
| `mountds`        | Finds containers that have docker socket mounted.                        | [docs](docs/auditors/mountds.md) |
| `netpols`        | Finds namespaces that do not have a default-deny network policy.         | [docs](docs/auditors/netpols.md) |
| `nonroot`        | Finds containers running as root.                                        | [docs](docs/auditors/nonroot.md) |
| `privesc`        | Finds containers that allow privilege escalation.                        | [docs](docs/auditors/privesc.md) |
| `privileged`     | Finds containers running as privileged.                                  | [docs](docs/auditors/privileged.md) |
| `rootfs`         | Finds containers which do not have a read-only filesystem.               | [docs](docs/auditors/rootfs.md) |
| `seccomp`        | Finds containers running without Seccomp.                                | [docs](docs/auditors/seccomp.md) |

### Global Flags

| Short   | Long           | Description                                                                                         |
| :------ | :------------- | :-------------------------------------------------------------------------------------------------- |
| -j      | --json         | Output audit results in JSON                                                                        |
| -c      | --kubeconfig   | Path to local Kubernetes config file. Only used in local mode (default is `$HOME/.kube/config`)     |
| -f      | --manifest     | Path to the yaml configuration to audit. Only used in manifest mode.                                |
| -n      | --namespace    | Only audit resources in the specified namespace. Only used in cluster mode.                         |
| -m      | --minseverity  | Set the lowest severity level to report (one of "ERROR", "WARN", "INFO") (default "INFO")           |

## Configuration File

Kubeaudit can be used with a configuration file instead of flags. See the [all command](docs/all.md).

## Override Errors

Security issues can be ignored for specific containers or pods by adding override labels. This means the auditor will produce `warning` results instead of `error` results. The labels are documented in each auditor's documentation, but the general format for auditors that support overrides is as follows:

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

If the `value` is set to a non-empty string, it will be displayed in the `warning` result as the `OverrideReason`:
```
-- [warning] AutomountServiceAccountTokenTrueAndDefaultSAAllowed
   Message: Default service account with token mounted. automountServiceAccountToken should be set to 'false' or a non-default service account should be used.
   Metadata:
      OverrideReason: SomeReason
```

As per Kubernetes spec, `value` must be 63 characters or less and must be empty or begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.

Multiple override labels (for multiple auditors) can be added to the same resource.

See the specific [auditor docs](#auditors) for the auditor you wish to override for examples.

To learn more about labels, see https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/

## CI/CD usage

kubeaudit will return exit code `2` whenever any errors are being found, so it can stop your pipeline.
If you do not want this to happen, run it as `kubeaudit all || true`

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
