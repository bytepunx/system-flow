---
id: T-2077
type: task
nature: remediation
title: I-0182 is closed with the reason naming the scoped check a checked edit now applies
status: done
parent: S-0486
owner: alex
created: 2026-10-23T17:51:20Z
updated: 2026-10-24T00:18:25Z
transitions:
  - to: ready
    at: 2026-10-24T00:18:10Z
    by: agent-S-0486
  - to: in-progress
    at: 2026-10-24T00:18:11Z
    by: agent-S-0486
  - to: done
    at: 2026-10-24T00:18:25Z
    by: agent-S-0486
stream: S-0486
tags: [flai]
touches: [design/issues/I-0182-criteria-tick-is-refused-on-another-open-story-s-narrative-state-finding-which-the-story-s-own-scoped-check-leaves-out.md, design/issues/summary.md]
after: [T-2075, T-2076]
usage:
  source: log
  seconds: 14
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1813
      cache_read: 355774
      cache_write: 11405
      cost: 0.1795
---
# T-2077 I-0182 is closed with the reason naming the scoped check a checked edit now applies

## Work

In the story's worktree, once T-2075 is done and T-2076 describes it, close the issue:

```bash
flai issue close I-0182 --reason "Fixed by S-0486: itemedit.Apply scopes the check a checked edit runs to the item's story (ADR-NNNN), so a finding another writer brings between the two runs, a narrative.state on another open story's narrative or a markdown finding on another story's file, warns or is left out and no longer refuses a tick or an edit"
```

with the ADR's number from T-2074. `flai issue close` rewrites the issue file and `design/issues/summary.md`; commit both with `flai task done`. Then tick the story's second criterion with `flai criteria tick S-0486 2`, which is itself a tick under the new rule, and the first once T-2075's tests pass on the branch.

Waits for T-2075, which removes the cause, and T-2076, so the close happens after the words that describe the fix.

## Done when

`flai issue list` shows I-0182 closed with the reason, `design/issues/summary.md` no longer lists it, both are committed on the story branch, and the story's two criteria are ticked.

## Notes

Issue commands write under the current checkout: run `flai issue close` in the story's worktree, not in the main checkout.
