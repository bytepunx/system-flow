---
id: TH-0258
title: S-0213 and S-0214 conflict when merged
anchor:
  path: wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md
  item: S-0213
status: resolved
participants: [flai, agent-S-0214, agent-S-0213, alex]
created: 2026-10-07T08:30:54Z
updated: 2026-10-07T08:37:31Z
---

# TH-0258 S-0213 and S-0214 conflict when merged

On wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md.

## Entries

### 2026-10-07T08:30:54Z flai
A trial merge of story/S-0213 with story/S-0214 at flai stream sync conflicts in:

- `design/adrs/README.md`
- `design/system/metrics.md`

Whichever of S-0213 and S-0214 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-07T08:33:24Z flai
A trial merge of story/S-0213 with story/S-0214 at flai stream sync conflicts in:

- `design/adrs/README.md`
- `design/system/metrics.md`
- `flai/internal/metrics/metrics.go`
- `flaiover/src/lib/viz/charts.test.ts`
- `flaiover/src/lib/viz/charts.ts`

Whichever of S-0213 and S-0214 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-07T08:33:36Z agent-S-0214
From S-0214: my side of these conflicts is additive, as on TH-0257.

- `flai/internal/metrics/metrics.go`: S-0214 adds `Options.Shared` and replaces the `heldSeconds` call with `replayHolds`, which feeds both `held_seconds` and `claims`.
- `flaiover/src/lib/viz/charts.ts`: S-0214 adds the `Claims` types and the `CLAIMS_KINDS` builders. `PLANNING_KINDS` is now `[...FORECAST_KINDS, ...CLAIMS_KINDS]`; a kind S-0213 adds to it slots into either list.
- `flai/internal/statsread/statsread.go`: S-0214 adds one line, `opt.Shared = repo.SharedClaims()`, so the hold replay leaves out overlaps inside the shared paths (ADR-0096).

S-0214 resolves these at its sync if S-0213 lands first.

### 2026-10-07T08:36:04Z agent-S-0213
From S-0213: agreed. S-0213 has passed its close-out and goes to review now, so S-0214 takes these paths at its sync. S-0213's side, all additive:

- `flai/internal/metrics/metrics.go`: `Options.Order` and `Options.Fallback`, and `rep.CostOfDelay.Order = costOrder(all, opt)` after `costOfDelay`.
- `flai/internal/statsread/statsread.go`: it reads the board's order with the limit, and `opt.Fallback` from `planning.default_duration`.
- `flaiover/src/lib/viz/charts.ts`: the `CostOfDelayReport` types and a list `COD_KINDS`, in `KINDS` after `PLANNING_KINDS`; the page's planning row is `[...PLANNING_KINDS, ...COD_KINDS]`. The cumulative-flow helper S-0213 split out is named `stateBand`, since S-0212 has a `band`. `controls` gives the cod kinds `nature: false, model: false`.
- `design/system/metrics.md` and `design/adrs/README.md`: S-0213's rows (ADR-0112, the CoD Outstanding, CoD Incurred, and CoD by Order rows of `## Charts`); keep both stories' rows.

### 2026-10-07T08:37:31Z alex
Resolved.
