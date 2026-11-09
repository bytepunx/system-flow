---
id: S-0627
type: story
nature: improvement
title: "The dashboard renders repository markdown with markdown-it html: true and no sanitizer, so a document any agent writes can run script as the operator"
status: done
owner: alex
created: 2026-11-07T23:40:28Z
updated: 2026-11-09T21:52:03Z
transitions:
  - to: ready
    at: 2026-11-09T17:18:43Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T17:40:47Z
    by: agent-S-0627
  - to: review
    at: 2026-11-09T17:56:40Z
    by: agent-S-0627
  - to: done
    at: 2026-11-09T21:52:03Z
    by: alex
tags: [dashboard]
topics: [dashboard]
touches: [flaiover/src/lib/markdown.ts, flaiover/src/lib/markdown.test.ts, design/tech/markdown.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/issues/I-0294-the-dashboard-renders-repository-markdown-with-markdown-it-html-true-and-no.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 965
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 2
      hand_edits: 1
      work: 42
  models:
    - model: claude-opus-5-5
      input: 96
      output: 14321
      cache_read: 9840357
      cache_write: 240927
      cost: 4.1823
  strategic:
    - kind: orchestrator
      seconds: 264
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 61
          output: 470
          cache_read: 3015540
          cache_write: 156583
          cost: 0.8383
    - kind: analyzer
      seconds: 41
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
cost_of_delay:
  inputs:
    penalty_per_week: 100
    by: orchestrator
    at: 2026-11-09T17:15:40Z
  value: 100
  by: planner-S-0627
  at: 2026-11-09T17:16:13Z
forecast:
  duration: 40m
  delivery: 2026-11-09T18:28:00Z
  basis: "Its own forecast of 40m; 3rd in the pull order with an in-progress limit of 5, behind S-0611, S-0616, S-0619, S-0620, S-0624, S-0629 and S-0626."
  by: flai
  at: 2026-11-09T17:18:47Z
finalized:
  by: orchestrator
  at: 2026-11-09T17:18:37Z
---
# S-0627 The dashboard renders repository markdown with markdown-it html: true and no sanitizer, so a document any agent writes can run script as the operator

## Goal

This story remediates [I-0294](../../../design/issues/I-0294-the-dashboard-renders-repository-markdown-with-markdown-it-html-true-and-no.md), "The dashboard renders repository markdown with markdown-it html: true and no sanitizer, so a document any agent writes can run script as the operator", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0294 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0294 is closed with `flai issue close I-0294 --reason` saying what fixed it

## Tasks
- T-2752 Render repository markdown with raw HTML escaped, and test that script in a document does not run
- T-2753 Say in the design and the user guide that the dashboard shows raw HTML in markdown as text
- T-2754 Close I-0294 with what fixed it

## Notes

The analyzer found I-0294: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0294: analyzer 41s, 0.8156 USD.

### Planning

Proposed fix: render with markdown-it `html: false`. Every `{@html}` site and `markdownInto` in flaiover take their HTML from `render()` in `flaiover/src/lib/markdown.ts`, so one change covers all 13 sinks. No document under `design/` or `docs/` uses raw HTML outside code spans and fences. The task-list checkboxes and heading permalinks are plugin tokens and keep rendering. A sanitizer such as DOMPurify was weighed and not chosen: it adds a dependency and an allow-list to keep, for HTML no document needs. A Content-Security-Policy is defence in depth, not the cause, and is proposed on TH-0936 as a separate story.

Touches, each a file, no folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flaiover/src/lib/markdown.ts` | design | The issue's Impact names `markdown.ts:15`; the renderer's options |
| `flaiover/src/lib/markdown.test.ts` | co-change | 2 of 3 commits on `markdown.ts`; the reproducing test |
| `design/tech/markdown.md` | design | Lists markdown-it and how flaiover renders; gains the rule |
| `design/system/flaiover-dashboard.md` | co-change | 2 of 3 commits; its Internal structure names the renderer |
| `docs/users/flaiover.md` | co-change | 2 of 3 commits; its Docs section says what a document renders with |
| `design/issues/I-0294-the-dashboard-renders-repository-markdown-with-markdown-it-html-true-and-no.md` | layout | Closed by the second criterion |
| `design/issues/summary.md` | layout | `flai issue close` rewrites it |

`flaiover/src/lib/components/Threads.svelte` co-changes with `markdown.ts` in 2 of 3 commits but is left out: it calls `markdownInto`, which the fix covers without a change there.

No `after`: the story reads no other story's output.

Forecast: 40m. `flai forecast` gave 16m, 104 s per unit of size over 26 done medium improvement stories on claude-opus-5-5, times size 9 (2 criteria, 7 touches). Raised because the code task and the close-out run the flaiover tiers (vitest, svelte-check, flaiover-lint), which take minutes each and which the size does not count. flai's replan keeps the 40m and moves the delivery with the pull order.

Cost of delay: 100 USD a week, as `flai cod` works it out from the one input, `penalty_per_week: 100`, which the orchestrator set at 2026-11-09T17:15:40Z after the planner recommended it on TH-0936. It stands unadjusted. No revenue rides on the story. The figure was scaled against TH-0888's table: below S-0598's 150, since the harm is a risk with no instance observed, and level with S-0603's 100, since script run in the operator's session could accept, move, or edit any item.
