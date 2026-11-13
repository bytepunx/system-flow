---
id: T-3301
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0771
owner: alex
created: 2026-11-13T13:41:40Z
updated: 2026-11-13T14:36:49Z
transitions:
  - to: ready
    at: 2026-11-13T14:36:22Z
    by: agent-S-0771
  - to: in-progress
    at: 2026-11-13T14:36:22Z
    by: agent-S-0771
  - to: done
    at: 2026-11-13T14:36:49Z
    by: agent-S-0771
stream: S-0771
tags: []
after: [T-3300]
usage:
  source: log
  seconds: 27
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 749
      cache_read: 326410
      cache_write: 14205
      cost: 0.1939
---
# T-3301 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0771 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran as `scripts/flai.sh release --evaluate --require-met --judged ""` (the manifest's `flai.command`) in the main checkout; exited 0. Output:

```text
release policy threshold: met
1 accepted story is not yet released, and their count, 1, is at or over the threshold of 1
value  986.84 USD/week of 1000 USD/week
count  1 of 1
pending:
  S-0661  986.84 USD/week  A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables
```
