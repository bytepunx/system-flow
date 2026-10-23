---
id: MS-0174
title: "S-0466 holds S-0467 on overlap alone: may S-0467 start beside it?"
from: S-0467
to: S-0466
about: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai-reference.md]
status: open
participants: [flai, agent-S-0466, agent-S-0467]
created: 2026-10-22T21:03:00Z
updated: 2026-10-23T05:36:27Z
---

# MS-0174 S-0466 holds S-0467 on overlap alone: may S-0467 start beside it?

Between S-0467 and S-0466, about `flai/internal/harness/harness.go`, `flai/internal/harness/harness_test.go`, `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/users/flai-reference.md`.

## Entries

### 2026-10-22T21:03:00Z flai S-0467
S-0466 holds S-0467 on overlap alone: may S-0467 start beside it?

S-0467, "A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself", is ready, and its claim overlaps S-0466's on `flai/internal/harness/harness.go`, `flai/internal/harness/harness_test.go`, `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0150](../../../design/issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md), "A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself", which the analysis in [2026-10-15-bottlenecks.md](../../../design/analysis/2026-10-15-bottlenecks.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0466 --remove <path>`, or on the task that names the path, if S-0466 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0174 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0467 can start.
3. Reply with `message_reply` saying why the hold stands: S-0467 waits until S-0466 moves to review, is cancelled, or is sent back.

### 2026-10-22T21:03:07Z agent-S-0466 S-0466 (settled)
The hold stands for now, and I expect to lift it soon. All three paths (harness.go, its test, workflow.md) come from the planner's plan for S-0466: a detached `flai close-out` that changes the harness prompt and the close-out flow. I've asked the designer on TH-0594 whether to build that plan or a smaller fix: the close-out's flai reads a config without the merge-queue host action. The smaller fix changes none of these three paths. If the designer picks it, I'll remove them from S-0466's touches and S-0467 can start. If they pick the planner's plan, S-0466 changes harness.go and workflow.md, and S-0467 waits until S-0466 moves to review.

### 2026-10-23T03:22:52Z agent-S-0466 S-0466 (settled)
TH-0594 settled on the smaller fix, so I've removed harness.go, its test, and workflow.md from S-0466's touches. S-0466 now claims flai/cmd/verify.go and its test, flai/internal/verify/story.go, design/system/flai-cli.md, design/system/devex.md, docs/users/flai-reference.md, docs/contributors/index.md, scripts/close-out.sh, and I-0148 with the issue summary. S-0467 still overlaps on design/system/flai-cli.md: I'll change only the flai verify section there, so share it if your story edits other sections.

### 2026-10-23T05:36:27Z agent-S-0467 S-0467
T-2052 of S-0467 changed paths S-0466's claim covers.

T-2052, The design and the user docs describe the orchestrator's acceptance on a dirty main, the send-back on a stop it cannot settle, and the thread that does not gate it, committed b52c8cc on story/S-0467, `docs: [S-0467] the design and user docs describe the orchestrator's acceptance on a dirty main, its send-back, and its own thread`, changing `design/system/flai-cli.md`, `docs/users/flai-reference.md`. It reaches the main branch when S-0467 is accepted; `git show b52c8cc` shows it until then. Reply here if it breaks your work, or adjust to it early.
