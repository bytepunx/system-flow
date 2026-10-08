---
id: T-1394
type: task
nature: improvement
title: The hold-and-ask moves out of permission_prompt into flai/internal/ask, with the bound passed by the caller
status: backlog
parent: S-0354
owner: alex
created: 2026-10-08T08:48:02Z
updated: 2026-10-08T08:48:02Z
transitions: []
stream: S-0354
tags: [cli]
touches: [flai/internal/ask/ask.go, flai/internal/ask/ask_test.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go]
---
# T-1394 The hold-and-ask moves out of permission_prompt into flai/internal/ask, with the bound passed by the caller

## Work

- `flai/internal/ask`: `Request` (the story, the agent, the tool, the path, the change text) and `Hold(ctx, repo, req, bound) (allowed bool, reason string)`, holding what `permission.go`'s `askOperator`, `askedBefore`, `awaitAnswer`, `answerOn`, `settle`, `allowed`, `answerers`, and `permissionRequest` do today: the thread `Allow <tool> <path>?`, the answerers of ADR-0097, the refusal that leaves the thread open after the bound (ADR-0124), and the answer taken on the same request again.
- `permission.go` keeps Claude Code's contract (`tool_name`, `input`, `tool_use_id`, `behavior`, `updatedInput`, `message`), the worktree and `.git` checks, and auto-approve, and calls `ask.Hold` with `permissionWait`.
- `ask_test.go` covers an allow, a deny, a timeout, an answer by the project's owner, and an answer taken on a second request; `permission_test.go` keeps its expectations.

First layer: it runs together with T-1395, which touches other files.

## Done when

- `flai test flai/internal/ask/ flai/internal/mcpserver/` passes with no expectation in `permission_test.go` changed.

## Notes

Layer 1 of S-0354.
