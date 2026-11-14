---
id: T-3359
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0786
owner: alex
created: 2026-11-14T17:32:28Z
updated: 2026-11-14T17:33:59Z
transitions:
  - to: ready
    at: 2026-11-14T17:33:47Z
    by: agent-S-0786
  - to: in-progress
    at: 2026-11-14T17:33:47Z
    by: agent-S-0786
  - to: done
    at: 2026-11-14T17:33:59Z
    by: agent-S-0786
stream: S-0786
tags: []
after: [T-3358]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 51
      cache_read: 500959
      cache_write: 2333
      cost: 0.1893
---
# T-3359 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0786 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-14: the command ran in the main checkout and exited 0. The release policy threshold was met: 1 accepted story was not yet released (S-0785, 12.5 USD/week), against a count threshold of 1.
