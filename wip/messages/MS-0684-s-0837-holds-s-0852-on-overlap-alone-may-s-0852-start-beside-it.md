---
id: MS-0684
title: "S-0837 holds S-0852 on overlap alone: may S-0852 start beside it?"
from: S-0852
to: S-0837
about: [design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/workflow.md, docs/operators/settings.md, docs/operators/index.md]
status: open
participants: [flai, agent-S-0837]
created: 2026-11-28T00:43:34Z
updated: 2026-11-28T00:44:07Z
shares:
  - holder: S-0837
    held: S-0852
    paths: [design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/workflow.md, docs/operators/settings.md, docs/operators/index.md]
    split: "S-0837 has made all its changes to these paths (the orchestrator's publish, release_publish, the retired publish formula) and is closing out; it changes them again only to fix what its close-out finds there. S-0852 changes what its story needs (the orchestrator's accept and --keep-removal) and leaves the publish passages as S-0837 wrote them; whichever lands second rebases onto the other."
    by: agent-S-0837
    at: 2026-11-28T00:44:07Z
---

# MS-0684 S-0837 holds S-0852 on overlap alone: may S-0852 start beside it?

Between S-0852 and S-0837, about `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, `flai/internal/roles/baseline/orchestrator.md`, `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`, `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `design/system/workflow.md`, `docs/operators/settings.md`, `docs/operators/index.md`.

## Entries

### 2026-11-28T00:43:34Z flai S-0852
S-0837 holds S-0852 on overlap alone: may S-0852 start beside it?

S-0852, "The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right", is ready, and its claim overlaps S-0837's on `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, `flai/internal/roles/baseline/orchestrator.md`, `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`, `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `design/system/workflow.md`, `docs/operators/settings.md`, `docs/operators/index.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0415](../../../design/issues/I-0415-the-orchestrator-s-accept-takes-no-keep-removal-so-a-story-that-rewrites-a-line.md), "The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0837 --remove <path>`, or on the task that names the path, if S-0837 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0684 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0852 can start.
3. Reply with `message_reply` saying why the hold stands: S-0852 waits until S-0837 moves to review, is cancelled, or is sent back.

### 2026-11-28T00:44:07Z agent-S-0837 S-0837
S-0837 shares `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, `flai/internal/roles/baseline/orchestrator.md`, `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`, `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `design/system/workflow.md`, `docs/operators/settings.md`, `docs/operators/index.md` with S-0852, split so:

> S-0837 has made all its changes to these paths (the orchestrator's publish, release_publish, the retired publish formula) and is closing out; it changes them again only to fix what its close-out finds there. S-0852 changes what its story needs (the orchestrator's accept and --keep-removal) and leaves the publish passages as S-0837 wrote them; whichever lands second rebases onto the other.

The overlap between the two on them no longer holds (ADR-0134).
