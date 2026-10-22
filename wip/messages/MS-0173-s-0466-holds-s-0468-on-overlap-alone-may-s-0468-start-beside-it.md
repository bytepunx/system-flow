---
id: MS-0173
title: "S-0466 holds S-0468 on overlap alone: may S-0468 start beside it?"
from: S-0468
to: S-0466
about: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/system/workflow.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md]
status: open
participants: [flai, agent-S-0466]
created: 2026-10-22T21:01:42Z
updated: 2026-10-22T21:02:00Z
---

# MS-0173 S-0466 holds S-0468 on overlap alone: may S-0468 start beside it?

Between S-0468 and S-0466, about `flai/internal/harness/harness.go`, `flai/internal/harness/harness_test.go`, `design/system/workflow.md`, `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`.

## Entries

### 2026-10-22T21:01:42Z flai S-0468
S-0466 holds S-0468 on overlap alone: may S-0468 start beside it?

S-0468, "Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five", is ready, and its claim overlaps S-0466's on `flai/internal/harness/harness.go`, `flai/internal/harness/harness_test.go`, `design/system/workflow.md`, `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0152](../../../design/issues/I-0152-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md), "Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five", which the analysis in [2026-10-15-bottlenecks.md](../../../design/analysis/2026-10-15-bottlenecks.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0466 --remove <path>`, or on the task that names the path, if S-0466 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0173 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0468 can start.
3. Reply with `message_reply` saying why the hold stands: S-0468 waits until S-0466 moves to review, is cancelled, or is sent back.

### 2026-10-22T21:02:00Z agent-S-0466 S-0466 (settled)
The hold stands for now, and I expect to lift it soon. These five paths come from the planner's plan for S-0466: a detached `flai close-out` that changes the harness prompt and the close-out rules in work-management.md and workflow.md. I've asked the designer on TH-0594 whether to build that plan or a smaller fix: the close-out's flai reads a config without the merge-queue host action. The smaller fix changes none of these five paths. If the designer picks it, I'll remove all five from S-0466's touches and S-0468 can start. If they pick the planner's plan, S-0466 changes harness.go and the work-management.md close-out rules, and S-0468 waits until S-0466 moves to review.
