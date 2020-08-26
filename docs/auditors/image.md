# Image Auditor (image)

Finds containers which do not use the desired version of an image (via the tag) or use an image without a tag.

## General Usage

```
kubeaudit image [flags]
```

### Flags
| Short   | Long      | Description                                               | Default                          |
| :------ | :-------- | :-------------------------------------------------------- | :------------------------------- |
| -i      | --image   | Image and tag to check against.                           |                                  |

Also see [Global Flags](/README.md#global-flags)

## Examples

The image and tag to look for are specified using the `-i/--image image:tag` flag. For example, `-i gcr.io/google_containers/echoserver:1.7` will look for containers using the `gcr.io/google_containers/echoserver` image which have a tag other than `1.7`.

```
$ kubeaudit image -i "scratch:1.6" -f "auditors/image/fixtures/image-tag-present.yml"

---------------- Results for ---------------

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: deployment

--------------------------------------------

-- [error] ImageTagIncorrect
   Message: Container tag is incorrect. It should be set to '1.6'.
   Metadata:
      Container: deployment
```

If the container image matches the provided image but the container image has no tag, a warning is produced:
```
$ kubeaudit image -i "scratch:1.6" -f "auditors/image/fixtures/image-tag-missing.yml"

---------------- Results for ---------------

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: deployment

--------------------------------------------

-- [warning] ImageTagMissing
   Message: Image tag is missing.
   Metadata:
      Container: container
```

The `image` auditor can be used to find all containers that use an image without a tag by omitting the `-i/--image` flag:
```
$ kubeaudit image -f "auditors/image/fixtures/image-tag-missing.yml"

---------------- Results for ---------------

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: deployment

--------------------------------------------

-- [warning] ImageTagMissing
   Message: Image tag is missing.
   Metadata:
      Container: container
```

## Override Errors

Overrides are not currently supported for `image`.
