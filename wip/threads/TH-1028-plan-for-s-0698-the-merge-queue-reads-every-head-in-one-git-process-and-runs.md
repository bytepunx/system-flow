---
id: TH-1028
title: "Plan for S-0698: the merge queue reads every head in one git process and runs git status only while a story in review passed verify; three tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0698-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md
  item: S-0698
status: resolved
participants: [planner-S-0698, orchestrator]
created: 2026-11-10T10:31:51Z
updated: 2026-11-10T11:01:48Z
---

# TH-1028 Plan for S-0698: the merge queue reads every head in one git process and runs git status only while a story in review passed verify; three tasks in three layers

On wip/kanban/stories/S-0698-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md.

## Entries

### 2026-11-10T10:31:51Z planner-S-0698
S-0698 is planned: three tasks in three layers, 1h forecast, 12.50 USD a week cost of delay.

Proposed fix (I-0333 asked for one before building). `mergequeue.ComputeAt` runs one `git rev-parse` per story in review, plus two main-head reads, one after another: 21 to 81 ms. It also runs `git status --porcelain` on every call: 52 to 148 ms. The fix:

- read every head in one `git for-each-ref`;
- run `git status` only when a story in review has a passing verify result, and run it at the same time as that read.

The queue's reasons and window do not change.

The tasks, in their layers:

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3000: one-process heads read in storygit | `flai/internal/storygit/storygit.go`, `storygit_test.go` | none |
| 2 | T-3003: `ComputeAt` uses it and gates and overlaps `git status`; a git-process count test reproduces I-0333 | `flai/internal/mergequeue/eligible.go`, `eligible_test.go` | T-3000 |
| 3 | T-3007: measure `board.get` before and after, record it in `server-performance.md` and `flai-cli.md`, close I-0333 | `design/system/server-performance.md`, `design/system/flai-cli.md`, the I-0333 file, `design/issues/summary.md` | T-3003 |

Assumptions:

- The start prompt said to read a thread named "asked". No such thread exists, and S-0698 had no earlier plan, so I planned it fresh.
- Caching `git status` is left out. flai serve's file watcher covers the manifest and `design`, `wip`, and `docs`, not the code tree, so it cannot tell when the main checkout's dirt changed.
- Caching the heads read on the main-refs watcher's word, as S-0604 did for show-ref, is left for later. It saves about 10 ms, and it would mean sharing `release`'s watch generation with `mergequeue`.
- The 13-at-once measure needs the operator's host. If the agent cannot reach `~/.flai/serve/projects.json`, T-3007 measures this repository's board alone and says so.
- If calls still miss 300 ms only because of the host's load, T-3007 asks here before ticking criterion 1, as D-0535 and D-0556 did.
- S-0698 has no epic. S-0604 and S-0636 were under E-0033. If it belongs there, the parent is the operator's to set.

Proposals: none to split, merge, or drop.

### 2026-11-10T10:59:42Z orchestrator
Approved as written: T-3000, then T-3003, then T-3007. The touches, the 1h forecast, and the 12.50 USD a week value fit the story. Both caches stay out, as you reasoned. The story stays without an epic: its parent is the operator's to set, if they want it under E-0033.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T11:01:48Z orchestrator
Resolved: Plan approved
