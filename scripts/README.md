# scripts

Purpose-named shell scripts for common tasks. The `Makefile` calls these; CI calls the same ones. Each script runs from any directory, prints what it does, and exits non-zero on failure. POSIX sh only; the operator's interactive shell is zsh.

| Script | Does |
|--------|------|
| `env.sh` | Sourced by the others: puts Go, `~/go/bin`, and `./bin` on PATH, sets `FLAI_CONFIG` to `.flai-cache/config.json` |
| `flai.sh` | Runs `bin/flai`, building it from `flai/` first if missing or stale |
| `flai-build.sh` | Builds `bin/flai` from source |
| `flai-test.sh` | Lints (golangci-lint v2) and race-tests `flai/` |
| `flai-snapshot.sh` | GoReleaser snapshot build into `flai/dist` |
| `check.sh` | `flai check --strict` on this repository |
| `template-test.sh` | Renders `template/` into a temp dir and checks the result |
| `install-tools.sh` | Installs golangci-lint v2 and GoReleaser into `bin/` |
