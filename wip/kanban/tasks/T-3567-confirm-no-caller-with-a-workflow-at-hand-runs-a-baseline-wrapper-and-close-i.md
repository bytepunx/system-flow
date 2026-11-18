---
id: T-3567
type: task
nature: remediation
title: Confirm no caller with a workflow at hand runs a baseline wrapper, and close I-0379
status: done
parent: S-0843
owner: alex
created: 2026-11-17T12:29:26Z
updated: 2026-11-18T10:51:58Z
transitions:
  - to: ready
    at: 2026-11-18T10:35:11Z
    by: agent-S-0843
  - to: in-progress
    at: 2026-11-18T10:35:11Z
    by: agent-S-0843
  - to: done
    at: 2026-11-18T10:51:58Z
    by: agent-S-0843
stream: S-0843
tags: [flai, workflow, issues]
touches: [design/issues/I-0379-callers-of-the-workflow-aware-guard-forecast-replay-and-share-hold-still-pass.md, design/issues/summary.md]
after: [T-3563, T-3565, T-3609, T-3610]
usage:
  source: log
  seconds: 1007
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 74
      output: 17523
      cache_read: 5115009
      cache_write: 231604
      cost: 2.7582
---
# T-3567 Confirm no caller with a workflow at hand runs a baseline wrapper, and close I-0379

## Work

Search the non-test Go under `flai/` again for each baseline wrapper, the functions whose comment says "under the baseline workflow" (`planning.Forecast`, `Replay`, `NewPuller`, `PullTimes`, `CostOfDelay`, `workitem.NewHolds`, `HoldsMemo.NewHolds`, `PlanOf`, `PlanCandidatesOf`, `itemedit.Covering`, `preview.ThreadItems`, `inbox.Describe`, `StarvedFrom`). Each caller that has a repository or a workflow at hand and whose result depends on a state name must call the `In` form. A caller of `NewHolds(nil, ...).Claim(&workitem.Item{Touches: ...})`, which reads no task, is the same under every workflow and stays.

The four callers the instance names were settled before this story: `flai/cmd/guard.go` gives the guard `repo.Flow()` for the orchestrator and every role a role file holds, `flai/internal/serve/replan.go` calls `planning.ReplayIn`, and `flai/internal/workitem/hold.go` names a holding story's state by the workflow's labels. Check each again and say so in the issue's reason. The instance's last sentence, `digest.Decisions` filtering on `Item.CompletedAt`, which reads `Item.Closed`, is I-0377's cause and S-0842's to fix: name S-0842 in the reason.

Waits for T-3563 and T-3565, whose changes the search must see.

## Done when

- The search finds no remaining caller that should take the workflow, or this story's agent adds a task for one it finds.
- `flai issue close I-0379 --reason "<what fixed it>"` has run, naming T-3561, T-3563, T-3565, the callers settled before, and S-0842 for the digest.
- `flai check --strict` is clean.

## Notes
