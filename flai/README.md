# flai

The system-flow CLI. Design: [design/system/flai-cli.md](../design/system/flai-cli.md). Libraries: [design/tech/go-libraries.md](../design/tech/go-libraries.md). User docs: [docs/users/flai.md](../docs/users/flai.md).

```bash
go build -o flai .              # local build
go test -race ./...             # tests
golangci-lint run               # lint (golangci-lint v2)
go run . version
```

Release builds inject version metadata:

```bash
go build -ldflags "-X github.com/bytepunx/system-flow/flai/internal/buildinfo.Version=0.1.0 \
  -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Commit=$(git rev-parse --short HEAD) \
  -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" .
```

Layout follows the internal structure in the design document: `cmd/` has one file per command, `internal/` has one package per concern.
