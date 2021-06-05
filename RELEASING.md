# Releasing

[GoReleaser](https://goreleaser.com/) is used to publish `teedy-cli`.

The release is triggered by GitHub Actions when a tag is pushed.

```bash
git tag -a v0.2.0 -m "New Release"
git push origin v0.2.0
```
