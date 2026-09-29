---
id: T-0585
type: task
nature: feature
title: A patch is read as runs of consecutive added and removed lines, each line's sign apart from its text
status: done
parent: S-0164
owner: alex
created: 2026-09-29T21:12:38Z
updated: 2026-09-29T21:16:02Z
transitions:
  - to: ready
    at: 2026-09-29T21:13:08Z
    by: agent-S-0164
  - to: in-progress
    at: 2026-09-29T21:13:09Z
    by: agent-S-0164
  - to: done
    at: 2026-09-29T21:16:02Z
    by: agent-S-0164
stream: S-0164
tags: []
touches: [flaiover/src/lib]
usage:
  source: log
  seconds: 173
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 15
      output: 5982
      cache_read: 842377
      cache_write: 26220
      cost: 1.0342
---

# T-0585 A patch is read as runs of consecutive added and removed lines, each line's sign apart from its text

## Work

In `flaiover/src/lib/review.ts`, read a unified patch as what the panel draws: each line with its kind, its sign (`+`, `-`, or none), and its text without the sign, and consecutive added lines or removed lines gathered into one run. The last, empty line a patch's final newline leaves is not a line. Tests beside it in `flaiover/src/lib/review.test.ts`.

## Done when

- A patch with two removed lines, then three added, then context, then one added gives a removed run of two, an added run of three, the context, and an added run of one.
- No line's text begins with its sign, and a hunk header and a "no newline" note keep their whole text.
- `vitest`, `prettier`, `eslint`, and `svelte-check` pass.

## Notes
