---
id: T-1595
type: task
nature: improvement
title: Record in a new ADR that an outdated MCP server's task_done closes the task with the installed flai once it has caught up
status: done
parent: S-0390
owner: alex
created: 2026-10-15T02:59:22Z
updated: 2026-10-15T03:03:13Z
transitions:
  - to: ready
    at: 2026-10-15T03:02:04Z
    by: agent-S-0390
  - to: in-progress
    at: 2026-10-15T03:02:05Z
    by: agent-S-0390
  - to: done
    at: 2026-10-15T03:03:00Z
    by: agent-S-0390
stream: S-0390
tags: [adr, design]
touches: [design/adrs/README.md, design/system/flai-cli.md]
usage:
  source: log
  seconds: 55
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 12
      output: 127
      cache_read: 863926
      cache_write: 19934
      cost: 0.4097
---
# T-1595 Record in a new ADR that an outdated MCP server's task_done closes the task with the installed flai once it has caught up

## Work

Write the remedy for I-0147 as an ADR with `flai adr new` that refines ADR-0142 and ADR-0158, and update `design/system/flai-cli.md` § Versions: the host's flai and the tree. This task waits for nothing. It comes first because the story's goal asks for a proposal before the build.

The cause, from the four instances: a story's agent's MCP server is a stdio child of its claude session. `flai serve` starts it with its own binary's path (`flai/internal/harness/adapters.go`, `mcpConfig`). The child keeps the binary it started with for the whole session. Even after ADR-0158 upgrades the host, the session's `task_done` still runs the old flai. ADR-0142 refuses it a minor version behind, so the agent falls back to the tree's `scripts/flai.sh`, which a project made from the template does not have.

The proposed remedy, to confirm or change:

- When ADR-0142's check would refuse, the MCP tool `task_done` reads the version of the binary now at its own path (`os.Executable`).
- When that version is at least the newest release in the project's history, `task_done` closes the task by running that binary's `task done` in the story's worktree, with the same arguments, agent, and session. It returns the binary's result and says which flai closed the task.
- Otherwise it refuses as today. The refusal says the installed flai is behind too, and that `flai serve` upgrades the host or asks the operator (ADR-0158). It tells the agent to call again once the host is upgraded, not to fall back to another flai.
- Other write tools keep ADR-0117's warn-only rule.

## Done when

- The ADR is accepted, or proposed, with its context, decision, consequences, and alternatives (re-exec the stdio server, restart the session, refuse only), and `design/adrs/README.md` lists it.
- `design/system/flai-cli.md` § Versions says what `task_done` does when the session's flai is behind and the installed flai is not.

## Notes

The ADR's number and slug are predicted; `flai task done` widens the touches with the real path.
