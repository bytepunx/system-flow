---
id: T-3133
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0729
owner: alex
created: 2026-11-12T03:39:40Z
updated: 2026-11-12T03:41:32Z
transitions:
  - to: ready
    at: 2026-11-12T03:41:11Z
    by: agent-S-0729
  - to: in-progress
    at: 2026-11-12T03:41:11Z
    by: agent-S-0729
  - to: done
    at: 2026-11-12T03:41:32Z
    by: agent-S-0729
stream: S-0729
tags: []
after: [T-3132]
usage:
  source: log
  seconds: 21
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 857
      cache_read: 473711
      cache_write: 17592
      cost: 0.2526
---
# T-3133 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0729 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-12: exited 0 in the main checkout. The policy threshold is met: 1 accepted story not yet released, S-0711 (2.5 USD/week), at the count threshold of 1.
