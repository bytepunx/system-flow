---
id: T-3267
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0762
owner: alex
created: 2026-11-13T07:24:38Z
updated: 2026-11-13T07:26:59Z
transitions:
  - to: ready
    at: 2026-11-13T07:26:31Z
    by: agent-S-0762
  - to: in-progress
    at: 2026-11-13T07:26:31Z
    by: agent-S-0762
  - to: done
    at: 2026-11-13T07:26:59Z
    by: agent-S-0762
stream: S-0762
tags: []
after: [T-3266]
usage:
  source: log
  seconds: 28
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 371
      cache_read: 164165
      cache_write: 7828
      cost: 0.1029
---
# T-3267 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0762 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: `scripts/flai.sh release --pending --require-pending --json` ran in the main checkout and exited 0, pushed true: flai 1.86.22 → 1.86.23 (tag `flai/v1.86.23`, patch: S-0683, S-0723, S-0726, S-0732); template 1.1.34 → 1.1.35 (patch: S-0726), published to the template repository as v1.1.35 at 2c78d45.
