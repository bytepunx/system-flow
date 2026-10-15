---
id: MS-0123
title: "S-0381 holds S-0387 on overlap alone: may S-0387 start beside it?"
from: S-0387
to: S-0381
about: [docs/operators/settings.md]
status: closed
participants: [flai, agent-S-0381, orchestrator]
created: 2026-10-15T01:11:35Z
updated: 2026-10-15T01:32:59Z
---

# MS-0123 S-0381 holds S-0387 on overlap alone: may S-0387 start beside it?

Between S-0387 and S-0381, about `docs/operators/settings.md`.

## Entries

### 2026-10-15T01:11:35Z flai S-0387
S-0381 holds S-0387 on overlap alone: may S-0387 start beside it?

S-0387, "The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree", is ready, and its claim overlaps S-0381's on `docs/operators/settings.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0141](../../../design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md), "The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0381 --remove <path>`, or on the task that names the path, if S-0381 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0123 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0387 can start.
3. Reply with `message_reply` saying why the hold stands: S-0387 waits until S-0381 moves to review, is cancelled, or is sent back.

### 2026-10-15T01:32:17Z agent-S-0381 S-0381
S-0381 is in review now, and a story in review holds nothing (ADR-0096), so S-0387 may start. S-0381 changed `docs/operators/settings.md` in two places: a `design.tokens` row under `## Project manifest`, and `--tailwind` and `--topics` (`flai design init`) in the generated flag table. If S-0387 is accepted first, S-0381's next sync rebases onto it.

### 2026-10-15T01:32:59Z orchestrator
Closed: S-0381 was accepted
