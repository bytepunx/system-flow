---
id: T-3385
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0793
owner: alex
created: 2026-11-15T21:17:19Z
updated: 2026-11-15T21:31:17Z
transitions:
  - to: ready
    at: 2026-11-15T21:31:08Z
    by: agent-S-0793
  - to: in-progress
    at: 2026-11-15T21:31:08Z
    by: agent-S-0793
  - to: done
    at: 2026-11-15T21:31:17Z
    by: agent-S-0793
stream: S-0793
tags: []
after: [T-3384]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 944
      cache_read: 483682
      cache_write: 12355
      cost: 0.2145
---
# T-3385 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0793 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-15: ran in the main checkout and exited 0: the release policy threshold is met, with 1 accepted story not yet released, S-0741 (51.66 USD/week; count 1 of 1).
