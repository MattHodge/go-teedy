# go-teedy

![CI](https://github.com/MattHodge/go-teedy/actions/workflows/actions.yml/badge.svg?branch=main) [![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](https://github.com/tterb/atomic-design-ui/blob/master/LICENSEs) [![codecov](https://codecov.io/gh/MattHodge/go-teedy/branch/main/graph/badge.svg?token=MZMQ45NV95)](https://codecov.io/gh/MattHodge/go-teedy)

## Releasing

[GoReleaser](https://goreleaser.com/) is used to publish `teedy-cli`.

The release is triggered by GitHub Actions when a tag is pushed.

```bash
git tag -a v0.2.0 -m "New Release"
git push origin v0.2.0
```