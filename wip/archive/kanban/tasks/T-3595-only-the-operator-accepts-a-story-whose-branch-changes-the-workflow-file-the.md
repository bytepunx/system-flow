---
id: T-3595
type: task
nature: remediation
title: "Only the operator accepts a story whose branch changes the workflow file: the acceptance preview, the orchestrator's blockers, and the merge queue's eligibility"
status: done
parent: S-0844
owner: alex
created: 2026-11-17T12:36:28Z
updated: 2026-11-28T00:21:33Z
transitions:
  - to: ready
    at: 2026-11-27T23:38:19Z
    by: agent-S-0844
  - to: in-progress
    at: 2026-11-27T23:38:19Z
    by: agent-S-0844
  - to: done
    at: 2026-11-28T00:21:33Z
    by: agent-S-0844
stream: S-0844
tags: [flai, acceptance]
touches: [flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/cmd/accept_workflow_file_test.go, flai/cmd/accept.go, flai/internal/preview/accept_test.go]
usage:
  source: log
  seconds: 2594
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 90
      output: 22435
      cache_read: 6798626
      cache_write: 187070
      cost: 2.9187
---
# T-3595 Only the operator accepts a story whose branch changes the workflow file: the acceptance preview, the orchestrator's blockers, and the merge queue's eligibility

## Work

This task waits for no other task. It shares no file with the other tasks of the first layer.

ADR-0301 §7 says a story whose branch changes `design/workflows/board.md` is accepted by the operator alone, never by the orchestrator or the merge queue. Today only the paths Claude Code protects (ADR-0106) get that rule (I-0380, instance 1).

- In `flai/internal/preview/accept.go`, `AcceptWith`, beside the ADR-0106 block, check whether the story's changed paths include the workflow file. Its path is `workflow.Rel` of the manifest's design folder, under `workflow.File`. When they do, set `OperatorOnly` and block any `by` that is not an operator, as the protected-path block does. The reason must name the workflow file and ADR-0301 §7, not "paths Claude Code protects".
- Do not add the file to `flai/internal/protected`. That list also drives `permission_prompt`, the guard, and OpenCode's deny rules (ADR-0106), so adding the file would change who is asked before every write to it.
- In `flai/internal/preview/orchestrator.go`, add a blocker of its own, for example `workflow_file`, beside `protected_paths`. The orchestrator's dry run then names it.
- In `flai/internal/mergequeue/eligible.go`, add a kind, for example `KindWorkflowFile`, that `ComputeWith` gives a story in review whose branch changes the file. The queue then never batches such a story, instead of refusing it only at landing. Read the change the way `undoesMain` does.
- Write `flai/cmd/accept_workflow_file_test.go`, after the pattern of `accept_protected_test.go`. It must cover three cases: a story whose branch changes the workflow file is refused to `--by orchestrator`, its preview lists the new blocker, and the operator accepts it. Add a case to `flai/internal/mergequeue/eligible_test.go` for the new kind.

## Done when

- `flai accept S-nnnn --by orchestrator --dry-run` on a story whose branch changes `design/workflows/board.md` lists the new blocker and names ADR-0301 §7, and the operator's acceptance of it passes.
- The merge queue lists that story as not eligible, with the new kind.
- A story that changes neither the workflow file nor a protected path is accepted as before.
- `flai test` on the changed paths passes.

## Notes
