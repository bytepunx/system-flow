---
id: TH-0474
title: "S-0427 needs a cost of delay input: recommend time lost 2h per cycle"
anchor:
  path: wip/kanban/stories/S-0427-experiment-with-local-storage-solutions-that-synchronize-with-the-documentation-in-the-file-system.md
  item: S-0427
status: open
participants: [planner-S-0427, orchestrator]
created: 2026-10-16T05:55:52Z
updated: 2026-10-16T05:56:06Z
---

# TH-0474 S-0427 needs a cost of delay input: recommend time lost 2h per cycle

On wip/kanban/stories/S-0427-experiment-with-local-storage-solutions-that-synchronize-with-the-documentation-in-the-file-system.md.

## Entries

### 2026-10-16T05:55:52Z planner-S-0427
S-0427 has no cost of delay inputs and no epic to take a share from, so `flai cod S-0427` cannot work out a value.

Recommendation: `time_lost_per_cycle: 2h`. The story is an experiment. Its value is the operator wait and agent overhead that file reads cost on each board, chart, and prime call, and a fast answer lets a follow-up story act on it sooner. It brings no revenue, and waiting on it carries no penalty.

Alternatives:

- `time_lost_per_cycle: 1h`, if the board and charts feel fast enough today and this is exploratory.
- `penalty_per_week: 250` (USD), if the read overhead already costs paid agent time you want priced.

Set it with `flai edit S-0427 --time-lost-per-cycle 2h` (or answer here). I am planning the touches, forecast, and tasks meanwhile, and will work out the value once the input is set.

### 2026-10-16T05:56:06Z orchestrator (recommendation)
**Recommendation:** set `time_lost_per_cycle: 2h` on S-0427, as the planner recommends: `flai edit S-0427 --time-lost-per-cycle 2h`.

## Why 2h

- **The goal asks for it.** It measures the "operator wait time and agent overhead" that file reads cost on each board, chart, and prime call. That is time lost, not revenue or a penalty.
- **Nothing in the story points to an alternative.** The story does not say the board feels fast enough, which would argue for 1h. It does not price agent time as a penalty either, which would argue for 250 USD.

## Why this needs you

I tried to set it, and flai refused: "S-0427 is ready, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". You moved S-0427 to ready, so its input is yours to set.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
