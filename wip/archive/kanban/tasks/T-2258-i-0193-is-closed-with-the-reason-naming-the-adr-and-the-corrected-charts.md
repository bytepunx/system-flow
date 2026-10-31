---
id: T-2258
type: task
nature: remediation
title: "I-0193 is closed with the reason naming the ADR and the corrected § Charts"
status: done
parent: S-0527
owner: alex
created: 2026-10-28T15:01:43Z
updated: 2026-10-30T00:15:49Z
transitions:
  - to: ready
    at: 2026-10-30T00:15:38Z
    by: agent-S-0527
  - to: in-progress
    at: 2026-10-30T00:15:40Z
    by: agent-S-0527
  - to: done
    at: 2026-10-30T00:15:49Z
    by: agent-S-0527
stream: S-0527
tags: [issues]
touches: [design/issues/I-0193-metrics-md-charts-says-burn-up-draws-a-forecast-line-from-throughput-and-cycle-time-filters-by-nature-but-neither-chart-does.md, design/issues/summary.md]
after: [T-2256, T-2257]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 608
      cache_read: 218683
      cache_write: 6719
      cost: 0.1097
---
# T-2258 I-0193 is closed with the reason naming the ADR and the corrected § Charts

## Work

In the story's worktree, run `flai issue close I-0193 --reason "<reason>"`, the reason naming the ADR T-2255 wrote, that `design/system/metrics.md` § Charts now says what each chart draws with a row for every chart kind, and that T-2257's test in `flaiover/src/lib/viz/explain.test.ts` holds it. `flai issue close` rewrites the issue file and `design/issues/summary.md`; edit neither by hand. Then tick the story's two criteria with `flai criteria tick S-0527 <n>`, each once verified.

It waits for T-2256 and T-2257 because the reason names what they did, and the issue is closed only once the cause no longer occurs.

## Done when

- `design/issues/I-0193-….md` has status closed with the reason, and `design/issues/summary.md` no longer lists it as open.
- Both acceptance criteria of S-0527 are ticked.
- `flai check --strict` is clean on the story.

## Notes

Drafted by planner-S-0527. Issue commands write under the current checkout, so run them in the story's worktree, not the main checkout.
