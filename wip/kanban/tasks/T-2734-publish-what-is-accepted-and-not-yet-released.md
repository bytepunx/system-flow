---
id: T-2734
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0641
owner: alex
created: 2026-11-09T16:56:49Z
updated: 2026-11-09T16:58:41Z
transitions:
  - to: ready
    at: 2026-11-09T16:58:00Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T16:58:01Z
    by: system-flow
  - to: done
    at: 2026-11-09T16:58:41Z
    by: agent-S-0641
stream: S-0641
tags: []
after: [T-2733]
usage:
  source: log
  seconds: 40
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 2772
      cache_read: 1555624
      cache_write: 44380
      cost: 0.7217
---
# T-2734 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes

- 2026-11-09: `flai release --pending --json` ran in the main checkout and exited 0, answering `{"plans": null, "pushed": false, "reason": "nothing pending", "tags": null}`.
- S-0610 and S-0615 had been released seconds before by a publish this step did not make: commit 7029b028a, "chore: publish flai 1.86.5, template 1.1.28", at 16:57:57Z under the operator's name, with `flai/v1.86.5` on `origin` and `origin/main` at 7029b028a. `flai release --evaluate` afterwards answers nothing accepted is waiting for a release. Recorded as I-0320.
