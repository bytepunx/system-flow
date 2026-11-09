---
id: T-2739
type: task
nature: remediation
title: I-0290 is closed, naming the ADR and the change that keeps the close-outs past a damaged journal line
status: done
parent: S-0624
owner: alex
created: 2026-11-09T17:04:52Z
updated: 2026-11-09T17:17:48Z
transitions:
  - to: ready
    at: 2026-11-09T17:17:39Z
    by: agent-S-0624
  - to: in-progress
    at: 2026-11-09T17:17:41Z
    by: agent-S-0624
  - to: done
    at: 2026-11-09T17:17:48Z
    by: agent-S-0624
stream: S-0624
tags: [issues]
touches: [design/issues/I-0290-one-line-holding-a-nul-byte-in-flai-serve-s-journal-drops-the-whole-close-outs.md, design/issues/summary.md]
after: [T-2738]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 653
      cache_read: 563592
      cache_write: 9105
      cost: 0.1986
---
# T-2739 I-0290 is closed, naming the ADR and the change that keeps the close-outs past a damaged journal line

## Work

Close I-0290 with `flai issue close I-0290 --reason "<reason>"`, the reason naming the ADR T-2737 wrote and the reader T-2738 changed: a serve journal line that holds NUL bytes or is not an entry is skipped and warned of, and the close-outs section stays. flai regenerates `design/issues/summary.md`.

It waits for T-2738 because the issue is closed only once the fix and its test pass.

## Done when

- I-0290's status is `closed`, with the reason recorded.
- `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
