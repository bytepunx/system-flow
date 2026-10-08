---
id: TH-0356
title: "Plan: S-0322 waits for S-0312's fix, tests I-0108's three-task case, and closes I-0108"
anchor:
  path: wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md
  item: S-0322
status: resolved
participants: [planner-S-0322, orchestrator]
created: 2026-10-08T04:33:28Z
updated: 2026-10-08T04:36:37Z
---

# TH-0356 Plan: S-0322 waits for S-0312's fix, tests I-0108's three-task case, and closes I-0108

On wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md.

## Entries

### 2026-10-08T04:33:28Z planner-S-0322
Recommendation: cancel S-0322 as a duplicate of S-0312, as TH-0306 proposes, and have S-0312's T-1285 close I-0108 with I-0104. If you keep S-0322, the plan below is ready to finalize.

I-0108 is the fault I-0104 records. S-0312 is ready and planned to fix it: T-1283 replaces `git add -A` in `flai/internal/taskdone/taskdone.go` with a commit of the closing task's paths. So S-0322 builds no fix of its own.

## Plan

S-0322 waits for S-0312 (`after`).

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1328: test three tasks of one layer closed in turn, each with a commit of its own (I-0108's S-0274 case), or name an S-0312 test that does | `flai/internal/taskdone/taskdone_test.go` |
| 2 | T-1329: `flai issue close I-0108`, after T-1328 | the I-0108 file, `design/issues/summary.md` |

| Figure | Value | Why |
|--------|-------|-----|
| Forecast | 20m, delivery 2026-10-08T12:13:00Z | flai gave 9m; raised for reading S-0312's change and the Go tiers at close-out |
| Cost of delay | 12.50 USD a week | `flai cod` from the 5m input, unchanged |

All touches name files; no folder touch.

## Assumptions

- S-0312 lands as planned. If T-1283 commits only the closing task's paths, I-0108's cause is gone and criterion 1 needs only the test.
- S-0312's own test covers two tasks. A three-task test matches I-0108's instance; T-1328 adds it only if S-0312 left none.
- `taskdone.go` is not claimed, because S-0312 changes it and this story only reads it.

## If you cancel instead

Add I-0108 to T-1285 on S-0312, so that it closes I-0108 with I-0104. That needs your word on that task, since it is not mine to rewrite.

### 2026-10-08T04:33:36Z orchestrator (recommendation)
Recommendation: cancel S-0322 as a duplicate of S-0312, and have S-0312's T-1285 close I-0108 along with I-0104.

- **One cause.** I-0108 and I-0104 record the same fault. S-0312's T-1283, whose remedy you confirmed on TH-0337, removes it: a task commits only the paths its touches cover. S-0322 builds no fix of its own. It only adds a three-task test and a close.
- **What it costs.** `flai move S-0322 cancelled --reason "duplicate of S-0312"`. Then add I-0108 to T-1285 on S-0312 (`flai edit T-1285`, or tell S-0312's agent when it starts). Optionally, ask for a three-task case in T-1283's test.
- **If you keep it.** The plan is sound. It waits for S-0312 by `after`, has two file-touch tasks, a 20m forecast, and a 12.50 USD a week value. I will finalize it on your word.

Changing S-0312's task and cancelling S-0322 are both yours, so this is a recommendation.

Source: wip/kanban/stories/S-0312-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md

### 2026-10-08T04:33:54Z orchestrator (recommendation)
Recommendation, updated: operator, at 04:33Z you cancelled S-0312, not S-0322. If you meant S-0322, undo it and cancel S-0322 instead:

```sh
flai move S-0312 backlog
flai move T-1282 backlog
flai move T-1283 backlog
flai move T-1284 backlog
flai move T-1285 backlog
flai move S-0322 cancelled --reason "duplicate of S-0312"
```

ADR-0055 lets a cancelled item move back to backlog, so this is reversible. S-0312 then needs your finalizing and promoting again. I can do both once it is back.

Why it matters: S-0322's plan builds no fix. It waits for S-0312 (`after: S-0312`) and only adds a test and closes I-0108. With S-0312 cancelled, nothing fixes I-0104 or I-0108, and S-0322 waits on a story that will never be done.

If you meant to keep S-0322 instead: say so here. The planner then needs to fold S-0312's four tasks into S-0322, with your confirmed remedy from TH-0337, and drop its `after`. I would ask it to replan.

Source: design/adrs/0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md § Decision

### 2026-10-08T04:35:05Z orchestrator (recommendation)
Recommendation: run `flai plan S-0322` (or Plan on the dashboard) so the planner replans it as the story that carries the fix. I cannot start it: `flai plan --candidates` does not list S-0322, because it already has tasks.

You kept S-0322 and added I-0104 to its criteria. Its plan still assumes S-0312 does the fix:

- `after: S-0312` points at a cancelled story. promotion waits for every story in `after` to be done, so S-0322 would be held for good.
- T-1328 and T-1329 only add a test and close I-0108. Nothing changes `flai/internal/taskdone/taskdone.go`, and nothing closes I-0104.

The replan should drop the `after`. It should take over S-0312's four tasks with the remedy you confirmed on TH-0337: an ADR refining ADR-0107, a commit of only the closing task's paths, the CLI and MCP changes, and the docs. It should add the I-0104 close, and set the touches, forecast, and value again. I will approve that plan and finalize S-0322 when it comes back.

Separately, criterion 2 shows `flai issue close I-0104 --reason` without the backticks criteria 1 and 3 use. That is cosmetic.

Source: wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md

### 2026-10-08T04:36:37Z orchestrator
Resolved: Settled by the operator: S-0312 cancelled, S-0322 kept and widened to I-0104, and replanned on TH-0358 to carry the fix
