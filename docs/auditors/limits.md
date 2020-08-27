# Limits Auditor (limits)

Finds containers which exceed the specified CPU and memory limits or do not specify any.

## General Usage

```
kubeaudit limits [flags]
```

### Flags
| Short   | Long      | Description                                               | Default                          |
| :------ | :-------- | :-------------------------------------------------------- | :------------------------------- |
|         | --cpu     | Max CPU limit                                             |                                  |
|         | --memory  | Max memory limit                                          |                                  |

Also see [Global Flags](/README.md#global-flags)

## Examples

The max CPU is specified using the `--cpu` flag:
```
$ kubeaudit limits --cpu 600m -f "auditors/limits/fixtures/resources-limit.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod

--------------------------------------------

-- [warning] LimitsCPUExceeded
   Message: CPU limit exceeded. It is set to '750m' which exceeds the max CPU limit of '600m'.
   Metadata:
      Container: container
      ContainerCpuLimit: 750m
      MaxCPU: 600m
```

The max memory is specified using the `--memory` flag:
```
$ kubeaudit limits --memory 384 -f "auditors/limits/fixtures/resources-limit.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod

--------------------------------------------

-- [warning] LimitsMemoryExceeded
   Message: Memory limit exceeded. It is set to '512Mi' which exceeds the max Memory limit of '384'.
   Metadata:
      MaxMemory: 384
      Container: container
      ContainerMemoryLimit: 512Mi
```

The CPU and memory can be audited at the same time by including both the `--cpu` and `--memory` flags:
```
$ kubeaudit limits --cpu 600m --memory 384 -f "auditors/limits/fixtures/resources-limit.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod

--------------------------------------------

-- [warning] LimitsCPUExceeded
   Message: CPU limit exceeded. It is set to '750m' which exceeds the max CPU limit of '600m'.
   Metadata:
      Container: container
      ContainerCpuLimit: 750m
      MaxCPU: 600m

-- [warning] LimitsMemoryExceeded
   Message: Memory limit exceeded. It is set to '512Mi' which exceeds the max Memory limit of '384'.
   Metadata:
      Container: container
      ContainerMemoryLimit: 512Mi
      MaxMemory: 384
```

The `limits` auditor can be used to find all containers which do not specify a max CPU or memory by omitting the `--cpu` and `--memory` flags:
```
$ kubeaudit limits  -f "auditors/limits/fixtures/resources-limit-nil.yml"

---------------- Results for ---------------

  apiVersion: v1
  kind: Pod
  metadata:
    name: pod

--------------------------------------------

-- [warning] LimitsNotSet
   Message: Resource limits not set.
   Metadata:
      Container: container
```

## Override Errors

Overrides are not currently supported for `limits`.
