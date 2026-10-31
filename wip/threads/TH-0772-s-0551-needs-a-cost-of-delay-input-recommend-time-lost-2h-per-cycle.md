---
id: TH-0772
title: "S-0551 needs a cost of delay input: recommend time lost 2h per cycle"
anchor:
  path: wip/kanban/stories/S-0551-a-story-starts-on-a-trial-merge-instead-of-an-overlap-hold-and-after-holds-a.md
  item: S-0551
status: open
participants: [planner-S-0551, orchestrator]
created: 2026-10-31T07:35:33Z
updated: 2026-10-31T07:37:02Z
---

# TH-0772 S-0551 needs a cost of delay input: recommend time lost 2h per cycle

On wip/kanban/stories/S-0551-a-story-starts-on-a-trial-merge-instead-of-an-overlap-hold-and-after-holds-a.md.

## Entries

### 2026-10-31T07:35:33Z planner-S-0551
**Recommendation:** set `time_lost_per_cycle: 2h` on S-0551: `flai edit S-0551 --time-lost-per-cycle 2h`.

S-0551 has no cost of delay inputs, and E-0031 has neither inputs nor a value, so `flai cod S-0551` cannot work out a value.

## Why 2h

- **The goal is time lost.** The ready wait at p85 was 5.4 hours in the fortnight, and the journal counts 124 overlap holds and 113 `after` holds against 62 real conflicts. Most of that wait held a story that would have merged cleanly.
- **Not all of it is lost.** Some holds were real conflicts, and the pull order keeps some wait anyway. Roughly a third to half of the p85 wait is the hold this story removes, so about 2h.

## Alternatives

- `1h`, if you judge most of the wait to be the in-progress limit rather than the holds.
- `4h`, if you count the idle time of the host too: nothing was in progress for 42 percent of the fortnight.

I am drafting the tasks, the touches, and the forecast meanwhile, and work out the value once you answer.

### 2026-10-31T07:37:02Z orchestrator (recommendation)
**Recommendation:** set `time_lost_per_cycle: 2h` on S-0551, as the planner recommends:

```bash
flai edit S-0551 --time-lost-per-cycle 2h
```

## The figure, and not 1h or 4h

- **The goal is time lost.** The story removes the hold on stories that would have merged cleanly. That is waiting time, not revenue or a penalty.
- **The story points to no alternative.** It does not say the in-progress limit causes most of the wait, which would argue for 1h. It does not count the host's idle time as lost, which would argue for 4h.
- **It matches D-0316.** You set 2h on S-0427, also a flow-overhead story.

## Why this needs you

I tried to set it, and flai refused: "S-0551 is ready, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". You moved S-0551 to ready at 2026-10-31T07:24:00Z, so its input is yours to set. A confirmation of this entry does not write the input. Running the command does.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
