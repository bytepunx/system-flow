---
title: "Bottlenecks: waiting on the operator, acceptances that stop, and holds in ready"
updated: 2026-10-15
status: active
focus: bottlenecks
from: 2026-09-15
to: 2026-10-15
---

# Bottlenecks in the flow of work, 2026-09-15 to 2026-10-15

Every figure is from `flai stats --json` run at `2026-10-15T22:41:42Z` over its default 30-day window (`window_start` `2026-09-15T22:41:42Z`), type `story`. Hours in brackets are the same seconds divided by 3600, for reading only.

## Where the flow stands

| Figure | Value |
|--------|-------|
| Completed / cancelled | 371 / 12 (`cancellation_rate` 0.0313) |
| Throughput per week | 86.57 (W38 64, W39 57, W40 86, W41 106, W42 58 to date) |
| Cycle time p50 / p85 / max | 2,303 s (0.64 h) / 12,990 s (3.6 h) / 407,617 s (113 h, S-0232) |
| Lead time p50 / p85 | 45,649 s (12.7 h) / 232,557 s (64.6 h) |
| Queue time p50 / p85 / max | 1,366 s / 18,985 s (5.3 h) / 305,070 s (84.7 h, S-0216) |
| Time-in-state share | backlog 0.782, ready 0.117, in-progress 0.050, review 0.051 |
| Flow efficiency | 0.978 (blocked time is S-0241 alone, 660,766 s, an experiment waiting for forks) |
| WIP now | 1 (S-0391), limits ready 10, in-progress 5, review 5 (`wip/kanban/board.md`) |
| Empty wakes | 0 |

Work itself is not the bottleneck: a story is worked in under an hour at the median and under four at p85. What stretches lead time is waiting, in three places: on the operator, on an acceptance that stopped, and in `ready` under holds. The cumulative flow shows the cost: `done` stood at 318 on 2026-10-08 and 2026-10-09 and 319 on 2026-10-10 and 2026-10-11, with 2 stories in progress and 3 to 4 in review each day, then rose to 350 on 2026-10-12. The turns per day over those days were 1,823, 90, 2, 0, 1,544.

## Finding 1: a story in review waits days when the orchestrator's acceptance stops

**Severity: high.** Issue: [I-0150](../issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md).

`waiting.weeks[]` for W42 (from 2026-10-12): `review.total_seconds` 973,317 (270 h) over 58 completed stories, `mean_seconds` 16,781, against 2,910 in W41 and 3,032 in W40. Three stories carry 933,886 s of it, 96 percent:

| Story | `kind` | `started` | `ended` | `seconds` | `awaited` |
|-------|--------|-----------|---------|-----------|-----------|
| S-0341 | review | 2026-10-08T09:41:28Z | 2026-10-12T14:12:27Z | 361,859 (100.5 h) | alex |
| S-0342 | review | 2026-10-08T21:21:04Z | 2026-10-12T14:33:30Z | 321,146 (89.2 h) | alex |
| S-0343 | review | 2026-10-09T16:52:20Z | 2026-10-12T14:33:41Z | 250,881 (69.7 h) | alex |

S-0341's `time_in_state_seconds` are backlog 17, ready 2,268, in-progress 3,851, review 361,859: 98 percent of its lead time was review. `orchestration.permissions.accept_reviews` has been on since `9c28ced8` (2026-10-07), and `wip/agents/orchestrator.md` shows the orchestrator ran on 2026-10-09, 2026-10-10, and until 07:32Z on 2026-10-11, so the permission was not the cause. Each acceptance stopped on something the orchestrator could not clear:

- S-0341: at 09:43Z on 2026-10-08 the orchestrator declined to accept and opened TH-0387, asking for a clean-tree gate on the verifier's reuse of tiers; it held the story again on 2026-10-09, 2026-10-10, and 2026-10-12. The gate was never built; the operator accepted the story as it was.
- S-0342: held while main carried S-0338's uncommitted acceptance (TH-0389, TH-0391), then refused in the dry run on the operator's uncommitted files (TH-0394); when the orchestrator finally ran `flai accept S-0342` on 2026-10-10 it stopped mid-rebase on `design/system/flai-cli.md` (TH-0397). `coordination.weeks[]` W42 records `conflicts.met_accept` 1, `met_accept_paths` 1. `design/system/flai-cli.md` is in the manifest's `claims.shared`, which never holds (ADR-0096), so the conflict surfaces at acceptance and the orchestrator cannot finish it.
- S-0343: held by the same dirty main (TH-0394), then needed a sync onto S-0289 (TH-0398).

