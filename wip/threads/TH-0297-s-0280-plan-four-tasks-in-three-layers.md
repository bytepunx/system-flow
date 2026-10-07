---
id: TH-0297
title: "S-0280 plan: four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0280-flai-check-finds-item-archive-outside-the-story-at-close-out.md
  item: S-0280
status: open
participants: [planner-S-0280]
created: 2026-10-07T15:05:21Z
updated: 2026-10-07T15:05:21Z
---

# TH-0297 S-0280 plan: four tasks in three layers

On wip/kanban/stories/S-0280-flai-check-finds-item-archive-outside-the-story-at-close-out.md.

## Entries

### 2026-10-07T15:05:21Z planner-S-0280
S-0280 is planned: four tasks in three layers, touches, and a 25m forecast. The cost of delay value waits on TH-0295.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1164 An ADR refining ADR-0085 and ADR-0115 records the remedy | none | `design/adrs` |
| 2 | T-1165 A close-out records no `item.archive` and leaves out one not naming the story | T-1164 | `flai/internal/check/scope.go`, `scope_test.go`, `flai/cmd/check.go`, `check_test.go`, `docs/users/flai-reference.md` |
| 3 | T-1166 I-0078 is closed with `flai issue close` | T-1165 | the I-0078 file, `design/issues/summary.md` |
| 3 | T-1167 The convention, the design, and the users' guide say so | T-1165 | both `work-management.md` copies, `template/CHANGELOG.md`, `continuous-improvement.md`, `flai-cli.md`, `docs/users/flai.md` |

T-1166 and T-1167 share no path and run together.

## Assumptions

- **The remedy.** All 38 instances name one file, S-0250's, cancelled and left unarchived until commit 1f7b7296. `item.archive` never names the closing story, so it is never the story's to fix. The plan follows ADR-0115's pattern for `wip.overlap`: a close-out records no `item.archive`, and its scoped check leaves it out. `flai check` in the main checkout still warns. T-1164's ADR may still choose otherwise. T-1165's Work names this remedy.
- **Not planned: archiving on cancel.** It would stop the condition, not just its recording. It reverses ADR-0055, which lets a cancelled item move back to backlog, and ADR-0028, which leaves a cancelled story's narrative in place. If you want it, I would make it a separate story rather than grow this one.
- **`flai/internal/check/check.go` is not touched**, since the rule itself stays. T-1165 widens its touches if it must change.
- **Forecast 25m**: flai's 21m raised by 4m for three serial layers, in line with S-0279's measured 26m.

## Proposed changes

None to split, merge, or drop.
