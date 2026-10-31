---
id: T-2385
type: task
nature: improvement
title: The board's queue and the dashboard's merge queue panel show the preflight hold with its files and findings and its thread
status: backlog
parent: S-0549
owner: alex
created: 2026-10-31T07:41:30Z
updated: 2026-10-31T07:41:30Z
transitions: []
stream: S-0549
tags: [flai, flaiover]
touches: [flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/components/flaiover/src/lib/components/MergeQueue.examples.json, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-2383]
---
# T-2385 The board's queue and the dashboard's merge queue panel show the preflight hold with its files and findings and its thread

## Work

Carry the preflight hold on `mergequeue.Queue` (`eligible.go`), which `OnBoard` gives `board.get`, `flai board --json`, and the MCP `board` tool: when it began, each file, line, rule, and message, and the thread once `flai serve` has opened it. S-0552 adds its window to `Queue` in the same file, so keep to one field and its filling (MS-0263).

In `MergeQueue.svelte`, add the field to `MergeQueueState` and show the hold in `warn` beside the operator's `Held by` line: that no batch is built until main is clean, each file with its findings, and the thread as a link. The dashboard's `/api/queue` passes `board.queue` through, so the route needs no change. Add a held example to `MergeQueue.examples.json` so `flai check`'s `components.examples` stays clean, and a test in `MergeQueue.svelte.test.ts`.

Describe the panel's new line in `design/system/flaiover-dashboard.md` (the `/board` row's merge queue sentences) and `docs/users/flaiover.md` (the **Merge queue.** paragraph under Board), in the same commit.

It waits for the preflight task, whose state field it reads. It shares no path with the landing task, so the two run together.

## Done when

`flai board --json` carries the hold under `queue`, the panel shows it with its files, findings, and thread, the flaiover tests and svelte-check pass, and the two documents say so. Story criterion 4 (the dashboard).

## Notes

Drafted by the planner of S-0549 on 2026-10-31.
