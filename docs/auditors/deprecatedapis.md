# Kubernetes Deprecated API Auditor (deprecatedapis)

Finds any resource defined with a deprecated API version.

## General Usage

```
kubeaudit deprecatedapis [flags]
```

### Flags
| Short   | Long                   | Description                                   | Default             |
| :------ | :--------------------- | :-------------------------------------------- | :------------------ |
|         | --current-k8s-version  | Kubernetes current version                    |                     |
|         | --targeted-k8s-version | Kubernetes version to migrate to              |                     |


Also see [Global Flags](/README.md#global-flags)

## Examples

The `deprecatedapis` auditor finds the deprecated APIs in use, reports the versions where they will be removed, and recommends replacement APIs.
```
$ kubeaudit deprecatedapis -f "auditors/deprecatedapis/fixtures/cronjob.yml"

---------------- Results for ---------------

  apiVersion: batch/v1beta1
  kind: CronJob
  metadata:
    name: hello

--------------------------------------------

-- [warning] DeprecatedAPIUsed
   Message: batch/v1beta1 CronJob is deprecated in v1.21+, unavailable in v1.25+, introduced in v1.8+; use batch/v1 CronJob
   Metadata:
      DeprecatedMajor: 1
      DeprecatedMinor: 21
      IntroducedMajor: 1
      IntroducedMinor: 8
      RemovedMajor: 1
      RemovedMinor: 25
      ReplacementKind: CronJob
      ReplacementGroup: batch/v1
```

The `deprecatedapis` auditor can be used with the `--current-k8s-version` flag. If the API is not yet deprecated for this version the auditor will produce an `info` otherwise a `warning`.
```
$ kubeaudit deprecatedapis --current-k8s-version 1.20  -f "auditors/deprecatedapis/fixtures/cronjob.yml"

---------------- Results for ---------------

  apiVersion: batch/v1beta1
  kind: CronJob
  metadata:
    name: hello

--------------------------------------------

-- [info] DeprecatedAPIUsed
   Message: batch/v1beta1 CronJob is deprecated in v1.21+, unavailable in v1.25+, introduced in v1.8+; use batch/v1 CronJob
   Metadata:
      DeprecatedMajor: 1
      DeprecatedMinor: 21
      IntroducedMajor: 1
      IntroducedMinor: 8
      RemovedMajor: 1
      RemovedMinor: 25
      ReplacementKind: CronJob
      ReplacementGroup: batch/v1
```

The `deprecatedapis` auditor can be used with the `--targeted-k8s-version` flag. If the API is not available for the targeted version the auditor will produce an `error` otherwise a `warning` or `info` if the API is not yet deprecated for this version. 
```
$ kubeaudit deprecatedapis --current-k8s-version 1.20 --targeted-k8s-version 1.25 -f "auditors/deprecatedapis/fixtures/cronjob.yml"

---------------- Results for ---------------

  apiVersion: batch/v1beta1
  kind: CronJob
  metadata:
    name: hello

--------------------------------------------

-- [error] DeprecatedAPIUsed
   Message: batch/v1beta1 CronJob is deprecated in v1.21+, unavailable in v1.25+, introduced in v1.8+; use batch/v1 CronJob
   Metadata:
      DeprecatedMajor: 1
      DeprecatedMinor: 21
      IntroducedMajor: 1
      IntroducedMinor: 8
      RemovedMajor: 1
      RemovedMinor: 25
      ReplacementKind: CronJob
      ReplacementGroup: batch/v1
```

## Override Errors

Overrides are not currently supported for `deprecatedapis`.
