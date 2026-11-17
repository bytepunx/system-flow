---
id: T-3461
type: task
nature: feature
title: flai upgrade writes the lines a project's queue and when stand for into the workflow file it adds to a project that has none
status: backlog
parent: S-0813
owner: alex
created: 2026-11-16T23:49:49Z
updated: 2026-11-16T23:49:49Z
transitions: []
stream: S-0813
tags: [flai, template]
touches: [flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go]
after: [T-3452]
---
# T-3461 flai upgrade writes the lines a project's queue and when stand for into the workflow file it adds to a project that has none

## Work

The title's "flai upgrade writes them into a project's file", within ADR-0307: `flai upgrade` adds the baseline file only to a project that lacks it, and never writes a project's file after. It waits for T-3452's derivation.

- In `flai/internal/upgrade/upgrade.go`, where `kept` and `Compute` handle the workflow file: when the project has no `design/workflows/board.md` and its manifest sets `queue` or `when`, write the template's baseline with T-3452's lines in place of the baseline's `verify:` and `verify_at.batch`. Record in the lock what the template gave, so that a later upgrade still keeps the file as the project's own.
- A project that has its own file is untouched, keys or none. `flai check`'s advice, T-3455, is its way to the lines.
- `Describe` says the file was added with the lines the keys stand for.

## Done when

- `upgrade_test.go` covers a project with no file and the keys (file added with the lines), with no file and no keys (baseline as today), and with its own file and the keys (file unchanged).
- `flai test` passes for `flai/internal/upgrade`.

## Notes
