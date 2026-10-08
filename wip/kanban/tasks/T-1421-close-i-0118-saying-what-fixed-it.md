---
id: T-1421
type: task
nature: remediation
title: Close I-0118 saying what fixed it
status: done
parent: S-0346
owner: alex
created: 2026-10-08T08:54:11Z
updated: 2026-10-08T09:36:22Z
transitions:
  - to: ready
    at: 2026-10-08T09:30:52Z
    by: agent-S-0346
  - to: in-progress
    at: 2026-10-08T09:35:55Z
    by: agent-S-0346
  - to: done
    at: 2026-10-08T09:36:22Z
    by: agent-S-0346
stream: S-0346
tags: [issues]
touches: [design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1420]
usage:
  source: log
  seconds: 27
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1114
      cache_read: 356162
      cache_write: 11941
      cost: 0.1891
---
# T-1421 Close I-0118 saying what fixed it

## Work

- Run `flai check --strict` in the main checkout and confirm that it finds no `markdown.MD034` under `wip/`.
- Close the issue with `flai issue close I-0118 --reason`. The reason says that all five instances were lines an installed flai older than 1.39.4 wrote, before S-0324 gave its MD034 the bare `www.` rule. It says that every agent's write to `wip/` is now refused on one. It names T-1420's fix for the one write flai makes from text no agent checks, and the test that reproduces it.
- Waits for T-1420, whose fix and test the reason names.

## Done when

- I-0118 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean.

## Notes
