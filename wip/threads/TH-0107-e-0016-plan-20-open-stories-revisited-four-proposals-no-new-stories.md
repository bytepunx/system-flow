---
id: TH-0107
title: "E-0016 plan: 20 open stories revisited, four proposals, no new stories"
anchor:
  path: wip/kanban/epics/E-0016-introduce-orchestration-planning-and-analysis-agents.md
  item: E-0016
status: resolved
participants: [planner-E-0016, alex]
created: 2026-10-04T04:45:38Z
updated: 2026-10-05T04:06:46Z
---

# TH-0107 E-0016 plan: 20 open stories revisited, four proposals, no new stories

On wip/kanban/epics/E-0016-introduce-orchestration-planning-and-analysis-agents.md.

## Entries

### 2026-10-04T04:45:38Z planner-E-0016
I revisited E-0016's 20 open stories against its outcome. The 17 not in progress now have predicted touches and a forecast. Their cost of delay values wait on `planning.hour_rate` (TH-0106). I created no stories: the existing ones cover the outcome.

## What changed

- **Touches** were extended on all 17 stories, and none were removed. They now include the code each story reaches, its docs, and its manifest settings.
- The chart stories S-0212 to S-0216 named `flaiover/src/lib/charts`, which does not exist. Charts live in `flaiover/src/lib/viz/charts.ts` and the menu in `flaiover/src/lib/sitemenu.ts`, so I added both. I kept the old path.
- **Forecasts** use agent time from E-0016's ten accepted stories (1197 to 5186 s, median about 3900 s), scaled by each story's criteria and new surfaces.
- **Not edited:** S-0209, S-0225 and S-0255, which are in progress.

## Order the forecasts assume

Most flai-side stories share `flai/internal/harness`, `hostapi`, `mcpserver` or `flai/cmd`, so overlap holds make them run one at a time. Parallelism comes from the dashboard lane.

| Lane | Stories in order | Last delivery |
|------|------------------|---------------|
| flai | S-0210 → S-0217 → S-0211 → S-0223 → S-0218 → S-0219 → S-0220 → S-0221 → S-0222 → S-0224 → S-0226 → S-0227 → S-0229 | 2026-10-05T16:00Z |
| dashboard | S-0212 → S-0214 → S-0215 → S-0213 → S-0228 → S-0216 | 2026-10-05T14:30Z |

This is the forecasts' assumption, not a reorder: the pull order stays yours.

## Assumptions

- S-0209, S-0225 and S-0255 are accepted by about 06:00Z today, which frees S-0210 and the chart lane.
- You finalize and promote each story as its `after` clears. A review takes about an hour, as observed on S-0199 to S-0208.
- S-0214's hold time by reason is a new metric, so it needs an ADR (metrics.md changes need one). The other charts' figures are already in `flai stats` from S-0205, except the cod-order projection.

## Proposals (none applied; each needs your word)

1. **Rewrite S-0226 and S-0227.** Their criteria were copied from S-0225 ("a orchestrator", and "shows the planner's expected cost"). S-0225 already makes the `strategic` usage entry generic per kind (`StrategicKinds` covers planner, orchestrator and analyzer) and shows it in stats and charts. What is left:
   - S-0226: charge each orchestrator activity to the items its `activity_log` names, and to the project when it names none.
   - S-0227: issues carry no usage today. Add usage to issues, with an ADR, and carry it to the story made from an issue.
2. **Make S-0213 wait for S-0217.** The cod-order chart compares orders by cost of delay and by WSJF. That is `flai order --by` arithmetic, and you chose on 2026-10-02 that this arithmetic lives in flai. The chart should read it, not recompute it.
3. **Make S-0223 wait for S-0211.** Both add a schedule to `flai serve`, which has no scheduler today. One scheduler, built once, is cheaper than two.
4. **Split S-0228.**
   - The Planner page needs only S-0208 and S-0225, so it could ship now.
   - The Orchestrator and Analyzer pages keep `after: S-0218, S-0223`.

   As it stands, the Planner page waits about 20 hours for the orchestrator.

I propose no merges and no drops.

## Risks

- `item_edit` refuses a touch that starts with a dot, such as `.claude/agents/planner.md`, though `flai touches` accepts it. I set those stories' touches (S-0210, S-0218, S-0219, S-0220, S-0223, S-0224) with `flai touches`.
- S-0218 and S-0223 add MCP tools for new roles. If the guard hook's matcher in `.claude/settings.json` has to change, a story agent cannot ship that change (I-0069), and you will need to paste it.
- `flai check --strict` reports nothing on the revisited stories. It does fail today, on overlaps among S-0209, S-0225 and S-0255 and on TH-0094.

## What you need to do

- Answer the hour rate on TH-0106.
- Reply here on proposals 1 to 4.

### 2026-10-04T04:48:36Z planner-E-0016
Two corrections and an update:

- **Counts:** E-0016 has 22 open stories. All 22 were revisited, and the 19 not in progress (S-0210 to S-0224 and S-0226 to S-0229) got touches and forecasts. My first entry said 20 and 17.
- **Cost of delay:** with `planning.hour_rate: 150`, E-0016's value is 1500 USD per week. Each of the 19 stories now has its share, by forecast duration (table on TH-0106).

Proposals 1 to 4 still wait for your word. Reply "yes" to all of them, or name the ones to apply, and I will make the edits: rewriting S-0226 and S-0227, adding `after` to S-0213 and S-0223, and splitting S-0228 with a new draft for the Planner page.

### 2026-10-04T04:51:47Z alex
yes

### 2026-10-04T04:53:44Z planner-E-0016
All four proposals are applied. The lint is clean, and `flai check --strict` reports nothing on these items.

1. **S-0226 and S-0227 are rewritten.**
   - S-0226 charges each orchestrator activity to the items it names, split evenly between them. An activity that names no item goes to a project strategic total, under an ADR that extends ADR-0083. I added `design/adrs` and `flai/internal/metrics` to its touches.
   - S-0227 gives issues `usage` with a `strategic` entry, under an ADR, and charges the analyzer's activities to the issues they name. The story made from an issue carries that usage, counted once in totals.
   - S-0227 now also waits for S-0226, which adds the project total.
2. **S-0213** now waits for S-0205 and S-0217.
3. **S-0223** now waits for S-0206, S-0207 and S-0211.
4. **S-0228 is split.**
   - New draft **S-0259**, the Planner page: `after: S-0208, S-0225`, which are both done. Forecast 45m, delivery about 09:30Z, value 36.59 USD/week.
   - S-0228 is retitled to the Orchestrator and Analyzer pages, and waits for S-0259 too. Forecast 1h15m, value 60.98 USD/week.
   - The two values add up to S-0228's old share, so the other stories' values stand.

S-0259 is a draft in the backlog, ready for you to finalize.

### 2026-10-05T04:06:46Z alex
Resolved.
