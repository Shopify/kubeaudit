# How to Create a Kubeaudit Release

1. Make sure you're on master and have the latest changes.

2. Find the [latest release](https://github.com/Shopify/kubeaudit/releases).

> We use [semver](https://semver.org/) versioning. In semver, version numbers have the format `v<MAJOR>.<MINOR>.<PATCH>`. However, because we still consider Kubeaudit to be in "alpha", the major number is always 0. This means that we do not maintain versions from before a breaking change, and updating to a new minor version may introduce a breaking change.

If the changes since the most recent release are bug fixes only, bump the last number (the patch version). If any of the changes since the last release include a new feature or breaking change, bump the second number (the minor version) and set the last number to 0 (the patch version). For example, if the latest release is `v0.11.5` and there were only bug fixes merged to master since then, the next version number will be `v0.11.6`. If there were new features added or a breaking change was made, the next version would be `v0.12.0`.

3. Update the `VERSION` file if necessary. You'll have to open / merge a PR to do this.

4. Create a tag with the new version and push it up to Github:

```
git tag -a <VERSION> -m "<VERSION>"
git push origin <VERSION>
```

For example:

```
git tag -a v0.11.6 -m "v0.11.6"
git push origin v0.11.6
```

5. Login to Docker

Login to the Shopify Docker account using `docker login`. Find the credentials in password manager.

6. You will need a Github token in order for Goreleaser to be able to create a release in Github. If you already have one, skip to the next step.

[Create a Github token](https://github.com/settings/tokens/new) with the `repo` scope.

7. Run Goreleaser

Make sure you have Docker running.

```
GITHUB_TOKEN=<YOUR TOKEN> goreleaser --rm-dist
```

8. Logout of docker

Logout of the Shopify Docker account: `docker logout`

9. Publish the release in Github

Goreleaser is set to draft mode which means it will create a draft release in Github, allowing you to double check the release and make changes to the Changelog. Find the [draft release](https://github.com/Shopify/kubeaudit/releases) and make sure there are no commits to master since the release.

> If there are commits to master since the release, this may mean you didn't make the tag on master or your master is out of date.

Click `Edit` on the right of the draft release and tidy up the Changelog if necessary. We like to add thank you's to external contributors, for example:

```
202e355 Fixed code quality issues using DeepSource (#315) - Thank you @withshubh for the contribution!
```

Click on `Publish release` at the bottom.
