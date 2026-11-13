---
id: T-3341
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0783
owner: alex
created: 2026-11-13T22:47:23Z
updated: 2026-11-13T22:49:25Z
transitions:
  - to: ready
    at: 2026-11-13T22:49:15Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T22:49:15Z
    by: system-flow
  - to: done
    at: 2026-11-13T22:49:25Z
    by: system-flow
stream: S-0783
tags: []
after: [T-3340]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 660
      cache_read: 319855
      cache_write: 15956
      cost: 0.2048
---
# T-3341 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0783 --step evaluate --exit 2`, giving it the step's output.

## Notes

Stopped the run with exit 2:

```text
release policy threshold: not met
no accepted story is waiting for a release, so the threshold is not met
value  0 USD/week of 1000 USD/week
count  0 of 1
nothing accepted is waiting for a release
{"ts":"2026-11-13T22:49:19.545Z","level":"FATAL","msg":"command failed","component":"cmd","err":"release policy threshold not met: no accepted story is waiting for a release, so the threshold is not met"}
```