The design says what should happen instead. `design/system/workflow.md` § What waits on an acceptance: "With the `orchestrate` host action on and `orchestration.permissions.accept_reviews` on, the orchestrator accepts a story in review under ADR-0093's conditions, and the stories waiting on it start without a person. A wait remains only for a story the orchestrator cannot vouch for, which it leaves in review with a thread saying what is missing". Here all three fell to "cannot vouch for", and each then waited for the operator's next session. I-0088, which recorded the operator-only review wait, is closed (count 2); this is the same wait returning through the orchestrator's refusals rather than its absence.

**Impact.** 933,886 s (259 h) of review wait in one week over three stories, and the board's `done` count moved by one story over 2026-10-09 to 2026-10-11. Estimated time lost per cycle (`planning.cycle` 168h): 96 h, the calendar time from S-0341's move to review to the last of the three acceptances, less the operator's working time on 2026-10-12.

## Finding 2: a pending recommendation waits a day or more for the operator while its story idles in progress

**Severity: high.** Issue: [I-0151](../issues/I-0151-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md).

`waiting.weeks[]` W42: `threads.total_seconds` 1,003,576 (278.8 h), `mean_seconds` 17,303 per completed story, against 1,492 in W41 and 3,137 in W40. `threads.confirmed` in W42: `count` 4, `total_seconds` 421,270 (117 h); the orchestrator's own answers ended one wait of 7 s. `waiting.longest[]`, every wait awaited `alex`:

| Story | Thread | `started` | `ended` | `seconds` | What waited |
|-------|--------|-----------|---------|-----------|-------------|
| S-0232 | TH-0324 | 2026-10-07T22:23:53Z | 2026-10-12T15:04:07Z | 405,614 (112.7 h) | The operator's cosign public key and two scope choices; the orchestrator recommended within 21 s |
| S-0337 | TH-0390 | 2026-10-08T20:57:12Z | 2026-10-12T15:29:06Z | 325,914 (90.5 h) | An `rm` of an empty ref under `.git` after the host's crash; recommended within 14 s, confirmed 2026-10-12 |
| S-0380 | TH-0445 | 2026-10-13T18:16:41Z | 2026-10-14T20:30:34Z | 94,433 (26.2 h) | Ten design choices; recommended at once, the operator overruled one |
| S-0378 | TH-0436 | 2026-10-13T16:59:05Z | 2026-10-14T18:23:40Z | 89,637 (24.9 h) | A plan with two decisions; recommended at once, confirmed |

S-0232's `cycle_time_seconds` is 407,617, of which `wait_threads_seconds` is 405,614: the window's maximum cycle time is one thread wait. Today's instance: `aging[]` lists S-0391 in progress for 63,454 s (17 h 37 m), `over_p85` true, waiting on TH-0461, where the orchestrator's recommendation (run `flai host upgrade` once) has been pending since 2026-10-15T05:09:11Z.

These questions needed the operator: a secret, a write under `.git`, a scope change. `design/system/metrics.md` § Waiting defines why they count: "A recommendation ends no wait" ([ADR-0091](../adrs/0091-a-recommendation-ends-no-thread-wait-and-flai-stats-reports-the-waits-the.md)), so each wait runs until the operator confirms. The orchestrator answers in seconds; the confirmation takes a day on a working day and four over a weekend, and each story's remaining tasks, and every ready story that names it in `after`, wait with it.

**Impact.** 421,270 s (117 h) of confirmed-recommendation waits in W42 across four stories, a mean of 105,317 s (29.3 h) each. Estimated time lost per cycle: 117 h of story time, taken as W42's `threads.confirmed.total_seconds`.

## Finding 3: folder-wide touches and `after` hubs held up to ten ready stories a day while three of five lanes stood empty

**Severity: medium, falling.** Issue: [I-0152](../issues/I-0152-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md).

`claims.weeks[].held_seconds`:

| Week | `overlap` | `after` | `no-touches` | `exact_share` |
|------|-----------|---------|--------------|---------------|
| W39 | 362,132 | 91,161 | 894 | 0.228 |
| W40 | 653,549 | 582,160 | 0 | 0.529 |
| W41 | 690,335 | 887,656 | 0 | 0.491 |
| W42 (to date) | 25,902 | 133,264 | 0 | 0.603 |

W41 held ready stories for 1,577,991 s (438 h) in all. `claims.days[]` over the same days: `held` 9 on 2026-10-03, 9 on 2026-10-05, 10 on 2026-10-06, while `in_progress` was 0, 1, and 2 against a limit of 5 (`claims.limit` 5). 221 of the 391 items carry `held_seconds`, 3,532,742 s together. The stories held longest, `items[].held_seconds`: S-0138 180,695; S-0137 172,010; S-0136 171,537; S-0135 170,672; S-0227 141,332; S-0229 141,310; S-0224 139,613; S-0226 120,787; S-0223 115,960.

