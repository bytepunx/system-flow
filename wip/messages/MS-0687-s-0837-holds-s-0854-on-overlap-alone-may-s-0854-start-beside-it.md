---
id: MS-0687
title: "S-0837 holds S-0854 on overlap alone: may S-0854 start beside it?"
from: S-0854
to: S-0837
about: [design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, flai/internal/roles/roles_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/harness/harness_test.go]
status: open
participants: [flai]
created: 2026-11-28T00:43:34Z
updated: 2026-11-28T00:43:34Z
---

# MS-0687 S-0837 holds S-0854 on overlap alone: may S-0854 start beside it?

Between S-0854 and S-0837, about `design/system/strategic-agents.md`, `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`, `flai/internal/roles/roles_test.go`, `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, `flai/internal/roles/baseline/orchestrator.md`, `flai/internal/harness/harness_test.go`.

## Entries

### 2026-11-28T00:43:34Z flai S-0854
S-0837 holds S-0854 on overlap alone: may S-0854 start beside it?

S-0854, "A story in review whose close-out or resolve job passed is not a reason that starts the orchestrator, so an eligible story waits for the operator or for an unrelated event", is ready, and its claim overlaps S-0837's on `design/system/strategic-agents.md`, `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`, `flai/internal/roles/roles_test.go`, `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, `flai/internal/roles/baseline/orchestrator.md`, `flai/internal/harness/harness_test.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0413](../../../design/issues/I-0413-a-story-in-review-whose-close-out-or-resolve-job-passed-is-not-a-reason-that.md), "A story in review whose close-out or resolve job passed is not a reason that starts the orchestrator, so an eligible story waits for the operator or for an unrelated event", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0837 --remove <path>`, or on the task that names the path, if S-0837 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0687 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0854 can start.
3. Reply with `message_reply` saying why the hold stands: S-0854 waits until S-0837 moves to review, is cancelled, or is sent back.
