---
id: D-0636
title: flai task done's text report prints the sync through printSyncChecks, which reads each branch's state on the baseline workflow, while flai stream sync calls printSyncChecksOn with the repository's workflow
kind: decision
sources: [S-0744, S-0755]
topics: [cli, go, flai, git]
paths: [flai/cmd/stream_sync.go, flai/cmd/task_done.go, flai/cmd/stream.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0636 flai task done's text report prints the sync through printSyncChecks, which reads each branch's state on the baseline workflow, while flai stream sync calls printSyncChecksOn with the repository's workflow

This entry no longer holds. Since S-0755, `flai task done`'s text report reads each branch's state on the repository's workflow, as `flai stream sync` does.

S-0744 moved `flai stream sync` onto the project's workflow (ADR-0301) but left `flai task done`'s report on the baseline, since `flai/cmd/task_done.go` was not its file. Its narrative's Decisions (`wip/archive/agents/S-0744.md`):

> printSyncChecks keeps its signature on the baseline for cmd/task_done.go (not this story's file), and flai stream sync uses printSyncChecksOn with the repository's workflow.

S-0755's task T-3390, "The MCP server, flai serve, and the commands that still read the baseline run the project's workflow" (commit 398e7b040), switched `flai/cmd/task_done.go` to `printSyncChecksOn(w, wf, ...)` with the repository's workflow. In `flai/cmd/stream_sync.go`, `printSyncChecks` still calls `printSyncChecksOn` with `workflow.Baseline()`. Only `flai/cmd/stream_sync_test.go` calls it now, so a story may drop it and point those tests at `printSyncChecksOn`. D-0592 covers the rest of what a new sync field needs in these files.

Sources: S-0744's narrative, `wip/archive/agents/S-0744.md`, § Decisions; S-0755's narrative, `wip/archive/agents/S-0755.md`, § Decisions (T-3390).