Their archived front matter (`wip/archive/`) shows both hold reasons at once. S-0207 to S-0211 and S-0220 to S-0229 declare folder-wide touches, `flai/cmd`, `flai/internal/{harness,hostapi,mcpserver,serve,guard,manifest}`, `flaiover/src/routes`, `docs/users`; S-0218 is an `after` hub that S-0220, S-0221, S-0222, S-0226, S-0228, and S-0229 all name, and S-0223 is another for S-0224 and S-0227. `design/system/metrics.md` § Claims and touches: "A hold counts under the one reason it is named by: `after` when a story it names in `after` is not done, whatever else holds it; otherwise the reason of the first story in progress, by ID, whose claim holds it". So one story in progress with a folder touch, or one hub not yet accepted, held the rest of the column. The `exact_share` of 0.49 to 0.60 says half the completed stories still changed files outside their touches or left declared touches unchanged, so the claims are both too wide to run in parallel and not wide enough to be true.

I-0087 recorded this once through `docs/users/flai.md` and folder touches and is closed (count 3); the shared-paths list and ADR-0096 answered that instance. The W40 and W41 holds came through the planner's touches and `after` on E-0016's stories instead. W42 shows the remedy working: `coordination.weeks[]` W42 `shares.started` 11, `shares.held_seconds` 179,046 (the hold each would have waited, avoided under ADR-0134), and `held_seconds` fell to 159,166 s.

**Impact.** 1,577,991 s (438 h) of story time held in W41 with three lanes free; 159,166 s (44 h) so far in W42. Estimated time lost per cycle: 44 h, W42's held time, the current rate.

## Finding 4: ready stories waited up to 3.5 days unheld while lanes were free

**Severity: low, recorded as an impression.** Issue: [I-0153](../issues/I-0153-ready-stories-waited-up-to-3-5-days-unheld-with-in-progress-lanes-free-and-review-under-its-limit.md).

`items[]` with the longest `queue_time_seconds` (`started - committed`) are S-0216 305,070, S-0215 303,610, S-0214 301,459, S-0213 298,706, and S-0212 295,977, the flaiover chart stories of E-0016, all completed on 2026-10-07. Their `held_seconds` are 35,693, 34,246, 32,110, 29,392, and 26,665: roughly 270,000 s (75 h) of each one's ready time was not a hold. Over 2026-10-04 to 2026-10-06, `claims.days[].in_progress` was 1, 1, 2 of 5 and `cfd[].counts.review` 1, 0, 0 of 5, so neither limit was full. `design/system/workflow.md` § The pull order says "`flai serve` starts each ready story's agent in this order ... only while the in-progress limit leaves room and review is under its limit, and not while the story is held". The metrics do not say why these were not started: flai serve's own runs are not in the activity documents, and the orchestrator ran only from 2026-10-07. The cause is not measured, so this is filed as an impression, with the figures, for the operator to confirm or dismiss.

**Impact.** About 75 h of unheld ready time on each of five stories in W41. Estimated time lost per cycle: 75 h, one such story's unheld wait, if it recurs.

## Not bottlenecks

- **Backlog share.** `time_in_state_share.backlog` is 0.782 of lead time, but it is the planning queue: 41 stories entered the backlog on 2026-10-02 (`cfd[]` backlog 0 on 2026-10-01, 41 on 2026-10-02), the planner drafted them, and the column is down to 1 today. `cost_of_delay.weeks[].incurred` was 767.89 USD in W40, 1,536.77 in W41, 378.47 in W42, and `without_value.backlog` fell from 13 on 2026-10-02 to 0 by 2026-10-05.
- **Agent ceremony.** `turns.total`: 5,772 turns, 5,169 `work`, 289 `ceremony`, 217 `hand_edits`, 97 `test_runs`, 0 `empty_wakes`.
- **Blocking.** Only S-0241 (an experiment that waits for forked sub-agents), S-0052 (19,523 s), and S-0190 (15,813 s) carry `blocked_seconds`.
- **Conflicts.** W42: `conflicts.found` 9 at trial merge, `met_sync` 10, `met_accept` 1. The one at acceptance is in Finding 1; the rest were met where the design puts them, at sync.

## What the operator should do next

1. Confirm or answer TH-0461 so S-0391 leaves `in-progress`; it has waited 17 h on a pending recommendation.
2. Choose at the next acceptance which of I-0150 to I-0153 become stories; the planner writes them.
