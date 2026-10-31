---
id: T-2394
type: task
nature: improvement
title: flai check reports an issues summary that differs from what its folder renders, and flai check --fix rewrites it
status: backlog
parent: S-0550
owner: alex
created: 2026-10-31T07:43:20Z
updated: 2026-10-31T07:43:20Z
transitions: []
stream: S-0550
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/fix.go, flai/internal/check/fix_test.go, flai/internal/check/issues_test.go, flai/cmd/check.go, flai/cmd/check_test.go]
---
# T-2394 flai check reports an issues summary that differs from what its folder renders, and flai check --fix rewrites it

## Work

`design/issues/summary.md` is already generated: `issues.Summary` and `issues.WriteSummary` (`flai/internal/issues/issues.go`) render it. `flai issue new`, `bump`, and `close` and the MCP issue tools rewrite it on every write, and `flai issue summary` writes it on demand. `flai check`'s `issues()` (`flai/internal/check/check.go`) checks only that open issues are linked and closed ones are not. It never compares the file with what the folder renders, and nothing fixes it.

- In `issues()`, compare the file with `issues.Summary` rendered from the issue files. When they differ, report an `issues.summary` finding that names `flai check --fix` and `flai issue summary`.
- Add `--fix` to `flai check` (`flai/cmd/check.go`). Put the fixers in a new `flai/internal/check/fix.go`: a table from a rule to the function that rewrites its generated file. `--fix` runs the fixers whose rule the check reported, then checks again and prints what it rewrote and what remains. It fixes generated files only and never edits a hand-written file. Register the issues summary here. The ADR index task registers `adr.index` after this task.
- Test it in `flai/internal/check/fix_test.go` and `issues_test.go`: a summary that differs from its folder is reported and then rewritten. Test the flag in `flai/cmd/check_test.go`.

First layer: it waits for nothing. The ADR index check task waits for it, since both change `check.go` and `fix.go`.

## Done when

A fixture whose `design/issues/summary.md` differs from its issue files gets an `issues.summary` finding from `flai check`. `flai check --fix` rewrites the file to what `flai issue summary` writes, and the check then reports nothing for it. `flai check --strict` on this repository is unchanged. Story criterion 1, for the issues summary.

## Notes

Drafted by the planner of S-0550 on 2026-10-31.
