---
id: T-2548
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0585
owner: alex
created: 2026-11-07T08:03:26Z
updated: 2026-11-07T08:15:59Z
transitions:
  - to: ready
    at: 2026-11-07T08:06:16Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:06:17Z
    by: system-flow
  - to: done
    at: 2026-11-07T08:15:59Z
    by: agent-S-0585
blocked:
  - from: 2026-11-07T08:06:26Z
    until: 2026-11-07T08:15:52Z
    reason: "step evaluate: flai release --evaluate --require-met exited 2 in the main checkout: release policy judgement not met, nothing accepted is waiting for a release"
stream: S-0585
tags: []
after: [T-2546]
usage:
  source: log
  seconds: 75
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 3585
      cache_read: 1586618
      cache_write: 49692
      cost: 0.7866
---
# T-2548 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes

- 2026-11-07: `flai release --evaluate --require-met` (run as `scripts/flai.sh`) in the main checkout exited 2. The release policy is `judgement`, which is never met by itself, and nothing accepted is waiting for a release. T-2549 was not run. Its whole output:

```text
release policy judgement: not met
the release policy is judgement, which is never met by itself: whether to release is the orchestrator's or the operator's call, with no accepted story waiting for a release
value  0 USD/week
count  0
nothing accepted is waiting for a release
{"ts":"2026-11-07T08:06:20.134Z","level":"FATAL","msg":"command failed","component":"cmd","err":"release policy judgement not met: the release policy is judgement, which is never met by itself: whether to release is the orchestrator's or the operator's call, with no accepted story waiting for a release"}
```
