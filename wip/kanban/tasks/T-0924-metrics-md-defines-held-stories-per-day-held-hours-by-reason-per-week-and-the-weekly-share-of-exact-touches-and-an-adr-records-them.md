---
id: T-0924
type: task
nature: feature
title: metrics.md defines held stories per day, held hours by reason per week, and the weekly share of exact touches, and an ADR records them
status: done
parent: S-0214
owner: alex
created: 2026-10-05T05:45:00Z
updated: 2026-10-07T08:23:03Z
transitions:
  - to: ready
    at: 2026-10-07T08:18:40Z
    by: agent-S-0214
  - to: in-progress
    at: 2026-10-07T08:18:41Z
    by: agent-S-0214
  - to: done
    at: 2026-10-07T08:23:03Z
    by: agent-S-0214
stream: S-0214
tags: [flai]
touches: [design/system/metrics.md, design/adrs]
usage:
  source: log
  seconds: 262
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 63
      output: 26524
      cache_read: 3067617
      cache_write: 115027
      cost: 1.8168
---
# T-0924 metrics.md defines held stories per day, held hours by reason per week, and the weekly share of exact touches, and an ADR records them

## Work

The three charts need figures `flai stats` does not give yet. `claims.days[]` has only `in_progress`, `items[].held_seconds` is not split by reason or by week, and `claims.drift[]` has no weekly share of stories whose touches were exact. ADR-0081 says the dashboard charts need no figures of their own, and `metrics.md` changes only with an ADR. So define these in `metrics.md` § Planning, waiting, and claims › Claims and touches:

- `claims.days[].held`: the stories in `ready` that the hold rules hold at the end of the day.
- `claims.weeks[]`: `week`, `start`, and `held`, the seconds stories spent held in the week per reason (`overlap`, `after`, `no-touches`, the empty claim), every reason present. A hold with more than one reason counts under its first, so that the reasons sum to the held time.
- Per week of the stories completed in it that a commit names: `stories`, `exact` (no file outside and no touch unchanged), and `exact_share`.

Give the precision of each value, and record the additions in a new ADR that refines ADR-0081 (`flai adr new`). Nothing to wait for: the other tasks build on this contract.

## Done when

- [ ] `metrics.md` defines each value with its precision, and its `updated` date is today
- [ ] A new ADR, refining ADR-0081, records the additions and the rule for a hold with more than one reason
- [ ] `flai check --strict` and the markdown lint pass

## Notes
