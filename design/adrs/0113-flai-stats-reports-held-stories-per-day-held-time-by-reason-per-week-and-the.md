---
id: ADR-0113
title: "flai stats reports held stories per day, held time by reason per week, and the share of stories whose touches were exact per week"
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0081]
---

# ADR-0113 flai stats reports held stories per day, held time by reason per week, and the share of stories whose touches were exact per week

## Context

S-0214 charts parallelism, hold time, and touches drift for E-0016. ADR-0081 has `flai stats` report holds and drift so that the dashboard charts need no figures of their own, and `design/system/metrics.md`, the contract between `flai stats` and the dashboard, changes only with an ADR. S-0205 gave `claims.days[].in_progress`, each story's `items[].held_seconds`, and `claims.drift[]` per story. The charts need what those lack: the stories held at the end of each day, the time held per week split by reason, and per week the share of completed stories whose touches were exact. A ready story can be held for more than one reason at once, by a story it names in `after:` and by an overlapping claim, or by two stories in progress for different reasons, so a split by reason must say where such a hold counts.

## Decision

`flai stats` reports the stories held at each day's end, the seconds held per week by the one reason each hold is named by, and per week the share of completed stories whose touches were exact.

It reports them as `metrics.md` § Claims and touches defines them, refining [ADR-0081](0081-flai-stats-reports-forecast-error-cost-of-delay-waiting-holds-touches-drift-and.md):

- `claims.days[].held`: the stories in `ready` that the hold rules hold at the end of the day, from the same replay as `items[].held_seconds`.
- `claims.weeks[]`, one point per ISO week of the window, with `week`, `start`, and `held_seconds`: per reason, `overlap`, `after`, and `no-touches`, the seconds stories spent held under it in the week, every reason present. A hold with more than one reason counts under the one reason it is named by, the code of `held (…)`: `after` when a story it names in `after:` is not done, whatever else holds it; otherwise that of the first story in progress, by ID, whose claim holds it. So the reasons sum to the time held, and no second counts twice.
- In the same points, of the stories in `claims.drift[]` completed in the week: `stories`, `exact`, those with no file outside their touches and no touch unchanged, and `exact_share`, `exact` over `stories`, absent when `stories` is 0. All three are absent when git cannot be read.

## Consequences

- The Parallelism, Hold Time, and Touches Drift charts read `flai stats --json` alone, and match it.
- The split by reason is the one `flai board` and the launcher show, so a reader can match a bar to the reason a held story displayed. A hold by `after` hides an overlap under it; the time an overlap alone would have cost is not reported.
- The weekly held time comes from the same replay as `items[].held_seconds`, and has the same limits: today's touches and `after`, no recorded holds.
- One weekly array carries both holds and drift, so a dashboard reads one series per week for the claims charts.

## Alternatives considered

- **Count a hold under every reason it has.** Shows every cause, but the reasons no longer sum to the time held and a stacked bar overstates it.
- **Split a hold's seconds evenly between its reasons.** Sums right, but gives fractions of seconds no hold displays and that a reader cannot match to the board.
- **A separate `claims.drift_weeks[]`.** Keeps drift apart from holds, at the cost of a second weekly series with the same weeks.
- **Leave the split to the dashboard.** It cannot replay the holds without the transitions of every item and the hold rules, which ADR-0081 keeps in flai.
