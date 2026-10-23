---
id: T-2042
type: task
nature: improvement
title: "A thread the orchestrator opened that nobody else wrote on holds the story until the orchestrator resolves it: the acceptance preview names it as the orchestrator's own to resolve"
status: done
parent: S-0467
owner: alex
created: 2026-10-22T20:59:46Z
updated: 2026-10-23T05:51:47Z
transitions:
  - to: ready
    at: 2026-10-23T04:51:14Z
    by: agent-S-0467
  - to: in-progress
    at: 2026-10-23T04:51:15Z
    by: agent-S-0467
  - to: done
    at: 2026-10-23T05:20:17Z
    by: agent-S-0467
stream: S-0467
tags: [flai]
touches: [flai/internal/preview/orchestrator.go, flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/cmd/accept_own_thread_test.go]
after: [T-2037]
usage:
  source: log
  seconds: 1742
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 11620
      cache_read: 2059053
      cache_write: 101262
      cost: 1.2035
---
# T-2042 A thread the orchestrator opened that nobody else wrote on holds the story until the orchestrator resolves it: the acceptance preview names it as the orchestrator's own to resolve

## Work

Waits for the ADR task, which settles the decision this builds: remedy 3 in the orchestrator's stricter form (ADR-0189, TH-0598), not the planner's.

Today `orchestratorBlockers` in `flai/internal/preview/orchestrator.go` blocks on every thread on the story or one of its tasks that `th.Open()` says is not resolved, whoever opened it, with one message. The orchestrator may already resolve a thread it opened (`flai guard`), but nothing tells it that a thread holding the story is its own to resolve, so a concern it raised became a gate the operator had to clear (TH-0387, S-0341, four days).

- Add to `flai/internal/threads/threads.go` one method that says whether a thread is the orchestrator's own: opened by the orchestrator (`workitem.IsOrchestrator`) with every entry the orchestrator's. One entry by anyone else, the operator or an agent, makes it a thread that holds until resolved by whoever may.
- `orchestratorBlockers` lists such a thread under its own blocker code, `own_thread_open`, whose message names the thread and says the orchestrator resolves it with a reason (`thread_resolve`, or `flai thread resolve`) once its concern is gone, or leaves the story in review. Every other open thread blocks as `thread_open`, as today.
- The merge queue's eligibility keeps both as holding the story; its reason may say the same of the orchestrator's own thread. Leave `eligible.go` unchanged if it reads well as is.

Tests: in `flai/internal/threads/threads_test.go` the method on a thread the orchestrator alone wrote, one the operator replied on, and one a story's agent opened; in a new `flai/cmd/accept_own_thread_test.go` the orchestrator's dry run lists its own untouched thread as `own_thread_open`, a thread the operator replied on as `thread_open`, and once the orchestrator resolves its own thread its acceptance goes through.

## Done when

- The orchestrator's acceptance preview names a thread it opened that nobody else wrote on as `own_thread_open`, its own to resolve, and still blocks on every open thread.
- Once the orchestrator resolves it, the acceptance goes through, with tests.
- `flai test` on the changed paths passes.

## Notes

Layer 2, beside the dirty-main task; no path in common with it. Rewritten by agent-S-0467 for the stricter form. The prompt and the convention are changed in the later prompt task, so that the orchestrator accepts on flai's blockers and the criteria, opens a thread on a story in review only for a reason not to accept it, and resolves it once the reason is gone.
