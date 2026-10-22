---
id: T-2042
type: task
nature: improvement
title: "A thread the orchestrator opened with no entry by anyone else does not gate the story: the acceptance preview notes it, the merge queue's eligibility leaves it out, and the acceptance resolves it"
status: backlog
parent: S-0467
owner: alex
created: 2026-10-22T20:59:46Z
updated: 2026-10-22T21:00:06Z
transitions: []
stream: S-0467
tags: [flai]
touches: [flai/internal/preview/orchestrator.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/cmd/accept_own_thread_test.go]
after: [T-2037]
---
# T-2042 A thread the orchestrator opened with no entry by anyone else does not gate the story: the acceptance preview notes it, the merge queue's eligibility leaves it out, and the acceptance resolves it

## Work

Waits for the ADR task, which settles the decision this builds.

Today `orchestratorBlockers` in `flai/internal/preview/orchestrator.go` blocks on every thread on the story or one of its tasks that `th.Open()` says is not resolved, whoever opened it, and `openThreads` in `flai/internal/mergequeue/eligible.go` makes the same thread the `thread_open` reason. So the thread the orchestrator opens at step 4 of its acceptance, saying what is missing, holds the story on its next pass until the orchestrator or the operator resolves it, and a concern it raised beyond the criteria became a gate the operator had to clear (TH-0387, S-0341, four days).

- Add to `flai/internal/threads/threads.go` one method that says whether a thread is the orchestrator's own gate: opened by the orchestrator (`Opener()` is the orchestrator, as `workitem.IsOrchestrator` tells it) with every entry the orchestrator's. One entry by anyone else, the operator or an agent, makes it a thread that holds.
- `orchestratorBlockers` leaves such a thread out of the `BlockThreadOpen` blockers and lists it in the preview as a note, naming the thread and that the acceptance resolves it, so the orchestrator's dry run says why it does not hold. Every other open thread blocks as today.
- `openThreads` in the merge queue's eligibility leaves the same threads out, so a story the orchestrator enqueues is not kept out of the batch by the thread it opened earlier.
- Check that the acceptance's archive resolves the orchestrator's thread with ADR-0109's entry as it resolves any open thread; add nothing if it does.

Tests: in `flai/internal/threads/threads_test.go` the method on a thread the orchestrator alone wrote, one the operator replied on, and one a story's agent opened; in `flai/internal/mergequeue/eligible_test.go` a story with only the orchestrator's thread is eligible and one the operator wrote on is not; in a new `flai/cmd/accept_own_thread_test.go` the orchestrator's dry run lists its own untouched thread as a note and no blocker, its acceptance goes through and resolves the thread, and a thread the operator replied on still blocks with the text `TestTheOrchestratorDoesNotAcceptAStoryWithAThreadOpen` asserts today.

## Done when

- The orchestrator's acceptance preview and the merge queue's eligibility leave out a thread the orchestrator opened that nobody else wrote on, and still block on every other open thread.
- The orchestrator's acceptance of such a story goes through and the thread is resolved at the archive, with tests at the three levels above.
- `flai test` on the changed paths passes.

## Notes

Layer 2, beside the dirty-main task; no path in common with it. The prompt and the convention are changed in the later prompt task, so that the orchestrator accepts on flai's blockers and the criteria alone and raises a concern beyond them on a thread that does not hold the story.
