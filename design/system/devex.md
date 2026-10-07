---
title: Developer experience and operations
updated: 2026-10-07
status: active
topics: [code, template]
---

# Developer experience and operations

What the template ships so a conforming repo is pleasant to work in from the first commit, and what this monorepo uses on top.

## Shipped by the template

| Component | File | Purpose |
|-----------|------|---------|
| Editor config | `.editorconfig` | UTF-8, LF, final newline, 2-space indent for yaml/md/json, tabs for Go and Makefile |
| Git ignore | `.gitignore` | OS and editor noise, `node_modules`, Go build output, `.flai-cache/` |
| Git attributes | `.gitattributes` | LF normalisation, markdown treated as text with diff, lockfiles marked generated |
| Markdown lint | `.markdownlint.yaml` | Relaxed line length, consistent headings and lists |
| Make | `Makefile` | `make check` runs `flai check` and markdownlint, `make dashboard` runs `flai dashboard`, `make board`, `make stats` |
| CI | `.github/workflows/system-flow-check.yml` | Runs `flai check --strict` and markdownlint on pull requests |
| Agent instructions | `CLAUDE.md` | The baseline described in [template.md](template.md) |
| Pull request template | `.github/pull_request_template.md` | Asks for the story ID and confirms the definition of done |

## Used by this monorepo

| Concern | Choice | Detail |
|---------|--------|--------|
| Go build and lint | `go build`, `golangci-lint` | Config in `flai/.golangci.yaml`; `scripts/flai-test.sh` names the golangci-lint it runs and refuses one that is not v2 |
| Tools and environment | `scripts/env.sh`, sourced by every script | `PATH` holds the tree's `bin/`, then the main checkout's `bin/`, where `scripts/install-tools.sh` installs the pinned golangci-lint and GoReleaser from any tree, so a story worktree uses them (S-0183); caches and `FLAI_CONFIG` default to the main checkout's `.flai-cache/` |
| Tier commands | `scripts/with-env.sh` | Runs one command with the tree's environment, as `env.sh` sets it, so a tier in `tests` reaches the golangci-lint `install-tools.sh` put in `bin/` though flai was started with another `PATH`; `scripts/flaiover-unit.sh` and `scripts/lint-md.sh` take the files a tier selected as arguments |
| Agent environment | `flai serve` | An agent gets flai serve's environment without `FLAI_HOST_URL`, `FLAI_HOST_TOKEN`, and `FLAI_CONFIG`, plus its story's variables ([flai-cli.md](flai-cli.md#commands)): it reaches neither the operator's host nor the config flai serve was given |
| Dashboard development | `scripts/flaiover-dev.sh`, `flaiover/compose.yaml` | Both take the login token from the file flai writes beside flai serve's state (`flai dashboard token --json`, `file`); compose mounts that file alone, never the repository (ADR-0031) |
| Test tiers | `system-flow.yaml`'s `tests`, `flai test`; `make test`, `make integration`, `make smoke` | The manifest's `tests` declare the tiers ([project-manifest.md](project-manifest.md)): gofmt, vet, golangci-lint, the short Go tests, vitest, and the markdown lint, each selected by the paths it covers, then integration and smoke under `--all` only. `flai test <path>...` runs the tiers the paths select, cheapest first, stops at the first that fails, and answers pass or its first findings, capped by `--max`; with no paths, what the checkout changed against the main branch. The MCP tool `test` and the host method `test.run` answer the same ([flai-cli.md](flai-cli.md#commands)). An agent runs it between tasks (S-0273). The `make` targets stay for people and CI: behavior tests beside the code (`-short`), integration tests need real git or the monorepo, smoke renders the template and checks the repo. `make test` also runs flaiover's vitest (`scripts/flaiover-unit.sh`) when `flaiover/node_modules` is present. `make flai-test`, and so close-out of a story that changes `flai/`, runs the lint, vitest, `make integration`, and `make smoke`, not `make test`: the full Go run holds every short test, so the Go tests run once, with the race detector (S-0267). See `design/conventions/code-quality.md`. |
| Go release | GoReleaser | `flai/.goreleaser.yaml`, tags `flai/v*` produce GitHub release binaries for linux, darwin, windows |
| Node package manager | pnpm | Lockfile committed, `corepack` pins the version |
| Dashboard image | GitHub Actions, `docker/build-push-action` | Tags `flaiover/v*` and `main` push to GHCR with `latest` and semver tags |
| Standard conformance | `.github/workflows/system-flow-check.yml` | Builds `flai` from source, runs `flai check --strict`, renders `./template` and checks the result |
| Task runner | Root `Makefile` | Delegates into `flai/` and `flaiover/` |
| Versioning | Independent per sub-project | Git tags prefixed with the project name |

Details and versions are in `design/tech/`.
