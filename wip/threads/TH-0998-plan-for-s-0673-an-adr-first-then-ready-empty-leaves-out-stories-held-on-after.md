---
id: TH-0998
title: "Plan for S-0673: an ADR first, then ready-empty leaves out stories held on after and a move to ready warns of the hold, then the guide and I-0285's close"
anchor:
  path: wip/kanban/stories/S-0673-after-chains-on-planned-stories-hold-ready-stories-for-days-while-in-progress.md
  item: S-0673
status: open
participants: [planner-S-0673]
created: 2026-11-10T08:50:10Z
updated: 2026-11-10T08:50:10Z
---

# TH-0998 Plan for S-0673: an ADR first, then ready-empty leaves out stories held on after and a move to ready warns of the hold, then the guide and I-0285's close

On wip/kanban/stories/S-0673-after-chains-on-planned-stories-hold-ready-stories-for-days-while-in-progress.md.

## Entries

### 2026-11-10T08:50:10Z planner-S-0673
S-0673 is planned: four tasks in three layers, touches, a 1h forecast, and a cost of delay value of 1,039 USD a week.

## What the instances show

- The long `after` holds were on stories planned before ADR-0230 (2026-10-31): E-0024's S-0420 to S-0425, created 2026-10-16, and S-0554. Since 2026-11-03 the longest `after` hold is S-0579's 6,259 s, and W46 has 0 s. The planner's rule already holds.
- S-0424, S-0425, and S-0554 were moved to ready by hand while a story they name in `after` was not done. `flai promote --candidates` leaves such stories out, and `flai move` to ready gives no warning.
- `readyDrained` in `flai/internal/serve/readyempty.go` counts held stories as ready. When ready held as many held stories as the in-progress room, no `ready-empty` fired, so the orchestrator was never started to promote. In-progress stayed at 1 to 4 of 5.

## Proposed remedy

- `ready-empty` counts only ready stories not held on `after`. A story held on overlap alone still counts.
- A move to ready warns of an open `after`. It warns rather than refuses, because the move is the operator's.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2928 the ADR, refining ADR-0211 §1 and ADR-0230 | `design/adrs/README.md`; the new ADR file is covered by the story's `design/adrs/` folder touch |
| 2 | T-2930 `ready-empty` count, with a test that reproduces I-0285 | `flai/internal/serve/readyempty.go`, `readyempty_test.go`, `agents_launch.go`, `orchestrate_reasons.go`, `orchestrate_reasons_test.go`, `flai/internal/inbox/events.go`, `flai/internal/mcpserver/server.go`, three design docs |
| 2 | T-2931 move-to-ready warning, with a test | `flai/internal/workitem/rules.go`, `rules_test.go`, `flai/cmd/move.go`, `docs/users/flai-reference.md`, `design/system/workflow.md` |
| 3 | T-2933 `docs/users/flai.md` and `flai issue close I-0285` | `docs/users/flai.md`, the I-0285 file, `design/issues/summary.md` |

T-2930 and T-2931 share no path, so they run together.

## Assumptions

- The ready WIP limit's count stays as it is. Changing it would change `metrics.md`, which needs its own ADR.
- The cost of delay value is lowered from `flai cod`'s 8,400 to 1,039 USD a week. That is the `after` hold rate since 2026-11-03, 6.93 h, at 150 USD an hour. The 56h input is unchanged.
- The story names no `after`.

## For the operator

Finalize the draft, or answer here to change the remedy before the story's agent writes the ADR.
