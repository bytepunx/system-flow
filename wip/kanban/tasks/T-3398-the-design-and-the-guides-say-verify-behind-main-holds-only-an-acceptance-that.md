---
id: T-3398
type: task
nature: remediation
title: The design and the guides say verify_behind_main holds only an acceptance that lands by hand, and I-0381 is closed
status: done
parent: S-0796
owner: alex
created: 2026-11-16T00:51:19Z
updated: 2026-11-16T01:31:00Z
transitions:
  - to: ready
    at: 2026-11-16T01:28:27Z
    by: agent-S-0796
  - to: in-progress
    at: 2026-11-16T01:28:27Z
    by: agent-S-0796
  - to: done
    at: 2026-11-16T01:30:59Z
    by: agent-S-0796
stream: S-0796
tags: [docs, acceptance, merge-queue]
touches: [design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md]
after: [T-3394, T-3396, T-3397]
usage:
  source: log
  seconds: 146
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 40
      output: 239
      cache_read: 1421019
      cache_write: 61373
      cost: 0.5616
---
# T-3398 The design and the guides say verify_behind_main holds only an acceptance that lands by hand, and I-0381 is closed

## Work

Say in each place that states ADR-0291's gate that, while the merge queue batches and is not held, the operator's preview and acceptance leave `verify_behind_main` to the batch and report the stale tiers, and that a held queue, `--without-queue-tiers`, `--without-kluster`, `--keep-removal`, and the orchestrator's acceptance keep it, linking the new ADR:

- `design/system/workflow.md` § What waits on an acceptance (S-0296), line 189.
- `design/system/flai-cli.md` § Commands (line 86) and § The acceptance gate (line 335).
- `docs/users/flai.md` § Accept and release (lines 3389 to 3399), with the `verify_behind_main_why` field beside `stale_tiers`.
- `docs/users/flaiover.md`'s Rebase and resolve paragraph (line 222), where "a verify that no longer stands" holds Accept: say it does not while the queue batches, and Accept enqueues.

Then close the issue: `flai issue close I-0381 --reason` naming the fix, the preview's predicate and the hold its two callers pass, and the tests that reproduce it.

It waits for T-3394, whose ADR it links, and for T-3396 and T-3397, whose behaviour it describes and whose tests the close reason names.

## Done when

- The four documents say it, each linking the new ADR, and the markdown lint passes on them.
- I-0381 is closed with its reason, and `design/issues/summary.md` no longer lists it.

## Notes
