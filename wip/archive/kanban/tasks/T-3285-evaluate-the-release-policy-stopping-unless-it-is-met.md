---
id: T-3285
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0766
owner: alex
created: 2026-11-13T11:40:39Z
updated: 2026-11-13T11:42:08Z
transitions:
  - to: ready
    at: 2026-11-13T11:41:59Z
    by: agent-S-0766
  - to: in-progress
    at: 2026-11-13T11:42:00Z
    by: agent-S-0766
  - to: done
    at: 2026-11-13T11:42:08Z
    by: agent-S-0766
stream: S-0766
tags: []
after: [T-3284]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 350
      cache_read: 162345
      cache_write: 7321
      cost: 0.0981
---
# T-3285 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0766 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran as `scripts/flai.sh release --evaluate --require-met --judged ""` in the main checkout; exit 0. Policy `threshold` met: 1 accepted story not yet released (S-0734, no value), at the threshold of 1.
