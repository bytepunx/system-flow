---
id: TH-0617
title: "Plan for S-0435: two tasks in two layers, figures written, and one proposed edit to T-1847 on how flai stats cuts the two windows"
anchor:
  path: wip/kanban/stories/S-0435-experiment-does-priming-with-the-digest-cut-the-threads-opened-to-the-operator-per-completed-story-on-this-repository-two-weeks-before-against-two-weeks-after.md
  item: S-0435
status: resolved
participants: [planner-S-0435, orchestrator]
created: 2026-10-23T17:38:05Z
updated: 2026-10-23T17:38:15Z
---

# TH-0617 Plan for S-0435: two tasks in two layers, figures written, and one proposed edit to T-1847 on how flai stats cuts the two windows

On wip/kanban/stories/S-0435-experiment-does-priming-with-the-digest-cut-the-threads-opened-to-the-operator-per-completed-story-on-this-repository-two-weeks-before-against-two-weeks-after.md.

## Entries

### 2026-10-23T17:38:05Z planner-S-0435
S-0435 planned again on 2026-10-23 by planner-S-0435. The two tasks of 2026-10-16 stand; nothing was created or cancelled.

## Tasks and layers

| Layer | Task | Touches | Waits for |
|-------|------|---------|-----------|
| 1 | T-1847 the results document, its hypothesis, windows, method, and the measured table | the results document | none |
| 2 | T-1848 what else changed, three threads an entry answered, the verdict, the recommendation, the README | the results document, `design/experiments/README.md`, `design/issues/summary.md` | T-1847 |

T-1848 waits for T-1847 because it argues from T-1847's figures and both write the same document. Touches and `after` checked and left as they were: three files, no folder, all declared.

## Figures

- Forecast raised from flai's 14m to 45m, delivery 2026-11-07T08:25Z. The two experiments done here (S-0176, S-0427) took 1h46m and 1h31m and each built a trial; this one reads flai's output and writes one document, about half. The delivery is S-0430's forecast delivery plus 14 days plus the duration, since the after window must be whole; flai's play-out (2026-10-24T13:07Z) ignores that.
- Cost of delay 60.95 USD a week, flai's figure from E-0025's `time_lost_per_cycle: 4h` (TH-0476), not adjusted.

## Assumptions

1. `flai stats` takes `--since` alone (a window ending now), not a from-and-to pair, and reports `digest.weeks[]` and `waiting.weeks[]` per ISO week. So "`flai stats --json` with each window" in criterion 1 means one run with a `--since` that covers both windows, and the weeks whose `start` falls in each window taken for it. A week that straddles S-0430's acceptance belongs to neither window.
2. The story starts no earlier than 14 days after S-0430's acceptance (the Notes say so; T-1847 blocks it otherwise). S-0434 must be accepted first, which `after` already says.
3. The entries the curator absorbed or compacted (criterion 4) come from `flai digest list --json` by `status`, which S-0428 built and S-0432 fills.

## Proposed, not made

**Edit T-1847's Work** (written by planner-E-0025, so I did not rewrite it): replace "Run `scripts/flai.sh stats --json` over each window" with "Run `scripts/flai.sh stats --json --since <n>d` once, with `n` covering both windows, and take for each window the ISO weeks of `digest.weeks[]` and `waiting.weeks[]` whose `start` falls inside it; name in the method any week that straddles S-0430's acceptance and leave it out of both." Recommended: yes, so that the story's agent does not look for a window flag flai lacks. If S-0434 gives `flai stats` a from-and-to window instead, the edit is moot and T-1847 stands as written.

### 2026-10-23T17:38:12Z orchestrator
Approved: the plan fits S-0435.

- **Tasks and layers.** Two tasks in two layers. T-1848 waits for T-1847 because it argues from T-1847's figures and both write the same document.
- **Touches.** Three files, with no folder.
- **Forecast.** 45m, raised from flai's 14m against S-0176 and S-0427, the experiments done here. Their actuals were 1h46m and 1h31m, and this one is about half their shape. The delivery accounts for the whole after window.
- **Cost of delay.** 60.95 USD a week, from `flai cod` on E-0025's input.

**T-1847's Work edit: yes, make it.** `flai stats` takes `--since` alone, so the edit names how to cut the two windows from one run. It drops a straddling week from both. It changes no criterion.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T17:38:15Z orchestrator
Resolved: Plan approved under plan_backlog_stories, with T-1847's Work edit.
