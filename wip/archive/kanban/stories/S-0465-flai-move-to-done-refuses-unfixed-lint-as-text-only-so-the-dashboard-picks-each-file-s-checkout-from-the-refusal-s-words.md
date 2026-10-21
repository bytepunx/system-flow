---
id: S-0465
type: story
nature: improvement
title: flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words
status: done
owner: alex
created: 2026-10-21T07:24:59Z
updated: 2026-10-21T09:04:46Z
transitions:
  - to: ready
    at: 2026-10-21T07:38:29Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-21T07:38:37Z
    by: agent-S-0465
  - to: review
    at: 2026-10-21T08:50:31Z
    by: agent-S-0465
  - to: done
    at: 2026-10-21T09:04:46Z
    by: orchestrator
tags: [flai, flaiover]
topics: [cli, dashboard]
touches: [flai/cmd/move.go, flai/cmd/move_refused_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, docs/users/flai.md, design/system/flai-cli.md, flaiover/src/lib/components/AcceptFindings.svelte, flaiover/src/lib/components/AcceptFindings.svelte.test.ts, design/components/flaiover/src/lib/components/AcceptFindings.examples.json, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", "flaiover/src/routes/api/items/[id]/move/+server.ts", flaiover/src/lib/filelinks.ts, flaiover/src/lib/filelinks.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/issues/I-0145-flai-move-to-done-refuses-unfixed-lint-as-text-only-so-the-dashboard-picks-each-file-s-checkout-from-the-refusal-s-words.md, design/issues/summary.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4331
  priced_by: harness
  turns:
    - day: 2026-10-21
      ceremony: 4
      hand_edits: 2
      work: 45
  models:
    - model: claude-opus-5-5
      input: 330
      output: 80749
      cache_read: 16358812
      cache_write: 673693
      cost: 8.9191
  strategic:
    - kind: orchestrator
      seconds: 483
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 42
          output: 575
          cache_read: 4155665
          cache_write: 23265
          cost: 1.0477
        - model: claude-sonnet-5-5
          input: 10
          output: 52
          cache_read: 115647
          cache_write: 44110
          cost: 0.1308
cost_of_delay:
  inputs:
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-10-21T07:24:59Z
  value: 22.5
  by: planner-S-0465
  at: 2026-10-21T07:37:38Z
forecast:
  duration: 1h
  delivery: 2026-10-22T01:57:00Z
  basis: "flai forecast's 36m (89 s per unit over 81 large-band improvement stories, size 24) raised to 1h: the change spans flai and flaiover, extracts a component with its examples, and runs both test tiers; S-0378, which built this refusal for the review page, took 2h."
  by: planner-S-0465
  at: 2026-10-21T07:37:38Z
finalized:
  by: orchestrator
  at: 2026-10-21T07:38:19Z
---
# S-0465 flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words

## Goal

This story remediates [I-0145](../../../design/issues/I-0145-flai-move-to-done-refuses-unfixed-lint-as-text-only-so-the-dashboard-picks-each-file-s-checkout-from-the-refusal-s-words.md), "flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0145 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0145 is closed with `flai issue close I-0145 --reason` saying what fixed it

## Tasks
- T-1979 flai move to done prints a refused acceptance as data with exit 4 under --json, and item.move answers it as Refused
- T-1980 AcceptFindings lists a refused acceptance's findings, each opened on its own checkout, and Review uses it
- T-1981 The board and the item page list a refused move to done's findings from its data, and filelinks no longer reads the refusal's words
- T-1982 Close I-0145 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0145. time_lost_per_cycle 9m: 10m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-10-13T18:47:09Z, 7.5 days before this story).

### Planning

The fix is the remedy the instance of I-0145 names. `flai move <story> done --json` refuses through `acceptRefusal`, as `flai accept --json` does. `item.move` maps exit 4 to `Refused`. The board and the item page then read the findings as data, as the review page does. `linkSegments` stops matching flai's wording.

Touches, by source:

| Touch | Source |
|-------|--------|
| `flai/cmd/move.go` | design: the issue names it; co-change |
| `flai/cmd/move_refused_test.go` | layout: a new test beside `accept_fix_test.go` |
| `flai/internal/hostapi/writes.go`, `writes_test.go` | design: the issue names `item.move`'s exits |
| `docs/users/flai.md`, `design/system/flai-cli.md` | co-change; documentation rules |
| `flaiover/src/lib/components/AcceptFindings.svelte`, its test, `design/components/flaiover/src/lib/components/AcceptFindings.examples.json` | layout: the findings list moved out of `Review.svelte`; ADR-0164 for the examples |
| `flaiover/src/lib/components/Review.svelte`, `Review.svelte.test.ts` | layout: Review holds the list today |
| `flaiover/src/routes/board/+page.svelte`, `board.svelte.test.ts` | design: the issue names the board; co-change |
| `flaiover/src/routes/items/[id]/+page.svelte`, `item.svelte.test.ts` | design: the issue names the item page; co-change |
| `flaiover/src/routes/api/items/[id]/move/+server.ts` | co-change: its comment documents the answer |
| `flaiover/src/lib/filelinks.ts`, `filelinks.test.ts` | design: the issue names `linkSegments` |
| `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md` | co-change; documentation rules |
| `design/issues/I-0145-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes them |

The story declared no touches. No folder touch was kept: every file is named.

Forecast: flai gave 36m (size 24, 89 s a unit). Raised to 1h. The change spans two components, moves a component out with its examples, and runs the Go and flaiover tiers. S-0378, which built this refusal for the review page, took 7163 s. Delivery is flai's, moved by the extra 24m.

Cost of delay: 22.50 USD a week, as `flai cod` gives it from flai's input of 9m lost a cycle at 150 USD an hour. It stands. The issue has one instance, and the only loss is a file opened on the wrong checkout after a wording change.

### Accepted by the orchestrator

- Verified: b41df5464e287a6e010a7e03251c9de72e9c731a
- Batch: 20261021-085133
- Enqueued verified: b41df5464e287a6e010a7e03251c9de72e9c731a
- At: 2026-10-21T09:04:46Z

Verdict: S-0465 meets both acceptance criteria at b41df546; flai verify passed every step there.
- 1: flai/cmd/move.go, flai/cmd/move_refused_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/components/AcceptFindings.svelte, flaiover/src/lib/components/AcceptFindings.svelte.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/routes/items/[id]/+page.svelte, flaiover/src/routes/items/[id]/item.svelte.test.ts, flaiover/src/routes/api/items/[id]/move/+server.ts, flaiover/src/lib/filelinks.ts, flaiover/src/lib/filelinks.test.ts
- 2: design/issues/I-0145-flai-move-to-done-refuses-unfixed-lint-as-text-only-so-the-dashboard-picks-each-file-s-checkout-from-the-refusal-s-words.md, design/issues/summary.md
