---
id: I-003
title: Files needed by tests or checks were missing in a fresh clone
class: defect
status: open
count: 2
cost: 15m
first_reported: 2026-09-15T18:21:56Z
last_reported: 2026-09-15T18:21:56Z
updated: 2026-09-15T22:40:33Z
---

# I-003 Files needed by tests or checks were missing in a fresh clone

## Description
The working tree passed but a clone did not: a test fixture under `bin/` was excluded by the root `.gitignore`, and `wip/kanban/tasks` was empty after archiving so git dropped the folder and `flai check` failed on the clone.

## Instances
### 2026-09-15T18:21:56Z
S-010: GoReleaser dry run in a clone ran `go test` and failed on the missing fixture; renamed `bin/` to `tools/`.

### 2026-09-15T18:21:56Z
S-010: same run, `flai check` reported `kanban/tasks` missing; added `.gitkeep` files.

## Remediation
Convention now requires verifying by clone when unsure (code-quality.md). A CI job on a clean checkout will catch the rest once the repo is on GitHub. Close after the first green CI run.
