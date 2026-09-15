# flai

The system-flow CLI. Design: [design/system/flai-cli.md](../design/system/flai-cli.md). Libraries: [design/tech/go-libraries.md](../design/tech/go-libraries.md). User docs: [docs/users/flai.md](../docs/users/flai.md).

```bash
go build -o flai .              # local build
go test -race ./...             # tests
golangci-lint run               # lint (golangci-lint v2)
go run . version
```

## Install

```bash
go install github.com/bytepunx/system-flow/flai@latest          # latest tagged release
go install github.com/bytepunx/system-flow/flai@flai/v0.1.0     # a specific release
```

Or download an archive from the GitHub release named `flai vX.Y.Z`: `flai_X.Y.Z_<os>_<arch>.tar.gz` (zip on Windows) for linux, darwin, and windows on amd64 and arm64, with `checksums.txt`.

## Release

Releases are cut by GoReleaser from tags named `flai/vX.Y.Z` (`.goreleaser.yaml`, workflow `release-flai.yml`). The tag prefix is stripped in templates because the open-source GoReleaser has no monorepo support; see the comment at the top of the config.

```bash
git tag flai/v0.1.0 && git push origin flai/v0.1.0   # CI builds and publishes
goreleaser release --snapshot --clean --skip=publish  # local dry run into dist/
```

Manual builds inject the same metadata:

```bash
go build -ldflags "-X github.com/bytepunx/system-flow/flai/internal/buildinfo.Version=0.1.0 \
  -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Commit=$(git rev-parse --short HEAD) \
  -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" .
```

Layout follows the internal structure in the design document: `cmd/` has one file per command, `internal/` has one package per concern.
