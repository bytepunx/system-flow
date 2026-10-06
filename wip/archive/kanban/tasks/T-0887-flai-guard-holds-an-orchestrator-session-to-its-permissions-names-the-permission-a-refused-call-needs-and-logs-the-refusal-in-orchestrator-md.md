---
id: T-0887
type: task
nature: feature
title: flai guard holds an orchestrator session to its permissions, names the permission a refused call needs, and logs the refusal in orchestrator.md
status: done
parent: S-0218
owner: alex
created: 2026-10-05T04:46:28Z
updated: 2026-10-05T07:47:04Z
transitions:
  - to: ready
    at: 2026-10-05T07:12:47Z
    by: agent-S-0218
  - to: in-progress
    at: 2026-10-05T07:12:48Z
    by: agent-S-0218
  - to: review
    at: 2026-10-05T07:47:04Z
    by: agent-S-0218
  - to: done
    at: 2026-10-05T07:47:04Z
    by: agent-S-0218
stream: S-0218
tags: [flai]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/cmd/plan.go, flai/cmd/plan_test.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go]
after: [T-0881]
usage:
  source: log
  seconds: 2056
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 136
      output: 55910
      cache_read: 7969676
      cache_write: 202528
      cost: 3.8461
---
# T-0887 flai guard holds an orchestrator session to its permissions, names the permission a refused call needs, and logs the refusal in orchestrator.md

## Work

With `FLAI_ROLE=orchestrate`, `flai guard` checks the session's own calls, as it does the planner's (`Guard.plan` in `flai/internal/guard/guard.go`), against rules of the orchestrator's that depend on `orchestration.permissions`. `flai/cmd/guard.go` reads the permissions from the manifest and hands them to `Guard`.

| Always passes | Passes only with |
|---------------|------------------|
| the reads a sub-agent has; `inbox`, `board`, `activity_log`, `wait_for_events`, `thread_open`; S-0217's reads: `flai order --by` without `--apply`, `flai promote --candidates`, `flai release --evaluate`, and their MCP tools | `plan_backlog_epics`: the MCP tool `plan` and `flai plan`, on an epic in the backlog |
| | `finalize_drafts`: `flai edit --no-draft` |
| | `promote_to_ready`: `item_move` and `flai move` to `ready` |
| | `order_ready`: `flai order` that writes, `--apply` among them |
| | `answer_threads`, `recommend` or `autonomous`: `thread_reply` and `flai reply` |
| | `accept_reviews`: `flai accept` |
| | `publish`: `flai release --pending` and `flai push` |

Everything else that writes is refused: other flai tools and commands, git's writes, and `Edit`, `Write`, and `NotebookEdit`. A refusal of a call a permission would allow names it, such as `needs orchestration.permissions.promote_to_ready`, and says to ask the operator on the item; other refusals say the orchestrator never does it. The orchestrator's sub-agents are held as every sub-agent is.

Each refusal is logged: `flai/internal/workitem/activity.go` gains an entry for a refusal in `wip/agents/orchestrator.md`, with the time, the call, and the permission it needs, and no seconds or cost, written under the same lock as other entries. `flai/cmd/plan.go` lets an `orchestrate` session started by `flai serve` ask for the planner, on a backlog epic only, where it refuses every other agent `flai serve` started.

This task waits for T-0881, for the permissions it reads. It runs with T-0883, whose paths it does not share.

## Done when

- a permitted call and a refused call pass and fail for each permission, on and off, in table tests in `guard_test.go`
- a refusal names the permission that would allow it, and a test reads its entry in a fixture `orchestrator.md`
- `flai guard` run as a hook with `FLAI_ROLE=orchestrate` reads the fixture manifest's permissions, in `cmd/guard_test.go`
- the MCP tool `plan` starts the planner for a backlog epic for the orchestrator with the permission and refuses it otherwise, in `cmd/plan_test.go`
- `go test ./internal/guard/ ./internal/workitem/ ./cmd/` passes

## Notes
