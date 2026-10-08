---
id: TH-0384
title: "S-0348 plan: a check scoped to a story leaves board.wip-limit out, then close I-0123"
anchor:
  path: wip/kanban/stories/S-0348-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md
  item: S-0348
status: open
participants: [planner-S-0348]
created: 2026-10-08T09:00:05Z
updated: 2026-10-08T09:00:05Z
---

# TH-0384 S-0348 plan: a check scoped to a story leaves board.wip-limit out, then close I-0123

On wip/kanban/stories/S-0348-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md.

## Entries

### 2026-10-08T09:00:05Z planner-S-0348
Plan for S-0348, a draft with 4 tasks in 3 layers.

Proposed fix: I-0123's one instance is S-0339's close-out finding `wip/kanban/board.md: 6 stories in in-progress, limit 5`. That count is the main checkout's state: no story branch changes or clears it, and it names no story. So `ScopeToStory` (`flai/internal/check/scope.go`) leaves `board.wip-limit` out, as it already leaves out `item.archive` (ADR-0122) and findings on another open story's narrative (ADR-0123, ADR-0125). A close-out then neither notes nor records it. `flai check --strict` in the main checkout still fails on it, so ADR-0073 stands.

Tasks:

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-1422: an ADR recording the decision | none | `design/adrs` |
| 2 | T-1423: drop the finding in `ScopeToStory`, add tests reproducing I-0123, update `recordOutside`'s comment and the help, regenerate the reference | T-1422 | `scope.go`, `scope_test.go`, `flai/cmd/check.go`, `check_test.go`, `docs/users/flai-reference.md` |
| 2 | T-1424: `flai-cli.md`, `continuous-improvement.md`, `docs/users/flai.md`, `work-management.md` and its template copy, `template/CHANGELOG.md` | T-1422 | those six files |
| 3 | T-1425: `flai issue close I-0123 --reason` | T-1423, T-1424 | the issue file, `design/issues/summary.md` |

T-1423 and T-1424 share no path, so they can run together.

Figures:

- Forecast: 21m, delivered 2026-10-08T14:34Z, as `flai forecast` gives it. Kept: the precedents S-0323, S-0280, and S-0318 took 12m, 22m, and 49m.
- Cost of delay: 12.50 USD a week, from the `time_lost_per_cycle: 5m` set on TH-0383.

Assumptions:

1. The finding is left out, not noted without an issue as `wip.overlap` is (ADR-0115). T-1422 lets the story's agent choose the note instead if the ADR argues for it.
2. All three columns are left out of the scoped check. Review over its limit is already advisory (ADR-0073), and leaving it out too keeps the rule simple.
3. The ADR is needed, as for each precedent. So `design/adrs` stays a folder touch: the ADR's number is not known yet. While this story is in progress, that touch holds S-0334, S-0337, and S-0297.
4. `flai/internal/check/check.go`, which produces the finding, does not change.

Proposed changes to the tasks: none.
