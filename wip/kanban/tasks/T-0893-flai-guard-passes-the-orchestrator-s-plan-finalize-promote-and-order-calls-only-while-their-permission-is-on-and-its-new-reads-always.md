---
id: T-0893
type: task
nature: feature
title: flai guard passes the orchestrator's plan, finalize, promote, and order calls only while their permission is on, and its new reads always
status: backlog
parent: S-0219
owner: alex
created: 2026-10-05T04:46:49Z
updated: 2026-10-05T04:47:11Z
transitions: []
stream: S-0219
tags: [flai]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/orchestrate_test.go, flai/cmd/guard.go]
after: [T-0884, T-0888, T-0890]
---
# T-0893 flai guard passes the orchestrator's plan, finalize, promote, and order calls only while their permission is on, and its new reads always

## Work

S-0218 gives `flai guard` the orchestrator's role (`FLAI_ROLE=orchestrate`) and its refusal shape: a call outside the permissions is refused naming the permission that would allow it, and logged. This task maps the four permissions of this story onto calls, in `flai/internal/guard/guard.go`, extending what S-0218 left and keeping its shape:

| Permission | Passes while on | Refused while off |
|------------|-----------------|-------------------|
| `plan_backlog_epics` | the MCP tool `plan` and `flai plan <E-nnnn>`, for an epic only | both, and `plan` for a story at any time |
| `finalize_drafts` | `item_edit` with `draft: false`, `flai edit <S-nnnn> --no-draft`, with no other field changed | both |
| `promote_to_ready` | `item_move` and `flai move` of a story to `ready` | both |
| `order_ready` | `flai order --by <policy> --apply` | it, and `flai order` placing one story by hand at any time |

The reads `flai plan --candidates`, `flai promote --candidates`, `flai promote --drafts`, `flai order --by` without `--apply`, and `flai release --evaluate` pass always (`cliReads`). Whether the story moved is held, a draft, or over the ready limit is checked by `flai move` (T-0896), not here: the guard reads the call, not the board.

This task waits for T-0884, T-0888, and T-0890, for the commands and flags it names. It runs with T-0896, whose paths it does not share.

## Done when

- table tests in `orchestrate_test.go` cover each permission on and off for each call above, and each refusal names the permission
- the new reads pass with every permission off
- the planner's and the story agent's rules are unchanged: `guard_test.go` passes as it was
- `go test ./internal/guard/ ./cmd/` passes

## Notes
