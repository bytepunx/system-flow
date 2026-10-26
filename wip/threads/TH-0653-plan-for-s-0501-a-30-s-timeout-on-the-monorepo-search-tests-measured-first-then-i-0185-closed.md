---
id: TH-0653
title: "Plan for S-0501: a 30 s timeout on the monorepo search tests, measured first, then I-0185 closed"
anchor:
  path: wip/kanban/stories/S-0501-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md
  item: S-0501
status: resolved
participants: [planner-S-0501, orchestrator, agent-S-0501]
created: 2026-10-26T04:32:05Z
updated: 2026-10-26T04:40:11Z
---

# TH-0653 Plan for S-0501: a 30 s timeout on the monorepo search tests, measured first, then I-0185 closed

On wip/kanban/stories/S-0501-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md.

## Entries

### 2026-10-26T04:32:05Z planner-S-0501
Plan for S-0501, two tasks in two layers.

Layer 1: T-2139 measures `bin/flai hostapi search.query` and `adrs.list` on the monorepo, idle and under load, and gives the `adrs on the monorepo` block of `flaiover/src/lib/server/search.test.ts` a 30 s timeout with a comment naming I-0185 and the reason.

Layer 2: T-2140, after T-2139, closes I-0185 with `flai issue close --reason` naming the cause, the figures, and the timeout, notes in the story why no reproducing test fits, and ticks the criteria.

The cause, as read from the code: the test's `flaiAsk` runs a one-shot `flai hostapi` process per call, which builds the whole repository's search index from nothing (3,850 markdown files, about 23 MB, the archive included) with no cache across processes; `flai serve` keeps that index, the test never does. Under four other stories' close-outs the build crossed vitest's 5 s default. `writes.test.ts` set 30 s for the same reason in S-0116, and this plan follows it.

Assumptions: the timing is unmeasured here, since `flai guard` refuses the planner `flai hostapi`, so T-2139 measures before it changes the test; a unit test cannot reproduce a timeout under host load, so the first criterion is met by removing the cause and the story's notes say so; `searching.go` and the index stay as they are, and if the idle build itself takes seconds the story's agent asks on a thread rather than widening the story. Touches are the test file and the two issue files; the co-change suggestions (dashboard design, package files, user docs, repo.ts, the search page) are left out as unrelated to a timeout. Forecast 25m over flai's 12m for the timed and repeated runs; cost of delay 12.50 USD a week from the operator's input, unchanged.

### 2026-10-26T04:32:11Z orchestrator
Approved. The plan fits S-0501.

- **The remedy is grounded.** Each one-shot `flai hostapi` builds the whole search index, and only `flai serve` keeps it, so the monorepo tests pay the full build under vitest's 5 s default. T-2139 measures that first, then sets 30 s with a reason, as `writes.test.ts` did in S-0116. If the idle build itself is slow, the agent asks rather than widening the story.
- **The tasks fit.** There are two tasks in two layers and three file touches, with the co-change files left out for stated reasons.
- **The figures fit.** The forecast is 25m, raised from flai's 12m for the timed runs. The value is 12.50 USD/week from `flai cod` on the issue's 5m input. I set no input.
- **The first criterion** is met by removing the cause, with the reason no reproducing unit test fits written in the story's notes.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-26T04:32:16Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, forecast, and value fit S-0501.

### 2026-10-26T04:33:01Z planner-S-0501
The story went to ready and in progress while the plan's last edit was in flight, so one line of the story's `### Planning` notes still reads "Forecast: see the forecast's basis" instead of the figures. The forecast block itself carries them: flai gave 12m at size 5 and delivery 2026-10-26T18:28Z; set to 25m, delivery 18:41Z, for the timed and repeated runs and the close-out's flaiover tier. The planner leaves the story to its agent now rather than editing it under them.

### 2026-10-26T04:35:51Z agent-S-0501
Noted; the forecast block is the record, and the notes line is left as it is. T-2139 is done: one-shot search.query measured at about 2 s of CPU on the monorepo and 4.8 to 6.1 s with 24 at once, and the block now has a 30 s timeout. The idle cost is asked on TH-0654.

### 2026-10-26T04:40:11Z orchestrator
Resolved: Plan approved and delivered: the later entries are notes. The forecast block is the record, and T-2139 and T-2140 are done. Nothing on this thread awaits an answer.
