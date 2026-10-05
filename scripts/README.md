# scripts

Purpose-named shell scripts for common tasks. The `Makefile` calls these; CI calls the same ones. Each script runs from any directory, prints what it does, and exits non-zero on failure. POSIX sh only; the operator's interactive shell is zsh.

| Script | Does |
|--------|------|
| `env.sh` | Sourced by the others: puts Go, `~/go/bin`, and `./bin` on PATH, sets `FLAI_CONFIG` to `.flai-cache/config.json` |
| `flai.sh` | Runs `bin/flai`, building it from `flai/` first if missing or stale |
| `flai-build.sh` | Builds `bin/flai` from source |
| `flai-reference.sh` | Regenerates `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md` from the command help (`go run`, leaves `bin/flai` alone) |
| `test.sh` | Behavior tests: `go test -race -short` in `flai/`, seconds, no external dependencies, then `flaiover-unit.sh` |
| `flaiover-unit.sh` | flaiover's vitest against a current `bin/flai`, when `flaiover/node_modules` is present; nothing otherwise |
| `integration.sh` | Integration tests: full `go test -race` including real git and the monorepo round-trip |
| `smoke.sh` | Smoke tests: render the template and check it, then check this repository |
| `flai-test.sh` | gofmt, vet, golangci-lint v2, then `flaiover-unit.sh`, `integration.sh`, and `smoke.sh`: not `test.sh`, whose short Go tests the full run holds, so each Go test runs once |
| `flai-snapshot.sh` | GoReleaser snapshot build into `flai/dist` |
| `check.sh` | `flai check --strict` on this repository |
| `close-out.sh` | Before a story goes to review, in its worktree: the tests its branch calls for, the markdown lint, `flai check --strict`, the narrative's `## Current state` and `## Next steps`, then the commit with the `git commit` options given; stops at the first step that fails |
| `lint-md.sh` | markdownlint-cli2 over every markdown file, with the CI globs and `.markdownlint.yaml` |
| `mdlint-fixtures.sh` | Regenerates `flai/internal/mdlint/testdata/cases/expected.txt` with markdownlint-cli2, the reference flai's own markdown lint is tested against |
| `template-test.sh` | Renders `template/` into a temp dir and checks the result |
| `install-tools.sh` | Installs golangci-lint v2 and GoReleaser into `bin/`, pnpm into `.flai-cache/pnpm` |
| `flaiover-install.sh` | `pnpm install --frozen-lockfile` in `flaiover/` |
| `flaiover-dev.sh` | Dev server against this repository (`PROJECT_DIR` defaults to the root) |
| `flaiover-build.sh` | Production build (adapter-node) |
| `flaiover-test.sh` | prettier, eslint, svelte-check, vitest |
| `flaiover-image.sh` | Builds the flaiover image locally as `flaiover:local` from the repo root |
