---
id: T-3304
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0772
owner: alex
created: 2026-11-13T13:42:42Z
updated: 2026-11-13T15:19:40Z
transitions:
  - to: ready
    at: 2026-11-13T15:19:32Z
    by: agent-S-0772
  - to: in-progress
    at: 2026-11-13T15:19:32Z
    by: agent-S-0772
  - to: done
    at: 2026-11-13T15:19:40Z
    by: agent-S-0772
stream: S-0772
tags: []
after: [T-3303]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 572
      cache_read: 321436
      cache_write: 11324
      cost: 0.1663
---
# T-3304 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0772 --step evaluate --exit 2`, giving it the step's output.

## Notes

Stopped the run with exit 2:

```text
release policy threshold: not met
no accepted story is waiting for a release, so the threshold is not met
value  0 USD/week of 1000 USD/week
count  0 of 1
nothing accepted is waiting for a release
{"ts":"2026-11-13T15:19:35.889Z","level":"FATAL","msg":"command failed","component":"cmd","err":"release policy threshold not met: no accepted story is waiting for a release, so the threshold is not met"}
```
