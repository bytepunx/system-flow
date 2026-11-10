---
id: T-2854
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0657
owner: alex
created: 2026-11-10T07:11:22Z
updated: 2026-11-10T07:55:30Z
transitions:
  - to: ready
    at: 2026-11-10T07:55:17Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T07:55:18Z
    by: system-flow
  - to: done
    at: 2026-11-10T07:55:30Z
    by: system-flow
stream: S-0657
tags: []
after: [T-2853]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 604
      cache_read: 303247
      cache_write: 13000
      cost: 0.1768
---
# T-2854 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0657 --step evaluate --exit 2`, giving it the step's output.

## Notes

Stopped the run with exit 2:

```text
release policy threshold: not met
no accepted story is waiting for a release, so the threshold is not met
value  0 USD/week of 1000 USD/week
count  0 of 1
nothing accepted is waiting for a release
{"ts":"2026-11-10T07:55:26.578Z","level":"FATAL","msg":"command failed","component":"cmd","err":"release policy threshold not met: no accepted story is waiting for a release, so the threshold is not met"}
```
