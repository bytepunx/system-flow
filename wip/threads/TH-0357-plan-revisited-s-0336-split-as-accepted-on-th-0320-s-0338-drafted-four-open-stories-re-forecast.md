---
id: TH-0357
title: "Plan revisited: S-0336 split as accepted on TH-0320, S-0338 drafted, four open stories re-forecast"
anchor:
  path: wip/kanban/epics/E-0018-improve-inter-agent-coordination-via-inboxes.md
  item: E-0018
status: resolved
participants: [planner-E-0018, alex]
created: 2026-10-08T04:35:12Z
updated: 2026-10-08T04:42:13Z
---

# TH-0357 Plan revisited: S-0336 split as accepted on TH-0320, S-0338 drafted, four open stories re-forecast

On wip/kanban/epics/E-0018-improve-inter-agent-coordination-via-inboxes.md.

## Entries

### 2026-10-08T04:35:12Z planner-E-0018
S-0336 is split as you accepted on TH-0320. The held card's part is a new draft, S-0338, with four tasks. I revisited S-0334 and S-0337 and refreshed every open story's touches, forecast, and value. `flai check --strict` reports nothing on them.

## Stories and their order

S-0330 to S-0333 and S-0335 are done.

| Story | State | after | Forecast | Value (USD a week) |
|-------|-------|-------|----------|-------------------:|
| S-0336 Messages view and story page | draft | S-0330 (done) | 40m | 217.39 |
| S-0334 a holding agent is asked, and may share | draft | S-0332 (done) | 53m | 288.04 |
| S-0337 coordination stats | finalized | S-0332, S-0333 (done), S-0334 | 51m | 277.17 |
| S-0338 the held card shows the ask and the share | draft, new | S-0334, S-0336 | 40m | 217.39 |

Each value is the story's share of E-0018's 1000 USD a week by forecast duration, out of 3h4m (`flai cod`). S-0336 and S-0334 wait for nothing open, so either can be pulled once you finalize it.

## Tasks and layers

Layers are separated by `→`. Tasks in brackets run together.

- S-0336: T-1211 host reads → T-1212 repo and route → T-1213 view and menu → T-1214 story page and inbox count → T-1215 docs
- S-0338: [T-1324 board fields in flai, T-1325 share in `Messages.svelte`] → T-1326 card → T-1327 docs
- S-0334: T-1202 ADR → T-1203 share and ask → [T-1204 hold honours shares, T-1205 serve asks] → T-1206 docs (unchanged)
- S-0337: T-1216 ADR and `metrics.md` → T-1217 conflict log → T-1218 stats → T-1219 chart → T-1220 docs (unchanged, but T-1220 now also touches `design/system/flai-cli.md`)

## What changed

- **S-0336:**
  - `after` is now S-0330, not S-0334.
  - The held-card criterion and "any share" are gone. The two `BoardCard.svelte` touches are gone.
  - T-1214 is narrowed to the story page and the inbox count.
  - T-1211, T-1213, and T-1215 no longer mention a share. `design/system/flai-cli.md` is added for the host reads.
  - Forecast 48m → 40m.
- **S-0338:** the board's card gains the asking conversation and the share from flai, so `flai board`, `board.get`, and the MCP `board` agree. `Messages.svelte` shows the share.
- **S-0337:** I changed only its touches, figures, and `### Planning` notes. I added `design/system/flai-cli.md`, which describes `flai stats`. Its words are yours and stay as they were.
- **Forecasts:** flai's rate fell from 114 to 104 s per unit over 33 done stories, so S-0334 went 58m → 53m and S-0337 54m → 51m.

## Assumptions

- `workitem` cannot import `messages`. T-1324 passes the asks and shares into `Holds`, as S-0334 passes shares, and widens its touches to each caller it changes. That is why I raised S-0338 from flai's 33m to 40m.
- S-0337 does not wait for S-0336 or S-0338: it reads the message files, not the dashboard.
- `design/adrs` stays a folder touch on S-0334 and S-0337, for ADRs not yet named. It is inside `claims.shared`, so it holds nothing.

## Still open for you

- TH-0319: whether a holding agent's share lifts a hold on its own (S-0334 as written, my recommendation) or waits for your confirmation. S-0338's card works either way.

## Proposals

- Nothing more to split, merge, or drop.
- From TH-0318, still not drafted: a report of the paths most often held on or in conflict, proposing seams. It needs S-0337's data from a few weeks of runs, so it is an analyzer report or a later research story.

### 2026-10-08T04:42:13Z alex
Resolved.
