---
id: T-2409
type: task
nature: improvement
title: flai check reports an ADR index behind its folder and --fix rewrites it, and a rebase that stops on the index regenerates it
status: backlog
parent: S-0550
owner: alex
created: 2026-10-31T07:44:26Z
updated: 2026-10-31T07:44:26Z
transitions: []
stream: S-0550
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/fix.go, flai/internal/check/fix_test.go, flai/internal/issues/generated.go, flai/internal/issues/generated_test.go, flai/internal/adrindex/adrindex.go, flai/internal/adrindex/adrindex_test.go]
after: [T-2394, T-2405]
---
# T-2409 flai check reports an ADR index behind its folder and --fix rewrites it, and a rebase that stops on the index regenerates it

## Work

`flai check`'s `adrIndex` (`flai/internal/check/check.go`, rule `adr.index`) warns only when a file has no row or a row has no file. A rebase that stops on the index merges it row by row through `adrindex.Merge`, wired into the rebase through `issues.FilesIn` in `flai/internal/issues/generated.go` (ADR-0157). The issues summary is in `issues.Generated`, which the rebase regenerates (ADR-0098) and `flai verify`'s patch id leaves out (ADR-0192).

- In `adrIndex`, compare the index with the render from the ADR index task. When they differ, report `adr.index` naming the rows that differ and `flai check --fix`. Register the index's fixer in `flai/internal/check/fix.go`, beside the issues summary's.
- Add the ADR index to `issues.Generated`, with the render as its `Regenerate`. Then a rebase that stops on it writes it from the folder of the rebased tree, `flai verify`'s patch id leaves it out, and the trial merge passes over it as it does now. Remove the row merge from the rebase path, and from `adrindex` if nothing else uses it, as the story's ADR decides.
- Test it in `fix_test.go`: an index behind its folder is reported and then rewritten. Test it in `generated_test.go`: a stop on the index regenerates it from both sides' ADR files.

Third layer: it waits for the issues summary task, since both change `check.go` and `fix.go`, and for the ADR index task, for the render it calls.

## Done when

`flai check` reports an index that differs from its folder, and `flai check --fix` rewrites it to what `flai adr index` writes. A rebase that stops on `design/adrs/README.md` because both sides added an ADR continues with both rows and no person's help. `flai check --strict` on this repository is clean, and `go test -race ./internal/check ./internal/issues ./internal/adrindex ./internal/storygit` passes. Story criterion 1.

## Notes

Drafted by the planner of S-0550 on 2026-10-31. If the rebase wiring lives in `flai/internal/storygit/sync.go` rather than in `generated.go`, change it there, and `flai task done` widens the touches.
