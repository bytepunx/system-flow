---
id: MS-0679
title: "S-0844 holds S-0824 on overlap alone: may S-0824 start beside it?"
from: S-0824
to: S-0844
about: [flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
status: closed
participants: [flai, agent-S-0844, orchestrator]
created: 2026-11-28T00:21:40Z
updated: 2026-11-28T00:43:34Z
---

# MS-0679 S-0844 holds S-0824 on overlap alone: may S-0824 start beside it?

Between S-0824 and S-0844, about `flai/internal/preview/accept.go`, `flai/internal/preview/accept_test.go`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`.

## Entries

### 2026-11-28T00:21:40Z flai S-0824
S-0844 holds S-0824 on overlap alone: may S-0824 start beside it?

S-0824, "A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted", is ready, and its claim overlaps S-0844's on `flai/internal/preview/accept.go`, `flai/internal/preview/accept_test.go`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0398](../../../design/issues/I-0398-a-rebase-and-resolve-job-runs-the-whole-close-out-again-each-time-main-moves.md), "A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted". The issue recommends this solution:
>
> Decide, in an ADR refining ADR-0313 and ADR-0156, what a rebase and resolve job owes when main moves under it, and build it. Options to weigh: the job passes at the base its last round ran on and leaves the rest to the acceptance's own `verify_behind_main` gate (ADR-0291, ADR-0303), which sends back only the tiers the new paths overlap; the later rounds apply ADR-0303's overlap rule in place of every tier main's paths select; or flai serve holds acceptances and resolve jobs in one order, so that main does not move under a job it is about to accept. ADR-0312's acceptance suites (S-0808) bear on the same cost. The job's record and the review page say which round the job is on and why it started another.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0844 --remove <path>`, or on the task that names the path, if S-0844 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0679 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0824 can start.
3. Reply with `message_reply` saying why the hold stands: S-0824 waits until S-0844 moves to review, is cancelled, or is sent back.

### 2026-11-28T00:23:39Z agent-S-0844 S-0844 (settled)
The hold stands for a short while. In preview/accept.go and accept_test.go, S-0844 adds the workflow-file operator-only check, ChangesWorkflowFile, and one row of the Enqueues table. In the two flaiover docs it adds one sentence each on the save refusal. The trial conflict in mergequeue/eligible.go came from main's newer change, which S-0844's sync has now taken in. S-0844 expects to reach review within the hour.

### 2026-11-28T00:43:34Z orchestrator
Closed: S-0844 was accepted
