---
id: T-2363
type: task
nature: remediation
title: The go-test tier names its timeout, and flai check warns of a tier that runs go test without -timeout
status: in-progress
parent: S-0548
owner: alex
created: 2026-10-31T07:25:01Z
updated: 2026-10-31T07:32:12Z
transitions:
  - to: ready
    at: 2026-10-31T07:32:11Z
    by: agent-S-0548
  - to: in-progress
    at: 2026-10-31T07:32:12Z
    by: agent-S-0548
stream: S-0548
tags: []
touches: [system-flow.yaml, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/check.go, flai/internal/check/tests_test.go, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, template/CHANGELOG.md]
---
# T-2363 The go-test tier names its timeout, and flai check warns of a tier that runs go test without -timeout

## Work

The `go-test` tier in `system-flow.yaml` runs `go test -race` with no `-timeout`, so go's 10-minute default panics naming only the package (I-0227). The `integration` tier's `scripts/integration.sh` passes `-timeout 20m`, and both `go test` lines in `.github/workflows/flai.yml` pass `-timeout 30m` already. Give the `go-test` tier `-timeout 20m`, and have `flai check` warn, under a rule name in the style of the others, of each tier whose `command` or `all_command` runs `go test` with no `-timeout` argument, naming the tier and saying to add one; a tier that runs a script is not inspected. Document the warning where `manifest.tests` is documented. Waits for no other task.

## Done when

- [ ] A test fails without the warning and passes with it; a tier with `-timeout` and a script tier get none
- [ ] `flai check --strict` passes on this repository

## Notes
