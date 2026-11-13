---
id: T-3308
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0773
owner: alex
created: 2026-11-13T15:20:30Z
updated: 2026-11-13T15:23:08Z
transitions:
  - to: ready
    at: 2026-11-13T15:22:56Z
    by: agent-S-0773
  - to: in-progress
    at: 2026-11-13T15:22:56Z
    by: agent-S-0773
  - to: done
    at: 2026-11-13T15:23:08Z
    by: agent-S-0773
stream: S-0773
tags: []
after: [T-3307]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 102
      cache_read: 498010
      cache_write: 1873
      cost: 0.1871
---
# T-3308 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0773 --step evaluate --exit 2`, giving it the step's output.

## Notes

Exited 0 in the main checkout:

```text
release policy threshold: met
1 accepted story is not yet released, and their count, 1, is at or over the threshold of 1
value  30 USD/week of 1000 USD/week
count  1 of 1
pending:
  S-0662  30 USD/week  The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first
```
