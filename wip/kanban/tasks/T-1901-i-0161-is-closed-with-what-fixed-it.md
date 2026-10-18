---
id: T-1901
type: task
nature: remediation
title: I-0161 is closed with what fixed it
status: done
parent: S-0446
owner: alex
created: 2026-10-18T02:38:37Z
updated: 2026-10-18T03:24:52Z
transitions:
  - to: ready
    at: 2026-10-18T03:24:36Z
    by: agent-S-0446
  - to: in-progress
    at: 2026-10-18T03:24:36Z
    by: agent-S-0446
  - to: done
    at: 2026-10-18T03:24:52Z
    by: agent-S-0446
stream: S-0446
tags: [issues]
touches: [design/issues/I-0161-flai-check-finds-components-examples-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1900]
usage:
  source: log
  seconds: 16
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1940
      cache_read: 364733
      cache_write: 16725
      cost: 0.213
---
# T-1901 I-0161 is closed with what fixed it

## Work

Close I-0161 in the story's worktree with `flai issue close I-0161 --reason "<reason>"`. The reason names the fix:

- A sub-project without svelte now gets one advisory `components.examples` note at its examples folder, not a warning for each file (T-1899).
- A story's scoped check leaves that note out when the story claims nothing under that folder (T-1900).
- It names the ADR that records both, and the test that reproduces the issue.

It waits for T-1900: the issue is closed only once its cause no longer occurs.

## Done when

I-0161 has status closed with that reason, `design/issues/summary.md` no longer lists it as open, and `flai check --strict` is clean.

## Notes

Run it in the worktree, not the main checkout, since `flai issue` writes under the checkout it runs in.
