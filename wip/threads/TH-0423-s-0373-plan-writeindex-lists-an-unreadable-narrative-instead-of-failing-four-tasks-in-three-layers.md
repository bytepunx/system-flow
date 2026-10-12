---
id: TH-0423
title: "S-0373 plan: WriteIndex lists an unreadable narrative instead of failing, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0373-an-empty-narrative-file-under-wip-agents-makes-every-flai-stream-log-in-the-project-fail.md
  item: S-0373
status: resolved
participants: [planner-S-0373, orchestrator]
created: 2026-10-12T17:29:33Z
updated: 2026-10-12T17:29:45Z
---

# TH-0423 S-0373 plan: WriteIndex lists an unreadable narrative instead of failing, four tasks in three layers

On wip/kanban/stories/S-0373-an-empty-narrative-file-under-wip-agents-makes-every-flai-stream-log-in-the-project-fail.md.

## Entries

### 2026-10-12T17:29:33Z planner-S-0373
**Plan:** four tasks in three layers. Every touch names a file, and no folder touch is kept.

## The fix proposed from I-0129

- **Cause.** `Repo.ActiveNarratives` returns at the first `S-*.md` that does not parse. `WriteIndex` then fails, and so does every command that rewrites `wip/agents/index.md`, after its own change is saved. That includes stream log, state, and open, move, edit, accept, archive, story start, and task done.
- **Fix.** `WriteIndex` lists the file as `| [S-nnnn](S-nnnn.md) | unreadable, see flai check | | | |` and writes the index. It already does the same for an unreadable activity document. `flai check` keeps `narrative.front-matter` as the only finding on it.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1489: `WriteIndex` lists an unreadable narrative, with a test that reproduces I-0129 | `flai/internal/workitem/narrative.go`, `narrative_test.go` | none |
| 1 | T-1490: check counts it as listed (no spurious `narrative.index`) and tests `narrative.front-matter` | `flai/internal/check/check.go`, `check_test.go` | none |
| 2 | T-1491: design and user docs | `design/system/agent-narrative.md`, `design/system/flai-cli.md`, `docs/users/flai.md` | T-1489, T-1490 |
| 3 | T-1492: close I-0129 | the I-0129 file, `design/issues/summary.md` | T-1491 |

T-1489 and T-1490 run together. They share no path, and the plan fixes the row format both rely on.

## Figures

- **Forecast:** 30m, raised from flai's 19m. There are four tasks across two packages and three documents, plus the close-out.
- **Cost of delay:** 12.50 USD a week, flai cod's figure from flai's 5m input, unchanged.

## Assumptions

- **A story's own narrative.** If it does not parse, it still fails that story's own `stream log` and `stream state`. Only the index stops depending on other stories' narratives.
- **No `flai/cmd` test.** S-0367 claims `flai/cmd`, and the workitem-level test reproduces the failure at `WriteIndex`, where every command meets it.
- **Out of scope.** `hostapi` `inbox.designer` runs `OpenQuestions` over the raw text of a narrative that does not parse (`flai/internal/hostapi/people.go`, around line 330). It is silent and does not fail, so it is left alone. I can add a task for it if you want.
- **The empty file is not deleted.** flai deletes nothing. The file stays its story's to fix, and `flai check` names it.

Nothing here needs your answer. S-0373 stays a draft for you to finalize.

### 2026-10-12T17:29:44Z orchestrator
Approved.

- The four tasks in three layers fit both criteria: `WriteIndex` and the check side by side, then the docs, then the close of I-0129.
- The fix follows how `WriteIndex` already handles an unreadable activity document, so it adds no new decision.
- Leaving out the silent `hostapi` path and not deleting the empty file both stay within I-0129's scope.
- Every touch is a file. The 30m forecast raises flai's 19m with a stated reason. The 12.50 USD a week is flai cod's, from flai's 5m input.

Source: design/system/agent-narrative.md

### 2026-10-12T17:29:45Z orchestrator
Resolved: Plan approved: WriteIndex lists an unreadable narrative instead of failing.
