---
id: T-3475
type: task
nature: feature
title: The move evaluates the kluster-verified rule where a transition names it, in place of a warning on every move into a waiting state
status: done
parent: S-0810
owner: alex
created: 2026-11-16T23:50:52Z
updated: 2026-11-18T09:41:06Z
transitions:
  - to: ready
    at: 2026-11-18T09:21:13Z
    by: agent-S-0810
  - to: in-progress
    at: 2026-11-18T09:21:13Z
    by: agent-S-0810
  - to: done
    at: 2026-11-18T09:41:06Z
    by: agent-S-0810
stream: S-0810
tags: [flai]
touches: [flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go, flai/internal/workitem/transition.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/cmd/move.go, flai/cmd/move_workflow_test.go]
usage:
  source: log
  seconds: 1193
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 29
      output: 13329
      cache_read: 3245129
      cache_write: 103919
      cost: 1.5352
---
# T-3475 The move evaluates the kluster-verified rule where a transition names it, in place of a warning on every move into a waiting state

## Work

Waits for no task: it changes only the move and its callers, which no other task of this story touches.

Today `ruleChecks[RuleKlusterVerified]` in `flai/internal/workitem/rules.go` returns nil, and `mcpserver.ReviewWarnings` (`flai/internal/mcpserver/server.go`, called from `item_move` and `flai/cmd/move.go` after `TransitionAll`) warns of the namespace on every move into a state of kind `waiting`, whether or not the transition names the rule.

1. Give `Repo.Move` a way to evaluate `kluster-verified` without `workitem` importing `kluster`: a checker the caller passes with the move (an option on the move in `flai/internal/workitem/transition.go`), which answers whether the story's namespace was deployed, and its checks passed while the `checks` host action is on, at the head of its branch.
2. Evaluate it only on a transition whose `rules` name `kluster-verified`, with the rule's severity: `warn` on the baseline, `refuse` where the file says so. With the kluster host action off, no `deploy` configured, or `kluster.agent_namespaces` 0, the rule passes with nothing to say, as ADR-0312 §5 asks.
3. `flai move` and `item_move` pass the checker built from `kluster.AtHead`; `ReviewWarnings` no longer warns of the namespace on its own, so the warning comes from the rule alone.
4. Tests: the baseline's move to review warns as today; a workflow whose transition into another state names the rule at `refuse` refuses an undeployed story; a transition that does not name it says nothing of the namespace; no cluster configured passes.

## Done when

- `flai test` passes on the changed paths.
- The baseline's move to review gives the same warning text as before this task, and a test shows a transition that does not name `kluster-verified` gives none.

## Notes
