---
id: T-2260
type: task
nature: remediation
title: I-0198 is closed with the reason naming the key's SvelteKit flag, the miss on malformed stamps, and their tests
status: done
parent: S-0528
owner: alex
created: 2026-10-28T18:04:14Z
updated: 2026-10-28T18:13:56Z
transitions:
  - to: ready
    at: 2026-10-28T18:05:33Z
    by: agent-S-0528
  - to: in-progress
    at: 2026-10-28T18:13:48Z
    by: agent-S-0528
  - to: done
    at: 2026-10-28T18:13:56Z
    by: agent-S-0528
stream: S-0528
tags: [flai, components]
touches: [design/issues/I-0198-the-component-props-cache-s-0478-can-serve-a-stale-skeleton-when-only-sveltejs-kit-is-added-or-removed-and-throws-instead-of-missing-on-malformed-stamp-entries.md, design/issues/summary.md]
after: [T-2259]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1552
      cache_read: 292820
      cache_write: 15625
      cost: 0.2006
---
# T-2260 I-0198 is closed with the reason naming the key's SvelteKit flag, the miss on malformed stamps, and their tests

## Work

In the story's worktree, run `flai issue close I-0198 --reason "<reason>"`. The reason says that `props.mjs` now keys its cache by whether the sub-project is SvelteKit and treats a cache file with malformed stamp entries as a miss, and names T-2259's two tests in `flai/internal/components/props_test.go`. `flai issue close` rewrites the issue file and `design/issues/summary.md`; edit neither by hand. Then tick the story's two criteria with `flai criteria tick S-0528 <n>`, each once verified.

It waits for T-2259 because the reason names what it did, and the issue is closed only once the cause no longer occurs.

## Done when

- `design/issues/I-0198-….md` has status closed with the reason, and `design/issues/summary.md` no longer lists it as open.
- Both acceptance criteria of S-0528 are ticked.
- `flai check --strict` is clean on the story.

## Notes

Drafted by planner-S-0528. Issue commands write under the current checkout, so run them in the story's worktree, not the main checkout.
