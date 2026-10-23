---
id: T-2080
type: task
nature: improvement
title: docs/users/flai.md says the patch id leaves out the generated paths, and the flai reference is regenerated
status: done
parent: S-0487
owner: alex
created: 2026-10-23T18:05:12Z
updated: 2026-10-23T18:27:07Z
transitions:
  - to: ready
    at: 2026-10-23T18:26:28Z
    by: agent-S-0487
  - to: in-progress
    at: 2026-10-23T18:26:29Z
    by: agent-S-0487
  - to: done
    at: 2026-10-23T18:27:07Z
    by: agent-S-0487
stream: S-0487
tags: [docs]
touches: [docs/users/flai.md, docs/users/flai-reference.md]
after: [T-2078]
usage:
  source: log
  seconds: 38
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 2204
      cache_read: 491832
      cache_write: 20427
      cost: 0.2745
---
# T-2080 docs/users/flai.md says the patch id leaves out the generated paths, and the flai reference is regenerated

## Work

In `docs/users/flai.md` § Verify a story before review › Resuming across a clean rebase, say that the branch's own change is taken without the generated paths, `design/issues/summary.md` today, so that a round after a rebase that regenerated it still reuses the tiers main's new paths do not select, and bump the file's `updated`. Regenerate `docs/users/flai-reference.md` with `make flai-reference`, since T-2078 changes `flai verify`'s help text.

Waits for T-2078, whose help text the reference is generated from.

## Done when

The section names what the patch id leaves out, `docs/users/flai-reference.md` matches the built flai's help, and `flai test docs/users` passes.

## Notes

The behaviour of the close-out script is unchanged: `scripts/close-out.sh` and the template's need no edit.
