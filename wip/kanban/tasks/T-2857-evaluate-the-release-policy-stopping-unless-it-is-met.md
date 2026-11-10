---
id: T-2857
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0658
owner: alex
created: 2026-11-10T07:56:53Z
updated: 2026-11-10T08:28:58Z
transitions:
  - to: ready
    at: 2026-11-10T08:28:21Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:28:22Z
    by: system-flow
  - to: done
    at: 2026-11-10T08:28:58Z
    by: agent-S-0658
stream: S-0658
tags: []
after: [T-2856]
usage:
  source: log
  seconds: 33
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 18
      output: 166
      cache_read: 1300039
      cache_write: 6279
      cost: 0.4806
---
# T-2857 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0658 --step evaluate --exit 2`, giving it the step's output.

## Notes

Exited 0 in the main checkout, run as `scripts/flai.sh release --evaluate --require-met --judged ""`, the tree's flai, as this project's tooling convention asks. The installed flai on `PATH` (1.86.6) refused the command first with exit 1, `unknown flag: --judged`, before evaluating anything.

```text
release policy threshold: met
3 accepted stories are not yet released, and their count, 3, is at or over the threshold of 1
value  102.5 USD/week of 1000 USD/week
count  3 of 1
pending:
  S-0651  70 USD/week   flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200
  S-0653  30 USD/week   TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports
  S-0654  2.5 USD/week  A checked edit of a work item writes, checks, and restores it with no lock, so a reader in between sees an edit the check refuses
```
