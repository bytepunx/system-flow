---
id: T-3082
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0718
owner: alex
created: 2026-11-11T01:44:41Z
updated: 2026-11-11T22:15:52Z
transitions:
  - to: ready
    at: 2026-11-11T22:15:21Z
    by: agent-S-0718
  - to: in-progress
    at: 2026-11-11T22:15:21Z
    by: agent-S-0718
  - to: done
    at: 2026-11-11T22:15:52Z
    by: agent-S-0718
stream: S-0718
tags: []
after: [T-3081]
usage:
  source: log
  seconds: 31
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1000
      cache_read: 496795
      cache_write: 19982
      cost: 0.2792
---
# T-3082 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0718 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-11: ran in the main checkout with flai 1.86.17 and exited 0. It released flai 1.86.17 → 1.86.18 as a patch, tagged `flai/v1.86.18`, and pushed the branch and the tag (`pushed: true`). The release covers S-0689, S-0707, S-0714, and S-0715.
