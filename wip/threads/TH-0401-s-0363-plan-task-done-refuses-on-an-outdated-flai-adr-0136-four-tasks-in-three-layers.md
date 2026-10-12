---
id: TH-0401
title: "S-0363 plan: task_done refuses on an outdated flai (ADR-0136), four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0363-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md
  item: S-0363
status: open
participants: [planner-S-0363, orchestrator]
created: 2026-10-12T14:37:23Z
updated: 2026-10-12T14:37:49Z
---

# TH-0401 S-0363 plan: task_done refuses on an outdated flai (ADR-0136), four tasks in three layers

On wip/kanban/stories/S-0363-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md.

## Entries

### 2026-10-12T14:37:23Z planner-S-0363
S-0363 is planned. It has 4 draft tasks, 12 file touches, a 20m forecast, and a cost of delay of 12.50 USD a week. The story's `### Planning` note says where each touch and figure came from.

## Proposed remedy

`taskdone.Run` refuses, changing nothing, while `release.FlaiOutdated` says the running flai is older than the newest `flai/vX.Y.Z` tag in the main checkout's history. The refusal names the upgrade command. A dev build is never refused.

- The MCP tool, `flai task done`, and the host's `task.done` all reach `taskdone.Run`, so one check covers the three.
- It is a refusal, not a warning: I-0127's instance ran with the inbox's `flai_outdated` warning showing, and the close went ahead.

## Tasks and layers

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1439 ADR-0136, the decision | none | the ADR, `design/adrs/README.md` |
| 2 | T-1440 the check in `taskdone.Run`, tests through Run and MCP, help and reference | T-1439 | `taskdone.go`, `taskdone_test.go`, `mcpserver/task_test.go`, `cmd/task_done.go`, `flai-reference.md` |
| 2 | T-1441 docs: `flai-cli.md`, `docs/users/flai.md`, `docs/operators/index.md` | T-1439 | those three |
| 3 | T-1442 close I-0127 | T-1440, T-1441 | the issue, `summary.md` |

T-1440 and T-1441 share no path, so they run together.

## Assumptions

1. A refusal departs from ADR-0117, whose minimum only warns. ADR-0136 limits the refusal to `task_done`, because `task_done` commits and moves a task to done in one call, which cannot be undone once the task is done. Every other command keeps warning.
2. A host one patch behind (1.40.1 against 1.40.3 today) stops every story's agent at its first task close until `flai host upgrade` runs. I take that cost as intended. If you want a softer form, say which and I will re-plan: warn-only, or refuse only when the minor version differs.
3. "The project's history" is the tags reachable from the main checkout's HEAD, as `flai_outdated` already reads them.
4. Once ready, S-0363 is held behind S-0362 (same two test files) and S-0337 (`design/adrs`, `flai-cli.md`, `docs/users/flai.md`). Each task says to sync after those are accepted.

## Proposals

No split, merge, or drop.

## What you need to do

Review the remedy, especially assumption 2, then finalize the draft.

### 2026-10-12T14:37:49Z orchestrator (recommendation)
**Recommendation:** take the planner's softer form. `task_done` refuses only while the running flai's major or minor version is behind the project's history, and warns on a patch behind. Keep the rest of the plan as drafted.

Why this is yours: ADR-0136 would be a new decision. It departs from ADR-0117 § Decision 7, under which an outdated flai warns and is never refused. A recorded decision is not mine to make.

Why the softer form:

- Under the planner's form, a host one patch behind stops every story's agent at its next task close. That lasts until `flai host upgrade` runs.
- Every release puts the host a patch behind. Today the host runs 1.40.1, and I have published 1.40.3 and 1.40.4 since. So under the strict form, each publish would stop all running stories until you upgrade.
- I-0127's instance ran 1.39.13 against 1.40.0 in the history, so a minor difference. The softer form still refuses it.
- The cost: a patch behind can still lack a fix, as 1.39.13 lacked S-0322's in 1.39.14. The inbox's `flai_outdated` warning keeps covering that case.

If you take this, the planner re-plans T-1439 and T-1440 to the minor-version threshold. The four tasks and three layers stay the same. If you prefer the strict form, confirm the plan as drafted.

Either way, I finalize S-0363 once you answer. The figures are flai's (20m, 12.50 USD a week), and its 12 file touches fit the criteria.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/adrs/0117-the-host-api-installs-a-release-the-dashboard-names-only-when-flai-lists-it-as.md § Decision
