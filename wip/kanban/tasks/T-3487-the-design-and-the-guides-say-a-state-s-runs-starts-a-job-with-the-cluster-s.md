---
id: T-3487
type: task
nature: feature
title: The design and the guides say a state's runs starts a job with the cluster's steps, which are skipped where no cluster can run them, and kluster-verified is checked where a transition names it
status: done
parent: S-0810
owner: alex
created: 2026-11-16T23:51:46Z
updated: 2026-11-18T11:11:37Z
transitions:
  - to: ready
    at: 2026-11-18T11:02:41Z
    by: agent-S-0810
  - to: in-progress
    at: 2026-11-18T11:02:41Z
    by: agent-S-0810
  - to: done
    at: 2026-11-18T11:11:37Z
    by: agent-S-0810
stream: S-0810
tags: [flai]
touches: [design/system/kluster.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/change-the-board.md, flai/internal/workflow/baseline/board.md, template/root/design/workflows/board.md, template/CHANGELOG.md]
after: [T-3475, T-3484]
usage:
  source: log
  seconds: 536
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 125
      output: 62416
      cache_read: 12252040
      cache_write: 426035
      cost: 6.0545
---
# T-3487 The design and the guides say a state's runs starts a job with the cluster's steps, which are skipped where no cluster can run them, and kluster-verified is checked where a transition names it

## Work

Waits for T-3475 and T-3484, whose behaviour it describes: the rule in the move, and the state job with its `kluster:` steps and skips.

1. `design/system/kluster.md`: `## What happens when` says the deploy and checks run as steps of the job a state's `runs:` starts, on `review` on the baseline, and are skipped with the reason where no deploy, host action, or namespace allows them; `## A failed deploy holds the story` says the story is sent back by the state's send-back transition.
2. `design/system/workflow.md`: under the workflow file, what a state's `runs:` starts, the job's steps, its pass and its send-back; under `## Transitions and who makes them`, that `kluster-verified` is evaluated by the move where a transition names it, with its severity.
3. `design/system/flai-cli.md`: under `### The close-out job`, the state job beside the close-out and resolve jobs, its record, and its log.
4. `docs/users/flai.md` (`## Run a project's kluster`, `### When a story's deploy or its checks fail`), `docs/operators/index.md` (`### The kluster host action (S-0577)`), and `docs/operators/runbooks/change-the-board.md`: a story entering a state whose `runs:` names `deploy` is deployed there, and the runbook's example names `runs:` on a state of its own.

## Done when

- `flai test` passes on the changed paths, the markdown lint among them.
- No document says the cluster deploys on entering review, or on any state of kind `waiting`, in code.

## Notes
