---
id: T-0785
type: task
nature: feature
title: "flai guard holds a planner session to planning: items and threads through flai, no move past backlog, no accept, publish, or edit"
status: done
parent: S-0208
owner: alex
created: 2026-10-04T00:43:46Z
updated: 2026-10-04T00:52:31Z
transitions:
  - to: ready
    at: 2026-10-04T00:44:31Z
    by: agent-S-0208
  - to: in-progress
    at: 2026-10-04T00:44:33Z
    by: agent-S-0208
  - to: done
    at: 2026-10-04T00:52:31Z
    by: agent-S-0208
stream: S-0208
tags: []
touches: [flai/internal/guard, flai/cmd/guard.go, flai/cmd/guard_test.go, template/root/.claude/settings.json, ".claude/settings.json"]
usage:
  source: log
  seconds: 478
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 67
      output: 28402
      cache_read: 4101140
      cache_write: 109292
      cost: 2.0235
---
# T-0785 flai guard holds a planner session to planning: items and threads through flai, no move past backlog, no accept, publish, or edit

## Work

- A session with `FLAI_ROLE=plan` in its environment is the planner (the hook runs in the session's environment). `flai guard` reads it and holds the planner's own calls (no agent ID) to planning; its sub-agents' calls (an agent ID) stay held as every sub-agent's are.
- The planner may: every read a sub-agent may; flai's MCP tools `inbox`, `item_new`, `item_edit`, `thread_open`, `thread_reply`, `activity_log`, `wait_for_events`, and `item_move` only to `backlog`; the CLI's `story new`, `epic new`, `edit`, `touches`, `thread new` and `reply`, `move <id> backlog`, `issue new` and `bump`, and the reads.
- The planner may not: move an item anywhere but `backlog`, accept, release, push, publish, archive, open or sync a stream, run a git command that writes, or call `Edit`, `Write`, or `NotebookEdit`.
- The hook's matcher in the template's and the repository's `.claude/settings.json` adds `Edit|Write|NotebookEdit` so the guard sees those calls; for a story's agent and its sub-agents nothing changes.
- The refusal says the planner cannot do it, why (`strategic-agents.md`), and what to do instead (a thread to the operator, or end with it in the summary).

Waits for nothing: the first layer, beside T-0784, with no path in common; the environment name `FLAI_ROLE=plan` is settled in the narrative's Decisions.

## Done when

- [ ] Tests show the planner allowed each write above and refused each forbidden call, `item_move` to `backlog` allowed and to `ready` refused, and a story's agent and a sub-agent decided as before
- [ ] Both `settings.json` files carry the same matcher

## Notes
