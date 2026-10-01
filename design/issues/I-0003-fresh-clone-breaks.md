---
id: I-0003
title: Files needed by tests or checks were missing in a fresh clone
class: defect
status: closed
count: 2
cost: 15m
first_reported: 2026-09-15T18:21:56Z
last_reported: 2026-09-15T18:21:56Z
updated: 2026-10-01T08:00:48Z
---

# I-0003 Files needed by tests or checks were missing in a fresh clone

## Description
The working tree passed but a clone did not: a test fixture under `bin/` was excluded by the root `.gitignore`, and `wip/kanban/tasks` was empty after archiving so git dropped the folder and `flai check` failed on the clone.

## Instances

### 2026-09-15T18:21:56Z
S-0010: GoReleaser dry run in a clone ran `go test` and failed on the missing fixture; renamed `bin/` to `tools/`.
S-0010: same run, `flai check` reported `kanban/tasks` missing; added `.gitkeep` files.

## Remediation
Convention now requires verifying by clone when unsure (code-quality.md). A CI job on a clean checkout will catch the rest once the repo is on GitHub. Close after the first green CI run.
Closed 2026-10-01T08:00:48Z: Fixed: CI runs on a clean checkout (.github/workflows/flai.yml, system-flow-check.yml) and has been green; no instance since 2026-09-15.
