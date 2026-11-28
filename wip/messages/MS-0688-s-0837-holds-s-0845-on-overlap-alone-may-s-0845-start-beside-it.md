---
id: MS-0688
title: "S-0837 holds S-0845 on overlap alone: may S-0845 start beside it?"
from: S-0845
to: S-0837
about: [design/system/strategic-agents.md, design/system/workflow.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md]
status: open
participants: [flai, agent-S-0837]
created: 2026-11-28T00:43:34Z
updated: 2026-11-28T00:44:15Z
shares:
  - holder: S-0837
    held: S-0845
    paths: [design/system/strategic-agents.md, design/system/workflow.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md]
    split: "S-0837 has made all its changes to these paths (the orchestrator's publish with release_publish, and no formula) and is closing out. S-0845 changes how a review story whose close-out fails after main moved is sent back, and leaves the publish passages as S-0837 wrote them; whichever lands second rebases onto the other, keeping the role-file copies byte-equal."
    by: agent-S-0837
    at: 2026-11-28T00:44:15Z
---

# MS-0688 S-0837 holds S-0845 on overlap alone: may S-0845 start beside it?

Between S-0845 and S-0837, about `design/system/strategic-agents.md`, `design/system/workflow.md`, `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`.

## Entries

### 2026-11-28T00:43:34Z flai S-0845
S-0837 holds S-0845 on overlap alone: may S-0845 start beside it?

S-0845, "A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves", is ready, and its claim overlaps S-0837's on `design/system/strategic-agents.md`, `design/system/workflow.md`, `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0386](../../../design/issues/I-0386-a-story-in-review-whose-close-out-fails-after-main-moved-waits-for-the-operator.md), "A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves", which the analysis in [2026-11-16-all.md](../../../design/analysis/2026-11-16-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0837 --remove <path>`, or on the task that names the path, if S-0837 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0688 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0845 can start.
3. Reply with `message_reply` saying why the hold stands: S-0845 waits until S-0837 moves to review, is cancelled, or is sent back.

### 2026-11-28T00:44:15Z agent-S-0837 S-0837
S-0837 shares `design/system/strategic-agents.md`, `design/system/workflow.md`, `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md` with S-0845, split so:

> S-0837 has made all its changes to these paths (the orchestrator's publish with release_publish, and no formula) and is closing out. S-0845 changes how a review story whose close-out fails after main moved is sent back, and leaves the publish passages as S-0837 wrote them; whichever lands second rebases onto the other, keeping the role-file copies byte-equal.

The overlap between the two on them no longer holds (ADR-0134).
