---
id: T-2402
type: task
nature: improvement
title: The fortnight before is recorded, overlap notices per acceptance and conflicts per week from 2026-10-15 to 2026-10-29, with how to measure the fortnight after
status: backlog
parent: S-0550
owner: alex
created: 2026-10-31T07:43:49Z
updated: 2026-10-31T07:43:49Z
transitions: []
stream: S-0550
tags: [cli]
touches: [design/system/workflow.md]
---
# T-2402 The fortnight before is recorded, overlap notices per acceptance and conflicts per week from 2026-10-15 to 2026-10-29, with how to measure the fortnight after

## Work

Criterion 4 compares the fortnight after landing with the one before. Neither figure is in the serve journal, which records host actions only (`journal.jsonl`). The overlap notices are in `.flai-cache/overlaps.jsonl` in the main checkout, written by `itemedit` and `flai accept` and trimmed to its last 500 lines. The conflicts are in `.flai-cache/conflicts.jsonl`, of kinds `trial`, `sync`, and `accept` (`flai/internal/storygit/conflicts.go`), and `flai stats --json` reports them per week under `coordination` (`flai/internal/metrics/coordination.go`, `metrics.md` § Coordination).

- Count, for 2026-10-15 to 2026-10-29, the acceptances, the overlap notices sent at acceptance and their paths, and the conflicts by kind and path, from those two files and `flai stats --json --since 2026-10-15`. Say which part of the window `overlaps.jsonl` still covers after its trim, and set the story's figures (248 notices over 97 acceptances, 62 conflicts, counted in the operator's conversation) beside yours.
- Write the figures and the commands that measure the fortnight after landing as a short section under `design/system/workflow.md` § Branches and collisions: the hub files, their notices, their conflicts, and the date the after can be read.
- Add a line to the story's Notes with `flai edit`: criterion 4 can be verified only a fortnight after acceptance, so it stays unticked at close-out, and the line says when and how it is read.

First layer: it waits for nothing. Take the figures before the story's changes land. The documentation task waits for it, since both change `workflow.md`.

## Done when

`design/system/workflow.md` holds the before figures, their sources, and the commands for the after. The story's Notes say why criterion 4 is unticked at close-out and when it can be ticked. Story criterion 4, its before half.

## Notes

Drafted by the planner of S-0550 on 2026-10-31.
