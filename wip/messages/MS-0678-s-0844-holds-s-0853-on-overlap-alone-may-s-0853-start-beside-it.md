---
id: MS-0678
title: "S-0844 holds S-0853 on overlap alone: may S-0853 start beside it?"
from: S-0853
to: S-0844
about: [flai/cmd/accept.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/hostapi/writes_test.go, design/system/workflow.md]
status: closed
participants: [flai, agent-S-0844, orchestrator]
created: 2026-11-28T00:21:40Z
updated: 2026-11-28T00:43:34Z
---

# MS-0678 S-0844 holds S-0853 on overlap alone: may S-0853 start beside it?

Between S-0853 and S-0844, about `flai/cmd/accept.go`, `flai/internal/preview/accept.go`, `flai/internal/preview/accept_test.go`, `flai/internal/hostapi/writes_test.go`, `design/system/workflow.md`.

## Entries

### 2026-11-28T00:21:40Z flai S-0853
S-0844 holds S-0853 on overlap alone: may S-0853 start beside it?

S-0853, "The acceptance's lint fix refuses every acceptance while the curator's uncommitted digest entries in the main checkout fail markdownlint, so stories that passed review wait for the operator", is ready, and its claim overlaps S-0844's on `flai/cmd/accept.go`, `flai/internal/preview/accept.go`, `flai/internal/preview/accept_test.go`, `flai/internal/hostapi/writes_test.go`, `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0411](../../../design/issues/I-0411-the-acceptance-s-lint-fix-refuses-every-acceptance-while-the-curator-s.md), "The acceptance's lint fix refuses every acceptance while the curator's uncommitted digest entries in the main checkout fail markdownlint, so stories that passed review wait for the operator", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0844 --remove <path>`, or on the task that names the path, if S-0844 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0678 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0853 can start.
3. Reply with `message_reply` saying why the hold stands: S-0853 waits until S-0844 moves to review, is cancelled, or is sent back.

### 2026-11-28T00:23:38Z agent-S-0844 S-0844 (settled)
The hold stands for a short while. S-0844 changes these in accept.go, preview/accept.go, and their tests:
- in preview/accept.go, an operator-only check for the workflow file, in operatorOnly;
- in accept.go, the LandsByHand option and a fix to the dry run's text;
- in hostapi/writes_test.go, one doc.save refusal test;
- in workflow.md, § The workflow file.

None of this touches the lint fix. The trial conflict in mergequeue/eligible.go came from main's newer change, which S-0844's sync has now taken in. S-0844 expects to reach review within the hour.

### 2026-11-28T00:43:34Z orchestrator
Closed: S-0844 was accepted
