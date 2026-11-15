---
id: T-3367
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0789
owner: alex
created: 2026-11-14T23:49:29Z
updated: 2026-11-15T01:58:38Z
transitions:
  - to: ready
    at: 2026-11-15T01:58:19Z
    by: agent-S-0789
  - to: in-progress
    at: 2026-11-15T01:58:19Z
    by: agent-S-0789
  - to: done
    at: 2026-11-15T01:58:38Z
    by: agent-S-0789
stream: S-0789
tags: []
after: [T-3366]
usage:
  source: log
  seconds: 19
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 1179
      cache_read: 656942
      cache_write: 26125
      cost: 0.364
---
# T-3367 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0789 --step evaluate --exit 2`, giving it the step's output.

## Notes

Stopped the run with exit 2:

```text
release policy threshold: not met
no accepted story is waiting for a release, so the threshold is not met
value  0 USD/week of 1000 USD/week
count  0 of 1
nothing accepted is waiting for a release
{"ts":"2026-11-15T01:58:23.048Z","level":"FATAL","msg":"command failed","component":"cmd","err":"release policy threshold not met: no accepted story is waiting for a release, so the threshold is not met"}
```
