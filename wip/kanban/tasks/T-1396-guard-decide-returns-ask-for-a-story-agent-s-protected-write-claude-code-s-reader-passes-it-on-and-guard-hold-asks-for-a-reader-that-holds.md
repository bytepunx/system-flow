---
id: T-1396
type: task
nature: improvement
title: guard.Decide returns ask for a story agent's protected write, Claude Code's reader passes it on, and guard.Hold asks for a reader that holds
status: backlog
parent: S-0354
owner: alex
created: 2026-10-08T08:48:15Z
updated: 2026-10-08T08:48:15Z
transitions: []
stream: S-0354
tags: [cli]
touches: [flai/internal/guard/call.go, flai/internal/guard/call_test.go, flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go]
after: [T-1394, T-1395]
---
# T-1396 guard.Decide returns ask for a story agent's protected write, Claude Code's reader passes it on, and guard.Hold asks for a reader that holds

## Work

- In `Decide`, a `FileWrite` by a story's own agent to a path `protected.Path` names, inside its worktree and not under `.git`, while auto-approve is off, returns `Ask` with the reason `permission_prompt` gives today. A sub-agent's such write keeps ADR-0102's deny.
- Claude Code's reader maps `Ask` to exit 0, so that Claude Code's own permission handling calls `permission_prompt` as today.
- `guard.Hold(ctx, call, bound)` asks through `ask.Hold` and returns allow or deny, for a reader whose harness cannot call a permission tool; a test drives it with a fake reader and a fake thread answer.

It waits for T-1394, whose `ask.Hold` it calls, and T-1395, whose list it reads.

## Done when

- `call_test.go` covers the ask verdict, the sub-agent's deny, and auto-approve on; `claudecode_test.go` covers the pass-through.
- `flai test flai/internal/guard/` passes.

## Notes

Layer 2 of S-0354.
