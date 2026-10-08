---
id: T-1306
type: task
nature: improvement
title: The tier lists in devex.md and tooling.md name flaiover-lint, and I-0105 is closed
status: done
parent: S-0319
owner: alex
created: 2026-10-08T00:15:16Z
updated: 2026-10-08T05:55:07Z
transitions:
  - to: ready
    at: 2026-10-08T05:54:48Z
    by: agent-S-0319
  - to: in-progress
    at: 2026-10-08T05:54:49Z
    by: agent-S-0319
  - to: done
    at: 2026-10-08T05:55:07Z
    by: agent-S-0319
stream: S-0319
tags: [docs, testing, issues]
touches: [design/system/devex.md, design/conventions/tooling.md, design/issues/I-0105-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md, design/issues/summary.md]
after: [T-1305]
usage:
  source: log
  seconds: 18
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1798
      cache_read: 658856
      cache_write: 17022
      cost: 0.3039
---
# T-1306 The tier lists in devex.md and tooling.md name flaiover-lint, and I-0105 is closed

## Work

Bring the documents that list this repository's tiers up to date with the tier T-1305 added:

- `design/system/devex.md`: the `Tier commands` row names `scripts/flaiover-lint.sh` among the scripts that take a tier's files, and the `Test tiers` and `Close-out` rows say that a change to `flaiover/` runs prettier and eslint on its files at each `flai test`, while svelte-check stays in the `flaiover` tier under `--all` and at the close-out.
- `design/conventions/tooling.md`, below the marker only: the project addition that lists what `flai test` runs adds the flaiover lint.

Then close I-0105 from the story's worktree: `flai issue close I-0105 --reason "<what fixed it>"`, naming the `flaiover-lint` tier and the script. That rewrites `design/issues/summary.md`.

It waits for T-1305 because the documents name the tier as T-1305 settles it, and the issue closes only once the tier runs.

## Done when

- [ ] `devex.md` and `tooling.md` name the `flaiover-lint` tier and say svelte-check stays in the `--all` tier
- [ ] I-0105 is closed with its reason, and `design/issues/summary.md` no longer lists it
- [ ] `flai test` on the changed markdown passes

## Notes

Drafted by the planner for S-0319. `docs/operators/settings.md` documents the `tests` keys, not this repository's tiers, so it does not change.
