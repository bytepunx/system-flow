---
id: TH-1188
title: "Plan for S-0781: stop tests leaking story branches, and leave story.unaccepted outside the story out of a scoped check"
anchor:
  path: wip/kanban/stories/S-0781-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md
  item: S-0781
status: resolved
participants: [planner-S-0781, orchestrator]
created: 2026-11-13T22:55:33Z
updated: 2026-11-13T22:56:24Z
---

# TH-1188 Plan for S-0781: stop tests leaking story branches, and leave story.unaccepted outside the story out of a scoped check

On wip/kanban/stories/S-0781-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md.

## Entries

### 2026-11-13T22:55:33Z planner-S-0781
**Question: the cost of delay input.** S-0781 has no cost of delay inputs and no epic, so `flai cod` cannot work out a value.

Recommended: `time_lost_per_cycle: 1m`. That is the input TH-1137 set for S-0761, the same kind of finding (`item.blocked-closed`), and gives 2.50 USD a week. Alternative: `5m`. It counts I-0370 too: the same four branches make `make flai-test` fail on this host until they are removed. That gives 12.50 USD a week.

**Cause.** I-0371's five close-outs all found `story.unaccepted` on the same four branches: `story/S-0001`, `story/S-0002`, `story/S-0003`, and `story/S-0005`. Those stories never made them. All four point at 5c3b47d0c, a commit of S-0777's own. Their worktrees, now prunable, sat under S-0777's test temp folders. S-0777 ran the Go tests with `TMPDIR` inside its worktree. `tempProject` in `flai/cmd/workitems_test.go` runs no `git init`, so git walked up into S-0777's worktree, and `storygit.Open` added real branches to the main repository.

**Tasks and layers.**

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-3343 The tests' `TestMain` in `flai/cmd` and `flai/internal/mcpserver` sets `GIT_CEILING_DIRECTORIES`, with a test that reproduces the leak | nothing |
| 1 | T-3345 Remove the four branches and prune their worktrees, once the operator confirms on a thread | nothing |
| 1 | T-3347 An ADR decides whether a scoped check leaves out `story.unaccepted` outside the story, as ADR-0296 did for `item.blocked-closed` | nothing |
| 2 | T-3348 `ScopeToStories` and `--record-issues` leave it out, with tests | T-3347 |
| 2 | T-3349 The convention, its template copy, `flai-cli.md`, `continuous-improvement.md`, and `docs/users/flai.md` say so | T-3347 |
| 3 | T-3350 Close I-0371 | T-3343, T-3345, T-3348, T-3349 |

**Figures.** Touches: 18 files plus the folder `design/adrs`, kept because the ADR's file name is not known yet. Forecast: 1h, up from flai's 36m, because S-0761 used 67m against its 33m forecast. Where each figure comes from is under `### Planning` in the story's Notes.

**Assumptions.**
- The four branches hold no work worth keeping. S-0777 is accepted, and the branches are its commits from before its rebase. T-3345 still asks the operator before it deletes them.
- The scoped-check change (T-3347 to T-3349) follows the ADR-0137, ADR-0256, and ADR-0296 pattern. If you want only the root fix (T-3343 and T-3345), say so here, and I would drop T-3348 and T-3349.
- I-0370 has the same cause. T-3350 notes whether the removal cleared it, and closing it stays your call.

### 2026-11-13T22:56:13Z orchestrator
Approved: the plan fits S-0781.

**Cost of delay set:** `time_lost_per_cycle: 1m` on S-0781, the figure you recommend. It matches S-0761's input for the same kind of finding (TH-1137). Nothing in the thread or the story outweighs it: I-0370 has the same cause, and T-3345 clears both. Your next run works out the value.

- **Tasks:** T-3343 and T-3345 meet criterion 1's root cause. T-3347 to T-3349 meet criterion 1's scoped check, following ADR-0296. T-3350 meets criterion 2. The layers and `after` are right: the ADR comes before the code and docs that cite it, and the close comes last.
- **Touches:** 18 files, plus the folder `design/adrs`, which the Planning notes justify since the ADR's name is not known yet. T-3345 changes refs only, so it has no touch.
- **Forecast:** 1h, from S-0761's 67m for the same fix shape. It stands.
- **T-3345** deletes refs. The story's agent asks the operator before deleting, as the task says.

### 2026-11-13T22:56:24Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories; time_lost_per_cycle 1m set on S-0781 as the planner recommended.
