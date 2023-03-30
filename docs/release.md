# How to Create a Kubeaudit Release

1. Make sure you're on main and have the latest changes.

2. Find the [latest release](https://github.com/Shopify/kubeaudit/releases).

> We use [semver](https://semver.org/) versioning. In semver, version numbers have the format `v<MAJOR>.<MINOR>.<PATCH>`. However, because we still consider Kubeaudit to be in "alpha", the major number is always 0. This means that we do not maintain versions from before a breaking change, and updating to a new minor version may introduce a breaking change.

If the changes since the most recent release are bug fixes only, bump the last number (the patch version). If any of the changes since the last release include a new feature or breaking change, bump the second number (the minor version) and set the last number to 0 (the patch version). For example, if the latest release is `v0.11.5` and there were only bug fixes merged to main since then, the next version number will be `v0.11.6`. If there were new features added or a breaking change was made, the next version would be `v0.12.0`.

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


5. Once you push the tag, the release Github action will be triggered and generate a draft release in Github, allowing you to double check it and make changes to the Changelog. Find the [draft release](https://github.com/Shopify/kubeaudit/releases) and make sure there are no commits to main since the release.

> If there are commits to main since the release, this may mean you didn't make the tag on main or your main is out of date.

6. Click `Edit` on the right of the draft release and tidy up the Changelog if necessary. We like to add thank you's to external contributors, for example:

```
202e355 Fixed code quality issues using DeepSource (#315) - Thank you @withshubh for the contribution!
```

Optionally, you can click on "Generate release notes", which adds Markdown for all the merged pull requests from the diff and contributors of the release.

7. Click on `Publish release` at the bottom.
