---
id: T-3363
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0787
owner: alex
created: 2026-11-14T18:00:18Z
updated: 2026-11-14T18:02:28Z
transitions:
  - to: ready
    at: 2026-11-14T18:02:05Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T18:02:06Z
    by: system-flow
  - to: done
    at: 2026-11-14T18:02:28Z
    by: agent-S-0787
stream: S-0787
tags: []
after: [T-3362]
usage:
  source: log
  seconds: 22
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 639
      cache_read: 337019
      cache_write: 16732
      cost: 0.2141
---
# T-3363 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0787 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-14: `scripts/flai.sh release --pending --require-pending --json` in the main checkout exited 0: S-0781 released as `flai/v1.95.6` (patch from 1.95.5) and template 1.1.41 (patch from 1.1.40, published to the template repository at 5220cba); tags and branch pushed.
- A first attempt redirected its output into a folder that did not exist yet, so the shell refused the redirect and the command never ran; the second run above is the step's run.
