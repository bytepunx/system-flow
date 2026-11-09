---
id: S-0626
type: story
nature: improvement
title: Design documents describe a /conventions page, a /streams page, a mounted-repository readiness check, a flai OTEL span, a dashboard scan of wip/kanban, and a single lockfile that are not there
status: done
owner: alex
created: 2026-11-07T23:40:21Z
updated: 2026-11-09T21:23:03Z
transitions:
  - to: ready
    at: 2026-11-09T17:18:42Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T17:26:31Z
    by: agent-S-0626
  - to: review
    at: 2026-11-09T17:46:51Z
    by: agent-S-0626
  - to: done
    at: 2026-11-09T21:23:03Z
    by: alex
tags: [design, docs, flaiover]
topics: [documentation]
touches: [design/system/flaiover-dashboard.md, flaiover/src/lib/server/design-routes.test.ts, design/conventions/telemetry.md, design/conventions/code-quality.md, design/system/metrics.md, design/system/dashboard-host-channel.md, design/tech/node.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/issues/I-0293-design-documents-describe-a-conventions-page-a-streams-page-a-mounted.md, design/issues/summary.md, design/issues/I-0322-a-change-to-flaiover-dashboard-md-alone-selects-no-vitest-tier-so-design-routes.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1230
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 1
      hand_edits: 2
      work: 65
  models:
    - model: claude-opus-5-5
      input: 246
      output: 43198
      cache_read: 26802083
      cache_write: 361396
      cost: 8.8641
  strategic:
    - kind: planner
      seconds: 23
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 37
          output: 2343
          cache_read: 2179078
          cache_write: 152963
          cost: 0.7694
    - kind: orchestrator
      seconds: 53
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 16
          output: 181
          cache_read: 966029
          cache_write: 31470
          cost: 0.2634
    - kind: analyzer
      seconds: 42
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
    time_lost_per_cycle: 15m
    by: orchestrator
    at: 2026-11-09T17:08:33Z
  value: 37.5
  by: planner-S-0626
  at: 2026-11-09T17:14:36Z
forecast:
  duration: 40m
  delivery: 2026-11-09T18:25:00Z
  basis: "Its own forecast of 40m; 2nd in the pull order with an in-progress limit of 5, behind S-0611, S-0616, S-0619, S-0620, S-0624 and S-0629."
  by: flai
  at: 2026-11-09T17:18:47Z
finalized:
  by: orchestrator
  at: 2026-11-09T17:18:35Z
---
# S-0626 Design documents describe a /conventions page, a /streams page, a mounted-repository readiness check, a flai OTEL span, a dashboard scan of wip/kanban, and a single lockfile that are not there

## Goal

This story remediates [I-0293](../../../design/issues/I-0293-design-documents-describe-a-conventions-page-a-streams-page-a-mounted.md), "Design documents describe a /conventions page, a /streams page, a mounted-repository readiness check, a flai OTEL span, a dashboard scan of wip/kanban, and a single lockfile that are not there", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0293 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0293 is closed with `flai issue close I-0293 --reason` saying what fixed it

## Tasks
- T-2755 The dashboard design's page and API tables name only routes flaiover has, and a test fails on a row with no route
- T-2756 telemetry.md and code-quality.md say what flai and flaiover expose and where flai's integration tests live
- T-2757 metrics.md, dashboard-host-channel.md, and node.md say how the dashboard gets its metrics, that the host channel is built, and which lockfile the build installs from
- T-2758 strategic-agents.md names formula run among the analyzer's refused commands, and flai-cli.md lists flai checks and flai serve checks
- T-2759 I-0293 is closed, naming the corrected documents and the route-table test

## Notes

The analyzer found I-0293: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0293: analyzer 42s, 0.8156 USD.

### Planning

Proposed fix, from the instances in Finding 12 of the analysis, each checked again on main on 2026-11-09. Every instance is a sentence that was true when written and stayed after the code moved on. So the fix makes each document say what is built, and builds nothing the documents describe. One instance has a test that fits: the dashboard design's page and API tables against `flaiover/src/routes`, where two rows (`/conventions`, `/streams`) name routes that are not there. T-2755 adds `flaiover/src/lib/server/design-routes.test.ts`, which fails on such a row, as `TestReferenceIsCurrent` does for `docs/users/flai-reference.md`. The other instances are prose, and no test fits them. `metrics.md` § Data access says how the dashboard reads data and defines no metric, so correcting it needs no ADR. The lockfile sentence in `design/tech/node.md` is corrected here, but removing the stray `package-lock.json` files is S-0630's (I-0297).

Touches, and where each came from (`flai touches suggest` had no declared touch to start from):

| Touch | Source | Why |
|-------|--------|-----|
| `design/system/flaiover-dashboard.md` | design | `/conventions` (line 126) and `/streams` (line 129) rows; no rows for `api/items/[id]/checks` and `checks/tail` |
| `flaiover/src/lib/server/design-routes.test.ts` | layout | The new reproducing test, beside `no-project-files.test.ts` |
| `design/conventions/telemetry.md` | design | Project additions: a flai OTEL span, "flai exposes nothing", a mounted-repository `/_ready` |
| `design/conventions/code-quality.md` | design | Project additions: `flai/tests/integration/` "created when a test outgrows the Go layout" |
| `design/system/metrics.md` | design | § Data access: a dashboard scan of `wip/kanban`, a template fixture |
| `design/system/dashboard-host-channel.md` | design | "Nothing here is built" and "S-0082, in progress" |
| `design/tech/node.md` | design | "a single lockfile format" |
| `design/system/strategic-agents.md` | design | § The analyzer's guard omits `formula run` |
| `design/system/flai-cli.md` | design | § Commands has no `flai checks` or `flai serve checks` row |
| `design/issues/I-0293-…-mounted.md` | criteria | Closed by criterion 2 |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

Left out: `design/tech/sveltekit.md`, whose lockfile line is true, and the template's conventions, since the stale sentences are only in this project's additions.

Tasks: T-2755, T-2756, T-2757, and T-2758 form the first layer and share no file, so they run together. T-2759, closing the issue, is the second, after all four.

Overlaps, left to the trial merge: S-0611 (in progress) touches `flaiover-dashboard.md`; S-0624 touches `metrics.md`; S-0619, S-0620, and S-0629 touch `flai-cli.md`; S-0620 touches `strategic-agents.md`.

Forecast: 40m. `flai forecast` gave 22m on size 13: median 100 s per unit of size over 138 done improvement stories on claude-opus-5-5 in the large band. Raised, because each of nine documents needs its claims read against the code before it is rewritten, and the new test parses two markdown tables and walks the routes. Its close-out also runs the flaiover tiers. Before the touches, it gave 3m on size 2.

Cost of delay: 37.50 USD a week, as `flai cod` gave it. That is 15m of time lost per 168h cycle at 150 USD an hour, 1.00 cycle a week. The input `time_lost_per_cycle: 15m` is the one recommended on TH-0930, which the orchestrator set. Kept as given.

No `after`: the story reads nothing another open story writes.
