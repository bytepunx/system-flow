---
id: T-1069
type: task
nature: improvement
title: flai test runs the manifest's tiers for paths or packages, or every tier with --all, and prints pass or the first findings as text and --json
status: backlog
parent: S-0273
owner: alex
created: 2026-10-06T22:51:58Z
updated: 2026-10-06T22:51:58Z
transitions: []
stream: S-0273
tags: [cli, testing]
touches: [flai/cmd/test.go, flai/cmd/test_test.go, flai/cmd/root.go, system-flow.yaml, scripts/flai-test.sh, scripts/test.sh]
after: [T-1067, T-1068]
---
# T-1069 flai test runs the manifest's tiers for paths or packages, or every tier with --all, and prints pass or the first findings as text and --json

## Work

Criterion 1. This task waits for T-1067 and T-1068. It maps the manifest's `tests` onto the verify package's tiers and runs them, so it needs both.

- Add `flai test [path|package]...` in `flai/cmd/test.go` and register it in `root.go`.
  - With no argument, it takes the paths changed in the current worktree against the main branch.
  - `--all` runs every tier, the `all_only` ones included.
  - `--max <n>` changes the cap on findings.
  - It runs in the current directory's checkout, so it works the same in a story's worktree.
- Text output: one line per tier run, with its state and duration, then each finding as `path:line name: message`, then how many were left out. `--json` prints the verify package's result.
- Exit status: 0 on pass, 1 on a failing tier, and another code for a manifest or usage error, kept apart from a test failure.
- Declare this repository's tiers in `system-flow.yaml` through `flai manifest set`, cheapest first, mirroring what `scripts/flai-test.sh` and `scripts/test.sh` run:
  - gofmt
  - go vet
  - golangci-lint
  - `go test -json -race -short` for the packages selected
  - vitest for `flaiover/`
  - the markdown lint
  - integration and smoke, `all_only`
- Add a comment to each of the two scripts naming the manifest's tiers they mirror, so they stay in step. Leave the scripts' behaviour and the `Makefile` targets as they are.

## Done when

- [ ] `flai test flai/internal/manifest` runs only the Go tiers for that package, and `flai test flaiover/src/lib/x.ts` runs only the flaiover tier (end-to-end tests through the CLI with a fixture project).
- [ ] A failing Go test in the fixture gives its name, path, line, and message in text and `--json`, and no more than the cap.
- [ ] `flai test --all` in this repository passes on a clean main branch.
- [ ] `make test` and `make flai-test` still pass.
- [ ] `go test ./cmd/...`, golangci-lint, and `flai check --strict` pass.

## Notes

Drafted by the planner for S-0273.
