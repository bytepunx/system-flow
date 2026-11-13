---
id: T-3298
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0770
owner: alex
created: 2026-11-13T13:16:45Z
updated: 2026-11-13T13:18:31Z
transitions:
  - to: ready
    at: 2026-11-13T13:18:19Z
    by: agent-S-0770
  - to: in-progress
    at: 2026-11-13T13:18:19Z
    by: agent-S-0770
  - to: done
    at: 2026-11-13T13:18:31Z
    by: agent-S-0770
stream: S-0770
tags: []
after: [T-3297]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 372
      cache_read: 160331
      cache_write: 7485
      cost: 0.0994
---
# T-3298 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0770 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran as `scripts/flai.sh release --evaluate --require-met --judged ""` (the manifest's `flai.command`) in the main checkout; exited 0. The threshold policy is met: 1 accepted story not yet released (S-0740, no value), count 1 of 1, value 0 of 1000 USD/week.
