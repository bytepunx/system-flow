---
id: TH-0988
title: "Plan for S-0662: one evidence reader, then the move, the close-out, item_get, the review page, and acceptance; nine tasks in four layers, 2h30m"
anchor:
  path: wip/kanban/stories/S-0662-the-move-to-review-the-close-out-and-the-review-page-say-what-was-deployed-and.md
  item: S-0662
status: resolved
participants: [planner-S-0662, orchestrator]
created: 2026-11-10T08:29:42Z
updated: 2026-11-10T08:32:25Z
---

# TH-0988 Plan for S-0662: one evidence reader, then the move, the close-out, item_get, the review page, and acceptance; nine tasks in four layers, 2h30m

On wip/kanban/stories/S-0662-the-move-to-review-the-close-out-and-the-review-page-say-what-was-deployed-and.md.

## Entries

### 2026-11-10T08:29:42Z planner-S-0662
**Plan.** Nine tasks in four layers. Forecast 2h30m, raised from `flai forecast`'s 1h27m. The cost of delay waits on TH-0983 and TH-0979.

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-2875 `kluster.Evidence`: a story namespace's last deploy and checks run, and whether both passed at the head | `flai/internal/kluster/evidence.go`, `evidence_test.go` | none |
| 1 | T-2877 work-management convention, its template copy, and `workflow.md`: deploy, check, and tick before review | the three, `template/CHANGELOG.md` | none |
| 1 | T-2878 the prime pack names the deploy and checks commands | `flai/internal/context/story.go`, `context.go`, `context_test.go` | none |
| 2 | T-2880 `flai move` and `item_move` warnings, criteria 1 and 2 | `flai/cmd/move.go`, `move_test.go`, `mcpserver/server.go`, `review_test.go` | T-2875 |
| 2 | T-2881 the close-out records the deploy and checks; `flai verify --last` prints them | `verify/story.go`, `verify/kluster.go`, `kluster_test.go`, `cmd/verify.go`, `verify_test.go` | T-2875 |
| 2 | T-2883 the orchestrator's evidence may name the record | `preview/orchestrator.go`, `evidence_test.go`, `cmd/accept.go`, `accept_orchestrator_test.go`, the strategic-agents design and convention, the role file's three copies, `template/CHANGELOG.md` | T-2875, T-2877 |
| 3 | T-2885 `item_get` carries a `kluster` field | `mcpserver/server.go`, `server_test.go` | T-2880, T-2881 |
| 3 | T-2886 the review page shows it beside the criteria, with log links | `hostapi/writes_host.go`, `writes_test.go`, `review.ts`, `review.test.ts`, `Review.svelte`, `Review.svelte.test.ts` | T-2881 |
| 4 | T-2888 `kluster.md`, `flai-cli.md`, `flaiover-dashboard.md`, `docs/users/flai.md`, `docs/users/flaiover.md` | those five | all code tasks |

Why each waits:

- T-2883 waits for T-2877 because both write `template/CHANGELOG.md`.
- T-2885 waits for T-2880 because both change `mcpserver/server.go`.

**Assumptions.**

- S-0659's ADR sets what the record holds and the evidence form. The tasks name a shape, such as `- <n>: kluster checks passed at <commit>`, only as an example.
- S-0661 writes a checks run record with its commit; today `ChecksRun` in `flai/internal/serve/checks.go` has none. T-2875 reads it by path, since `serve` imports `kluster`, not the reverse.
- The warning is added by the callers of `Repo.Move`. `workitem` cannot see kluster, so `rules.go` is not touched.
- Criterion 5 needs code, not only text. `unevidenced` in `flai/internal/preview/orchestrator.go` refuses a criterion with no changed file today.
- A container story verifies inside its container, where `.flai-cache/kluster` is not mounted (ADR-0274). T-2881 leaves the choice to the story's agent: read where the close-out can reach, or fill the part on the host.
- "In `design/system`" in criterion 6 means `design/system/workflow.md` § Definition of done.
- The role file changes in three copies: `design/workflows/roles/orchestrator.md`, its template copy, and the embedded `flai/internal/roles/baseline/orchestrator.md`.

**Proposed, not created.** Nothing to split, merge, add, or drop.

### 2026-11-10T08:32:12Z orchestrator
**Plan approved.** Nine tasks in four layers fit S-0662's eight criteria:

| Criteria | Task |
|----------|------|
| 1, 2 | T-2880, on T-2875's evidence reader |
| 3 | T-2881 (close-out, `flai verify --last`) and T-2885 (`item_get`) |
| 4 | T-2886 |
| 5 | T-2883 |
| 6 | T-2877 |
| 7 | T-2878 |
| 8 | T-2888 |

- The touches are files, no folder. The two shared paths, `template/CHANGELOG.md` and `flai/internal/mcpserver/server.go`, are ordered by `after` (T-2883 on T-2877, T-2885 on T-2880).
- The forecast of 2h30m is flai's 1h27m raised with a stated reason.
- `after: [S-0661]` alone is right: S-0661 names S-0659 and S-0660 already.
- The cost of delay waits on TH-0983 and TH-0979, which have recommendations awaiting the operator. Nothing for me to set here.

The assumptions about S-0659's ADR and S-0661's record shapes stand: the story's agent follows those stories where they decide otherwise. Resolving this thread.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-11-10T08:32:25Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and forecast fit S-0662's criteria; cost of delay waits on TH-0983 and TH-0979.
