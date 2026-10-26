---
id: S-0501
type: story
nature: remediation
title: search.test.ts 'finds a body phrase across wip and design' times out at 5 s on the monorepo while other stories' close-outs run
status: review
owner: alex
created: 2026-10-26T04:01:35Z
updated: 2026-10-26T04:39:32Z
transitions:
  - to: ready
    at: 2026-10-26T04:32:27Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-26T04:32:33Z
    by: agent-S-0501
  - to: review
    at: 2026-10-26T04:39:32Z
    by: agent-S-0501
tags: [dashboard]
topics: [dashboard, testing]
touches: [flaiover/src/lib/server/search.test.ts, design/issues/I-0185-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md, design/issues/summary.md, design/issues/I-0208-search-test-ts-lists-adrs-with-the-supersession-chain-fails-on-main-since-adr-0202-s-status-became-rejected-a-status-outside-the-four-the-test-and-the-adr-template-name.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 434
  priced_by: harness
  turns:
    - day: 2026-10-26
      ceremony: 2
      hand_edits: 2
      work: 33
  models:
    - model: claude-opus-5-5
      input: 76
      output: 13503
      cache_read: 4868334
      cache_write: 144388
      cost: 2.3991
  strategic:
    - kind: orchestrator
      seconds: 725
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 118
          output: 1875
          cache_read: 25179317
          cache_write: 35917
          cost: 5.5282
        - model: claude-sonnet-5-5
          input: 6
          output: 31
          cache_read: 45577
          cache_write: 36352
          cost: 0.056
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-26T04:01:35Z
  value: 12.5
  by: planner-S-0501
  at: 2026-10-26T04:31:10Z
forecast:
  duration: 25m
  delivery: 2026-10-26T18:41:00Z
  basis: "flai forecast gave 12m (median 144 s per unit over 11 done remediation stories on claude-opus-5-5 in the small band, size 5: 2 criteria, 3 touches) and delivery 2026-10-26T18:28Z at the 22nd place in the pull order; raised to 25m for the timed runs idle and under load, the repeated vitest runs, and the close-out's flaiover tier, which dominates a story this small; delivery moved by the same 13m."
  by: planner-S-0501
  at: 2026-10-26T04:31:10Z
finalized:
  by: orchestrator
  at: 2026-10-26T04:32:16Z
---
# S-0501 search.test.ts 'finds a body phrase across wip and design' times out at 5 s on the monorepo while other stories' close-outs run

## Goal

This story remediates [I-0185](../../../design/issues/I-0185-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md), "search.test.ts 'finds a body phrase across wip and design' times out at 5 s on the monorepo while other stories' close-outs run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0185 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0185 is closed with `flai issue close I-0185 --reason` saying what fixed it

## Tasks
- T-2139 Measure the one-shot search.query on the monorepo and give the monorepo search tests a timeout that covers a loaded host
- T-2140 Close I-0185 with the reason, note why no reproducing test fits, and tick the criteria

## Notes

Cost of delay inputs set by flai from I-0185. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T12:58:18Z, 4.6 days before this story; under one cycle counts as one).

### No reproducing test

Criterion 1 is met by removing the cause, with no test that reproduces it. The timeout comes from host load: 24 one-shot `search.query` calls at once took 4.8 to 6.1 s each on a 24-core host, while one took about 2 s. A unit test cannot set the host's load, and a test that started dozens of flai processes to make its own load would be slow and still depend on the core count. The 30 s timeout in `flaiover/src/lib/server/search.test.ts` holds the measurement in its comment.

### Planning

The proposed solution: the test's `flaiAsk` (`flaiover/src/lib/server/testing.ts`) runs one `flai hostapi search.query` process per call, and that one-shot process builds the search index from nothing every time: it walks and parses every markdown file under `design/`, `wip/` with its archive, and `docs/`, 3,850 files and about 23 MB today, and the cache in `flai/internal/hostapi/searching.go` (`indexes[root]`) lives in the process, so it helps `flai serve`, which pays the build once, and never the test. The fixture tests in the same file are cheap because the fixture is small; the monorepo tests pay the whole repository, with vitest's default 5 s, and under four other stories' close-outs the build crossed it. `writes.test.ts` met the same thing on a loaded host and set 30 s with the reason as a comment (S-0116). So the story does the same: it gives the `adrs on the monorepo` block a timeout that covers a loaded host, 30 s, with a comment saying why, after measuring the one-shot `search.query` on the monorepo, idle and under parallel load, to confirm the index build is where the time goes. A test that reproduces a timeout under host load does not fit a unit test, so the story records that in its notes instead of inventing one. If the measurement shows the one-shot build itself takes seconds on an idle host, the story's agent says so in the narrative and asks on a thread whether the index should be made cheaper, rather than widening the story.

Touches: `flaiover/src/lib/server/search.test.ts` (design: the test I-0185 names, where the timeout and its comment go), `design/issues/I-0185-...md` and `design/issues/summary.md` (criteria: `flai issue close` writes both). `flai touches suggest`, seeded with `search.test.ts` and `search.ts` since the story declared none, listed `design/system/flaiover-dashboard.md`, `flaiover/package.json`, `flaiover/pnpm-lock.yaml`, `docs/users/flaiover.md`, `flaiover/src/lib/server/repo.ts`, and the search page and route (2 to 4 of 8 commits each): left out, since they co-changed when search was built or moved to flai, and a test timeout changes none of them. `flai/internal/hostapi/searching.go` is left out too: the story changes no index code unless the thread above says so. No folder touch is kept. The timing itself is unmeasured in this plan, since `flai guard` refuses the planner `flai hostapi`; the first task measures it.

Forecast: see the forecast's basis; `flai forecast` is run again after the touches are set and the figure adjusted there.

Cost of delay: 12.50 USD a week from the operator's input, 5m lost per 168h cycle at 150 USD an hour, 1.00 cycles a week. Stands: one occurrence so far, and the cost is the close-out it stops each time it recurs.
