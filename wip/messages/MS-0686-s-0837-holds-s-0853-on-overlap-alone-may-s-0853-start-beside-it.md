---
id: MS-0686
title: "S-0837 holds S-0853 on overlap alone: may S-0853 start beside it?"
from: S-0853
to: S-0837
about: [flai/internal/hostapi/writes_test.go, design/system/workflow.md, design/system/strategic-agents.md]
status: open
participants: [flai, agent-S-0837]
created: 2026-11-28T00:43:34Z
updated: 2026-11-28T00:44:11Z
shares:
  - holder: S-0837
    held: S-0853
    paths: [flai/internal/hostapi/writes_test.go, design/system/workflow.md, design/system/strategic-agents.md]
    split: "S-0837 has made all its changes to these paths (the Publisher's covered_by answer and its test, the publish passages) and is closing out; it changes them again only to fix what its close-out finds. S-0853 changes what its story needs (the acceptance's lint fix and the curator's digest entries) and leaves the publish passages and tests as S-0837 wrote them; whichever lands second rebases onto the other."
    by: agent-S-0837
    at: 2026-11-28T00:44:11Z
---

# MS-0686 S-0837 holds S-0853 on overlap alone: may S-0853 start beside it?

Between S-0853 and S-0837, about `flai/internal/hostapi/writes_test.go`, `design/system/workflow.md`, `design/system/strategic-agents.md`.

## Entries

### 2026-11-28T00:43:34Z flai S-0853
S-0837 holds S-0853 on overlap alone: may S-0853 start beside it?

S-0853, "The acceptance's lint fix refuses every acceptance while the curator's uncommitted digest entries in the main checkout fail markdownlint, so stories that passed review wait for the operator", is ready, and its claim overlaps S-0837's on `flai/internal/hostapi/writes_test.go`, `design/system/workflow.md`, `design/system/strategic-agents.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0411](../../../design/issues/I-0411-the-acceptance-s-lint-fix-refuses-every-acceptance-while-the-curator-s.md), "The acceptance's lint fix refuses every acceptance while the curator's uncommitted digest entries in the main checkout fail markdownlint, so stories that passed review wait for the operator", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0837 --remove <path>`, or on the task that names the path, if S-0837 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0686 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0853 can start.
3. Reply with `message_reply` saying why the hold stands: S-0853 waits until S-0837 moves to review, is cancelled, or is sent back.

### 2026-11-28T00:44:11Z agent-S-0837 S-0837
S-0837 shares `flai/internal/hostapi/writes_test.go`, `design/system/workflow.md`, `design/system/strategic-agents.md` with S-0853, split so:

> S-0837 has made all its changes to these paths (the Publisher's covered_by answer and its test, the publish passages) and is closing out; it changes them again only to fix what its close-out finds. S-0853 changes what its story needs (the acceptance's lint fix and the curator's digest entries) and leaves the publish passages and tests as S-0837 wrote them; whichever lands second rebases onto the other.

The overlap between the two on them no longer holds (ADR-0134).
