---
id: MS-0677
title: "S-0844 holds S-0822 on overlap alone: may S-0822 start beside it?"
from: S-0822
to: S-0844
about: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
status: open
participants: [flai, agent-S-0844]
created: 2026-11-28T00:21:40Z
updated: 2026-11-28T00:23:36Z
---

# MS-0677 S-0844 holds S-0822 on overlap alone: may S-0822 start beside it?

Between S-0822 and S-0844, about `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`.

## Entries

### 2026-11-28T00:21:40Z flai S-0822
S-0844 holds S-0822 on overlap alone: may S-0822 start beside it?

S-0822, "A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it", is ready, and its claim overlaps S-0844's on `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0396](../../../design/issues/I-0396-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md), "A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it". The issue recommends this solution:
>
> When a rebase and resolve job passes, `settle` in `flai/internal/serve/resolvejob.go` reads the story again, as `sendBack` does. For a story no longer in review it records an outcome of its own, not `stayed in review`, and a reason that names the state the story is in, who moved it there and when, and the command that moves it back to review, since its close-out passed at its head. The review page shows that outcome as it shows the others. Decide in the story, with the operator, whether a move out of review while the job runs should also stop the job; record the decision in an ADR refining ADR-0313 if it changes what ADR-0313 says. A test in the resolve job's tests moves the story out of review while the job runs and reads the record.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0844 --remove <path>`, or on the task that names the path, if S-0844 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0677 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0822 can start.
3. Reply with `message_reply` saying why the hold stands: S-0822 waits until S-0844 moves to review, is cancelled, or is sent back.

### 2026-11-28T00:23:36Z agent-S-0844 S-0844 (settled)
The hold stands for a short while. In design/system/flaiover-dashboard.md and docs/users/flaiover.md, S-0844 adds one sentence each where the document save is described: the editor saves design/workflows/board.md only on a story's branch. Whichever story lands second keeps both. The trial conflict in mergequeue/eligible.go came from main's newer change, which S-0844's sync has now taken in. S-0844 expects to reach review within the hour.
