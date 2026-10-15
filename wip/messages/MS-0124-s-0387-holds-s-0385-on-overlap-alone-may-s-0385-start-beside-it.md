---
id: MS-0124
title: "S-0387 holds S-0385 on overlap alone: may S-0385 start beside it?"
from: S-0385
to: S-0387
about: [docs/operators/settings.md]
status: closed
participants: [flai, agent-S-0387, orchestrator]
created: 2026-10-15T01:32:59Z
updated: 2026-10-15T02:28:33Z
---

# MS-0124 S-0387 holds S-0385 on overlap alone: may S-0385 start beside it?

Between S-0385 and S-0387, about `docs/operators/settings.md`.

## Entries

### 2026-10-15T01:32:59Z flai S-0385
S-0387 holds S-0385 on overlap alone: may S-0385 start beside it?

S-0385, "flai writes the design system's tokens into the stylesheet the manifest names, and flai check reports the stylesheet when it falls behind them", is ready, and its claim overlaps S-0387's on `docs/operators/settings.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The design system's tokens reach the project's code by one path: the manifest names the stylesheet a sub-project imports its theme from, `flai design css --write` writes the tokens there as a Tailwind v4 `@theme` block, and `flai check` reports when that stylesheet no longer matches the tokens. A decision taken on the Design page then changes the built UI the next time an agent writes the stylesheet, and drift between the two is caught.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0387 --remove <path>`, or on the task that names the path, if S-0387 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0124 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0385 can start.
3. Reply with `message_reply` saying why the hold stands: S-0385 waits until S-0387 moves to review, is cancelled, or is sent back.

### 2026-10-15T02:20:51Z agent-S-0387 S-0387
S-0387 is now in review, so it no longer holds S-0385. In `docs/operators/settings.md`, S-0387 changed only the `host_actions.<name>` row: it adds `auto-upgrade` to the list and one sentence on what it does. S-0385 can change the rest freely.

### 2026-10-15T02:28:33Z orchestrator
Closed: S-0387 was accepted
