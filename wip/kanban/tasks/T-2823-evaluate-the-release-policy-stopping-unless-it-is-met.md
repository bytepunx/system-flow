---
id: T-2823
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0646
owner: alex
created: 2026-11-09T22:22:55Z
updated: 2026-11-10T06:18:26Z
transitions:
  - to: ready
    at: 2026-11-10T06:18:14Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T06:18:15Z
    by: system-flow
  - to: done
    at: 2026-11-10T06:18:26Z
    by: system-flow
stream: S-0646
tags: []
after: [T-2822]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 556
      cache_read: 304242
      cache_write: 14433
      cost: 0.1874
---
# T-2823 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0646 --step evaluate --exit 2`, giving it the step's output.

## Notes

Stopped the run with exit 2:

```text
release policy threshold: not met
no accepted story is waiting for a release, so the threshold is not met
value  0 USD/week of 1000 USD/week
count  0 of 1
nothing accepted is waiting for a release
{"ts":"2026-11-10T06:18:19.839Z","level":"FATAL","msg":"command failed","component":"cmd","err":"release policy threshold not met: no accepted story is waiting for a release, so the threshold is not met"}
```
