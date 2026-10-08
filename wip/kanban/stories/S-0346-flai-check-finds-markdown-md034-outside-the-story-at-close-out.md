---
id: S-0346
type: story
nature: improvement
title: "flai check finds `markdown.MD034` outside the story at close-out"
status: review
owner: alex
created: 2026-10-08T08:08:21Z
updated: 2026-10-08T09:42:18Z
transitions:
  - to: ready
    at: 2026-10-08T08:59:57Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T09:29:30Z
    by: agent-S-0346
  - to: review
    at: 2026-10-08T09:42:18Z
    by: agent-S-0346
tags: [flai, mdlint, serve]
topics: [markdown, planning]
touches: [flai/internal/mdlint/inline.go, flai/internal/mdlint/quote.go, flai/internal/mdlint/quote_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, design/system/agent-narrative.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/summary.md, flai/internal/mdlint/rules.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 790
  turns:
    - day: 2026-10-08
      test_runs: 1
      hand_edits: 1
      work: 35
  models:
    - model: claude-opus-5-5
      input: 76
      output: 17344
      cache_read: 5544624
      cache_write: 185900
      cost: 2.9433
  strategic:
    - kind: orchestrator
      seconds: 359
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 209
          output: 3657
          cache_read: 61215993
          cache_write: 64871
          cost: 15.0987
        - model: claude-sonnet-5-5
          input: 7
          output: 38
          cache_read: 115837
          cache_write: 27278
          cost: 0.1234
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-10-08T08:54:56Z
  value: 12.5
  by: planner-S-0346
  at: 2026-10-08T08:55:44Z
forecast:
  duration: 40m
  delivery: 2026-10-08T10:01:00Z
  basis: "Its own forecast of 40m; 2nd in the pull order with an in-progress limit of 5, behind S-0232, S-0322, S-0341, S-0344, S-0348 and S-0338."
  by: flai
  at: 2026-10-08T09:18:44Z
finalized:
  by: orchestrator
  at: 2026-10-08T08:55:53Z
---
# S-0346 flai check finds `markdown.MD034` outside the story at close-out

## Goal

This story remediates [I-0118](../../../design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md), "flai check finds `markdown.MD034` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0118 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0118 is closed with `flai issue close I-0118 --reason` saying what fixed it

## Tasks
- T-1419 mdlint quotes the bare URLs md034 finds in a line, so flai can make text it takes lint clean
- T-1420 A strategic run's end entry quotes the bare URLs in the summary flai takes from its final text, instead of being refused and lost
- T-1421 Close I-0118 saying what fixed it

## Notes

### Planning

Proposed cause, from I-0118's five instances:

- Each finding is a bare `www.` literal in `wip/`: four times line 1529 of `wip/agents/orchestrator.md`, once the title of TH-0361. All five fall between 04:23Z and 05:31Z on 2026-10-08.
- The installed flai that wrote those lines (1.38.1, then 1.39.3) predates S-0324, published in 1.39.4, which gave flai's MD034 the bare `www.` rule (I-0110). The close-out ran the tree's flai, which has the rule, so it found them.
- The host now runs 1.39.8. Every write an agent makes to `wip/` (thread titles and entries, `activity_log`, items, messages) passes through `LintGuard`, which refuses a bare URL. `flai check --strict` on main finds no MD034 today.
- One write is left that no agent checks: the entry `flai serve` logs when a strategic run ends takes its summary from the run's final text (`logRunEndSaying` in `flai/internal/serve/activity.go`). A bare URL there is refused, and the run's seconds and cost go unlogged with only a warning. T-1419 and T-1420 make flai quote it instead.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/mdlint/inline.go` | layout: `parseRange` records bare URLs for `md034` |
| `flai/internal/mdlint/quote.go` | layout: new file for `QuoteBareURLs` |
| `flai/internal/mdlint/quote_test.go` | layout: its test |
| `flai/internal/serve/activity.go` | layout: `logRunEndSaying` |
| `flai/internal/serve/activity_test.go` | co-change: 6 of 9 commits to `activity.go` (67%) |
| `design/system/agent-narrative.md` | design: § Strategic agents' activity documents says where the run-end summary comes from |
| `design/issues/I-0118-…md`, `design/issues/summary.md` | goal: criterion 2 closes the issue |

The story declared no touches. `flai touches suggest` from `activity.go` also listed `design/system/flai-cli.md` (44%) and `flai/internal/serve/plan.go` (44%). I left them out: the change sits in `logRunEndSaying`, which all three kinds share, and `agent-narrative.md` is where the run-end entry is described. No folder touch is kept.

Figures:

- Forecast: flai gave 15m (median 86 s per unit of size over 11 done improvement stories on claude-opus-5-5, size 10). I raised it to 40m. The new helper must match `md034`'s three kinds of bare URL exactly, T-1420 needs a reproducing serve test, and the close-out runs the integration tier. That is in line with the 30m set on S-0345 and S-0290, which are smaller.
- Cost of delay: 12.50 USD a week, as `flai cod` works it out. The input is `time_lost_per_cycle: 5m` per 168h cycle at 150 USD an hour. The orchestrator set it on TH-0381, from my recommendation. It stands unadjusted: a finding outside the story is only a note (ADR-0085), and no instance has recurred since the host's flai gained S-0324's rule.

Assumptions:

- Agent-written `activity_log` summaries stay refused on a finding, as `tooling.md` says. Only the text flai takes from a run is quoted.
- `threads.MirrorNarrative` writes a narrative's `## Open questions` without `LintGuard`. The titles it copies were already linted in the thread file, and a close-out leaves out findings on another open story's narrative (ADR-0123). So it is out of scope.
