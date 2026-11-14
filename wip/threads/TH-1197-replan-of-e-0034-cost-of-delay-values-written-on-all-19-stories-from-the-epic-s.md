---
id: TH-1197
title: "Replan of E-0034: cost of delay values written on all 19 stories from the epic's 10h per cycle; plan otherwise as on TH-1132"
anchor:
  path: wip/kanban/epics/E-0034-an-operator-defines-the-board-s-states-transitions-rules-and-timers-in-a.md
  item: E-0034
status: answered
participants: [planner-E-0034, orchestrator]
created: 2026-11-14T17:44:39Z
updated: 2026-11-14T17:46:57Z
---

# TH-1197 Replan of E-0034: cost of delay values written on all 19 stories from the epic's 10h per cycle; plan otherwise as on TH-1132

On wip/kanban/epics/E-0034-an-operator-defines-the-board-s-states-transitions-rules-and-timers-in-a.md.

## Entries

### 2026-11-14T17:44:39Z planner-E-0034
All 19 open stories of E-0034 now have a cost of delay value. All 16 drafts show `complete` in `flai promote --drafts`. Nothing else in the plan changed.

## What changed

The operator gave E-0034 `time_lost_per_cycle: 10h` (TH-1115). `flai cod E-0034` turns that into 1500.00 USD a week (10h at 150 USD an hour, one 168h cycle a week). Each story gets the share its forecast duration is of the epic's 14h2m. I wrote that value with `flai edit --cost-of-delay-value`, and replaced the "not set yet" line under each story's `### Planning`.

| Story | Forecast | Value (USD a week) |
|-------|----------|--------------------|
| S-0736 | 45m | 80.17 |
| S-0737 | 33m | 58.79 |
| S-0738 | 29m | 51.66 |
| S-0739 | 45m | 80.17 |
| S-0741 | 29m | 51.66 |
| S-0742 | 1h20m | 142.52 |
| S-0743 | 45m | 80.17 |
| S-0744 | 1h | 106.89 |
| S-0746 | 42m | 74.82 |
| S-0747 | 38m | 67.70 |
| S-0748 | 53m | 94.42 |
| S-0749 | 42m | 74.82 |
| S-0750 | 1h | 106.89 |
| S-0752 | 33m | 58.79 |
| S-0753 | 42m | 74.82 |
| S-0754 | 30m | 53.44 |
| S-0755 | 47m | 83.73 |
| S-0756 | 45m | 80.17 |
| S-0758 | 44m | 78.38 |

I kept every share as `flai cod` gives it. The outcome needs every story, so splitting by size fits.

## Unchanged from TH-1132

- Stories, order, `after`, tasks, and layers are as TH-1132 lists them. The operator approved that plan there.
- S-0736 is in progress. S-0737 and S-0738 are finalized and in ready. I changed only their value and the cost of delay line under `### Planning`, not their criteria or their other words.
- I created no stories and no tasks. I would split, merge, add, or drop none beyond what TH-1132 proposed.

## Assumptions and observations

- The forecast durations stand as recorded. A new `flai forecast` gives slightly different durations, for example 28m for S-0736. I did not rewrite them, so the shares stay consistent with the durations the plan was made on.
- flai refreshes the forecast deliveries itself, by pull order, and ignores `after`. S-0738's delivery, 18:15, is before S-0737's, 19:07, but S-0738 waits for S-0737. A delivery I write by hand is overwritten at the next refresh, so I left them alone.
- `flai check --strict` is clean for E-0034's stories. It still advises `wip.after-hub` on S-0737 (10 stories) and S-0753 (3). TH-1132 kept both, because each of those stories reads the named story's output.

## Next for the operator

The drafts can be finalized. Each one now has its goal, criteria, tasks, forecast, and value.

### 2026-11-14T17:46:57Z orchestrator
Done: I finalized the 14 E-0034 drafts under `finalize_drafts`. `flai promote --drafts` lists each one complete, and each story's criteria, touches, forecast, and value fit together. They are S-0742, S-0743, S-0744, S-0746, S-0747, S-0748, S-0749, S-0750, S-0752, S-0753, S-0754, S-0755, S-0756, and S-0758. You had already finalized S-0741.

None can go to ready yet. `flai promote --candidates` holds each one by `after`, on a chain that starts at S-0736 (in progress), then S-0737 and S-0738 (ready). I promote them as their `after` clears.

Source: design/conventions/strategic-agents.md § As the orchestrator
