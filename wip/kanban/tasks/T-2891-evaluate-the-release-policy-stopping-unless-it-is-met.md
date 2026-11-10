---
id: T-2891
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0664
owner: alex
created: 2026-11-10T08:33:20Z
updated: 2026-11-10T08:34:56Z
transitions:
  - to: ready
    at: 2026-11-10T08:34:32Z
    by: agent-S-0664
  - to: in-progress
    at: 2026-11-10T08:34:32Z
    by: agent-S-0664
  - to: done
    at: 2026-11-10T08:34:56Z
    by: agent-S-0664
stream: S-0664
tags: []
after: [T-2890]
usage:
  source: log
  seconds: 24
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 2
      output: 8
      cache_read: 162697
      cache_write: 1018
      cost: 0.0602
---
# T-2891 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0664 --step evaluate --exit 2`, giving it the step's output.

## Notes

Exited 0 in the main checkout, run as `scripts/flai.sh release --evaluate --require-met --judged ""`, the tree's flai, as S-0658 did: the installed flai on `PATH` (1.86.6) lacks `--judged`.

```text
release policy threshold: met
2 accepted stories are not yet released, and their count, 2, is at or over the threshold of 1
value  202.5 USD/week of 1000 USD/week
count  2 of 1
pending:
  S-0637  175 USD/week   Stories that each add a section at one spot in a shared design document conflict at every sync and acceptance, and the acceptance sends them back
  S-0652  27.5 USD/week  TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load
```
