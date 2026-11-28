---
id: MS-0675
title: "S-0844 holds S-0852 on overlap alone: may S-0852 start beside it?"
from: S-0852
to: S-0844
about: [flai/cmd/accept.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, design/system/workflow.md]
status: open
participants: [flai, agent-S-0844]
created: 2026-11-28T00:21:40Z
updated: 2026-11-28T00:23:32Z
---

# MS-0675 S-0844 holds S-0852 on overlap alone: may S-0852 start beside it?

Between S-0852 and S-0844, about `flai/cmd/accept.go`, `flai/internal/preview/accept.go`, `flai/internal/preview/accept_test.go`, `flai/internal/mergequeue/eligible.go`, `flai/internal/mergequeue/eligible_test.go`, `design/system/workflow.md`.

## Entries

### 2026-11-28T00:21:40Z flai S-0852
S-0844 holds S-0852 on overlap alone: may S-0852 start beside it?

S-0852, "The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right", is ready, and its claim overlaps S-0844's on `flai/cmd/accept.go`, `flai/internal/preview/accept.go`, `flai/internal/preview/accept_test.go`, `flai/internal/mergequeue/eligible.go`, `flai/internal/mergequeue/eligible_test.go`, `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0415](../../../design/issues/I-0415-the-orchestrator-s-accept-takes-no-keep-removal-so-a-story-that-rewrites-a-line.md), "The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0844 --remove <path>`, or on the task that names the path, if S-0844 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0675 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0852 can start.
3. Reply with `message_reply` saying why the hold stands: S-0852 waits until S-0844 moves to review, is cancelled, or is sent back.

### 2026-11-28T00:23:32Z agent-S-0844 S-0844 (settled)
The hold stands for a short while. S-0844 changes accept.go, preview/accept.go, mergequeue/eligible.go, and their tests. It adds a workflow_file blocker, the queue kind workflow_file, and AcceptOptions.LandsByHand, which keeps a story that changes the workflow file out of the queue. S-0852 will build on these. The trial conflict in eligible.go came from main's newer change there, which S-0844's sync has now taken in. S-0844 expects to reach review within the hour. The in-progress limit is full meanwhile.
