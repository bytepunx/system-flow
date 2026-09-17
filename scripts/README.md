# scripts

Purpose-named shell scripts for common tasks. The `Makefile` calls these; CI calls the same ones. Each script runs from any directory, prints what it does, and exits non-zero on failure. POSIX sh only; the operator's interactive shell is zsh.

| Script | Does |
|--------|------|
| `env.sh` | Sourced by the others: puts Go, `~/go/bin`, and `./bin` on PATH, sets `FLAI_CONFIG` to `.flai-cache/config.json` |
| `flai.sh` | Runs `bin/flai`, building it from `flai/` first if missing or stale |
| `flai-build.sh` | Builds `bin/flai` from source |
| `test.sh` | Behavior tests: `go test -race -short` in `flai/`, seconds, no external dependencies |
| `integration.sh` | Integration tests: full `go test -race` including real git and the monorepo round-trip |
| `smoke.sh` | Smoke tests: render the template and check it, then check this repository |
| `flai-test.sh` | gofmt, vet, golangci-lint v2, then all three tiers in order |
| `flai-snapshot.sh` | GoReleaser snapshot build into `flai/dist` |
| `check.sh` | `flai check --strict` on this repository |
| `template-test.sh` | Renders `template/` into a temp dir and checks the result |
| `install-tools.sh` | Installs golangci-lint v2 and GoReleaser into `bin/`, pnpm into `.flai-cache/pnpm` |
| `flaiover-install.sh` | `pnpm install --frozen-lockfile` in `flaiover/` |
| `flaiover-dev.sh` | Dev server against this repository (`PROJECT_DIR` defaults to the root) |
| `flaiover-build.sh` | Production build (adapter-node) |
| `flaiover-test.sh` | prettier, eslint, svelte-check, vitest |
