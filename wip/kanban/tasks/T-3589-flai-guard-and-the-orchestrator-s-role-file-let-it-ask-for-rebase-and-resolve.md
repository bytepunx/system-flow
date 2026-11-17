---
id: T-3589
type: task
nature: improvement
title: flai guard and the orchestrator's role file let it ask for rebase and resolve under accept_reviews when a close-out did not pass after main moved
status: backlog
parent: S-0845
owner: alex
created: 2026-11-17T12:34:14Z
updated: 2026-11-17T12:34:14Z
transitions: []
stream: S-0845
tags: [flai, orchestrator]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_rules.go, flai/internal/guard/orchestrate_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md]
after: [T-3588]
---
# T-3589 flai guard and the orchestrator's role file let it ask for rebase and resolve under accept_reviews when a close-out did not pass after main moved

## Work

Waits for T-3588, whose MCP tool `resolve` this task names in the role file and the guard. In `design/workflows/roles/orchestrator.md` and its copy under `template/root/`, add `mcp__flai__resolve` to the definition's tools and `resolve` to the `orchestration.permissions.accept_reviews` permission's tools beside `queue_add` and `accept`. Rewrite the paragraph that says "A run that did not pass is a blocker" and the matching step 12 so that, when the story's last close-out did not pass and main has moved past its base, the orchestrator asks for rebase and resolve with `resolve`, logs it with `activity_log`, opens no thread, and leaves the story for flai serve's job to keep in review or send back, as T-3587's ADR says. In `flai/internal/guard/guard.go` and `guard_rules.go`, let the orchestrator's `flai accept --resolve --by orchestrator` and the `resolve` tool pass only while `accept_reviews` is on, and test both ways in `orchestrate_test.go`.

## Done when

- The role file and its template copy name the tool and say when to use it, and `flai check --strict` passes on the role files.
- The guard passes the orchestrator's resolve request with `accept_reviews` on and refuses it with it off, each with a test.
- `flai test` passes on the paths changed.

## Notes
