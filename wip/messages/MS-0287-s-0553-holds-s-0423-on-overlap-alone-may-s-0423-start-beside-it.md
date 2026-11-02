---
id: MS-0287
title: "S-0553 holds S-0423 on overlap alone: may S-0423 start beside it?"
from: S-0423
to: S-0553
about: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/strategic-agents.md, docs/operators/index.md, flai/internal/check/check.go]
status: open
participants: [flai]
created: 2026-11-02T04:07:15Z
updated: 2026-11-02T04:07:15Z
---

# MS-0287 S-0553 holds S-0423 on overlap alone: may S-0423 start beside it?

Between S-0423 and S-0553, about `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`, `design/system/strategic-agents.md`, `docs/operators/index.md`, `flai/internal/check/check.go`.

## Entries

### 2026-11-02T04:07:15Z flai S-0423
S-0553 holds S-0423 on overlap alone: may S-0423 start beside it?

S-0423, "The template ships a baseline publish formula that fetches, evaluates orchestration.release, and publishes what is pending, the orchestrator starts it in place of calling release_publish, and it is run once here", is ready, and its claim overlaps S-0553's on `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`, `design/system/strategic-agents.md`, `docs/operators/index.md`, `flai/internal/check/check.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Publishing a release stops being prose. A baseline formula, `publish`, spells out what `work-management.md` and `strategic-agents.md` § Publishing a release say in words: fetch, evaluate the release policy, publish what is accepted and not yet released, and record the versions and tags. The orchestrator, under its `publish` permission, starts the formula when `release_evaluate` is met, in place of calling `release_publish` itself, so each publish is a story worked, measured, and archived. It is run once in this repository before the epic is accepted.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0553 --remove <path>`, or on the task that names the path, if S-0553 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0287 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0423 can start.
3. Reply with `message_reply` saying why the hold stands: S-0423 waits until S-0553 moves to review, is cancelled, or is sent back.
