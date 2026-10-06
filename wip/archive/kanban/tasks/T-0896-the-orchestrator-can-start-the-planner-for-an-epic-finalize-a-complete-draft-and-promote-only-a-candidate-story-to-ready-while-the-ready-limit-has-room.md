---
id: T-0896
type: task
nature: feature
title: The orchestrator can start the planner for an epic, finalize a complete draft, and promote only a candidate story to ready while the ready limit has room
status: done
parent: S-0219
owner: alex
created: 2026-10-05T04:47:01Z
updated: 2026-10-06T03:14:52Z
transitions:
  - to: ready
    at: 2026-10-06T03:03:31Z
    by: agent-S-0219
  - to: in-progress
    at: 2026-10-06T03:03:31Z
    by: agent-S-0219
  - to: done
    at: 2026-10-06T03:14:52Z
    by: agent-S-0219
stream: S-0219
tags: [flai]
touches: [flai/internal/mcpserver/plan.go, flai/internal/mcpserver/plan_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/agents.go, flai/cmd/move.go, flai/cmd/move_orchestrate_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/cmd/plan.go, flai/cmd/plan_test.go, flai/cmd/plan_orchestrate_test.go, flai/internal/workitem/promotable.go]
after: [T-0890]
usage:
  source: log
  seconds: 681
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 191
      output: 71034
      cache_read: 10776529
      cache_write: 248533
      cost: 4.987
---
# T-0896 The orchestrator can start the planner for an epic, finalize a complete draft, and promote only a candidate story to ready while the ready limit has room

## Work

The guard says whether the orchestrator may make a call; this task makes the calls work for it and holds them to the criteria, so that they hold whoever runs them as the orchestrator. Each applies only to `FLAI_ROLE=orchestrate`; every other caller is unchanged.

- **Plan.** The MCP tool `plan` (`mcpserver/plan.go`) is refused to an agent `flai serve` started. Let the orchestrator through for an epic, and record its runs under `plans` with the trigger `orchestrator` in place of `asked` (`serve/plan.go`), so its activity entry in `planner.md` says who asked.
- **Finalize.** `item_edit` refuses `draft: false` (`mcpserver/items_write.go`), and `flai edit --no-draft` stamps `finalized` (`cmd/edit.go`). Let the orchestrator finalize, stamped `finalized.by` with its agent name, only a draft that the first layer's draft check (`flai promote --drafts`) finds complete; refuse an incomplete one naming what it lacks.
- **Promote.** `flai move` and `item_move` to `ready` move a held story with a warning (`cmd/move.go`). For the orchestrator, refuse a story that is not among `flai promote --candidates` (a draft or a held story among them) with the candidates' reason, and refuse any move while the ready column is at its WIP limit.

This task waits for T-0890, for the draft check. It runs with the guard task, whose paths it does not share.

## Done when

- with `FLAI_ROLE=orchestrate`: `plan` starts the planner for an epic and records the trigger `orchestrator`, and is refused for a story
- with it: finalizing a complete draft stamps the orchestrator as who finalized it, and an incomplete one is refused naming what it lacks
- with it: moving a candidate to ready succeeds; a held story, a draft, and any story with ready at its limit are refused with the reason
- without it, each call behaves as before: the existing tests pass unchanged
- `go test ./internal/mcpserver/ ./internal/serve/ ./cmd/` passes

## Notes
