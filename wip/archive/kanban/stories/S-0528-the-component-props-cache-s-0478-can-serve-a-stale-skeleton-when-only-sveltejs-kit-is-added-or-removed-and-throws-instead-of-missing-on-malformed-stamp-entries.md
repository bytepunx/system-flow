---
id: S-0528
type: story
nature: remediation
title: The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries
status: done
owner: alex
created: 2026-10-28T12:54:18Z
updated: 2026-10-28T18:45:41Z
transitions:
  - to: ready
    at: 2026-10-28T18:05:03Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-28T18:05:08Z
    by: agent-S-0528
  - to: review
    at: 2026-10-28T18:15:18Z
    by: agent-S-0528
  - to: done
    at: 2026-10-28T18:45:41Z
    by: orchestrator
tags: [flai, components]
topics: [components]
touches: [flai/internal/components/props.mjs, flai/internal/components/props_test.go, design/system/component-library.md, design/issues/I-0198-the-component-props-cache-s-0478-can-serve-a-stale-skeleton-when-only-sveltejs-kit-is-added-or-removed-and-throws-instead-of-missing-on-malformed-stamp-entries.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 619
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 4
      hand_edits: 1
      work: 23
  models:
    - model: claude-opus-5-5
      input: 92
      output: 22670
      cache_read: 4277658
      cache_write: 228261
      cost: 2.9306
  strategic:
    - kind: orchestrator
      seconds: 1828
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 168
          output: 2923
          cache_read: 18061636
          cache_write: 45251
          cost: 3.9776
        - model: claude-sonnet-5-5
          input: 8
          output: 51
          cache_read: 70943
          cache_write: 36759
          cost: 0.074
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-10-28T18:03:09Z
  value: 12.5
  by: planner-S-0528
  at: 2026-10-28T18:04:34Z
forecast:
  duration: 40m
  delivery: 2026-10-29T08:50:00Z
  basis: "Planner's 40m over flai's 14m (median 118 s per unit over 33 done remediation stories on claude-opus-5-5 in the medium band, size 7): a small props.mjs fix, but two Go tests that run node and the fixture's typescript and must fail first, a design edit, the issue close, and a close-out that runs the components tests; 37th in the pull order with an in-progress limit of 5, delivery from flai's play-out start plus 40m times the cycle factor 1.35."
  by: planner-S-0528
  at: 2026-10-28T18:04:34Z
finalized:
  by: alex
  at: 2026-10-28T12:57:41Z
---
# S-0528 The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries

## Goal

This story remediates [I-0198](../../../design/issues/I-0198-the-component-props-cache-s-0478-can-serve-a-stale-skeleton-when-only-sveltejs-kit-is-added-or-removed-and-throws-instead-of-missing-on-malformed-stamp-entries.md), "The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0198 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0198 is closed with `flai issue close I-0198 --reason` saying what fixed it

## Tasks
- T-2259 props.mjs keys its cache by whether the sub-project is SvelteKit and misses on a cache file whose stamps are malformed
- T-2260 I-0198 is closed with the reason naming the key's SvelteKit flag, the miss on malformed stamps, and their tests

## Notes

### Planning

Proposed fix, from I-0198's one instance: add the answer of `isSvelteKit()` to the cache key in `cachedSkeletons`, and make `readCache` return a miss for a cache file whose stamp entries are malformed. T-2259 describes both and the tests that reproduce them.

Touches, all files, no folder:

| Touch | Source | Why |
| --- | --- | --- |
| `flai/internal/components/props.mjs` | layout | The issue names `isSvelteKit()`, `readCache`, and the key in this script. |
| `flai/internal/components/props_test.go` | layout | The S-0478 cache tests live here; the two new tests go beside them. |
| `design/system/component-library.md` | co-change (2 of 3, 67%), design | § Props lists what the key holds and what a miss is. |
| `design/issues/I-0198-….md` | criterion 2 | `flai issue close` rewrites it. |
| `design/issues/summary.md` | criterion 2 | `flai issue close` regenerates it. |

Left out: `flai/internal/components/props.go`, co-change 3 of 3. `scriptHash` is the script's own SHA-256, so a changed `props.mjs` keys new cache files without a Go change.

Forecast: 40m, up from flai's 14m. The fix is a few lines, but each new test runs `node` with the fixture's `typescript`, must be seen failing first, and the close-out runs the components package's tests. S-0525 and S-0527 were planned at 45m and 40m for work of this shape.

Cost of delay: 12.50 USD a week from `flai cod`, on the input `time_lost_per_cycle: 5m` that the orchestrator set on TH-0699, as recommended there. It stands: I-0198 occurred once and costs a re-run or a wrong example skeleton, not lost work.

No `after`: the story needs nothing that is not already on main.

### Accepted by the orchestrator

- Verified: 54322d12b18d9d96f4d9945fed898c8b8299bfab
- Batch: 20261028-183207
- Enqueued verified: 80f8ff534bc19af7be50c40e69b23588ae15ad53
- At: 2026-10-28T18:45:41Z

Verdict: pass — props.mjs keys its cache by the SvelteKit answer and takes malformed stamps as a miss, two Go tests reproduce both causes, and I-0198 is closed with a reason.
- 1: flai/internal/components/props.mjs, flai/internal/components/props_test.go, design/system/component-library.md
- 2: design/issues/I-0198-the-component-props-cache-s-0478-can-serve-a-stale-skeleton-when-only-sveltejs-kit-is-added-or-removed-and-throws-instead-of-missing-on-malformed-stamp-entries.md, design/issues/summary.md
