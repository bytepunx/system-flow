---
id: T-2111
type: task
nature: improvement
title: The design and the guides say that task done's check runs the flai the tree builds when the manifest names it
status: done
parent: S-0495
owner: alex
created: 2026-10-25T20:33:29Z
updated: 2026-10-25T20:58:01Z
transitions:
  - to: ready
    at: 2026-10-25T20:36:41Z
    by: agent-S-0495
  - to: in-progress
    at: 2026-10-25T20:56:05Z
    by: agent-S-0495
  - to: done
    at: 2026-10-25T20:58:01Z
    by: agent-S-0495
stream: S-0495
tags: [docs, cli, mcp]
touches: [design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md]
after: [T-2110]
usage:
  source: log
  seconds: 116
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 7338
      cache_read: 1070681
      cache_write: 57388
      cost: 0.7209
---
# T-2111 The design and the guides say that task done's check runs the flai the tree builds when the manifest names it

## Work

Describe what T-2110 built where task done and the MCP `task_done` are described, each with the ADR of T-2108 linked and `updated` bumped:

- `design/system/flai-cli.md` § Versions: the host's flai and the tree gains a bullet after "An outdated MCP server closes with the installed flai": the check step runs the flai the manifest's key names, in the worktree, with the command and arguments, `checked_by`, and the rule for a command that fails, naming I-0194 and the ADR. The `flai task done` and `task_done` rows of § Commands say the check may run the tree's flai.
- `design/system/agent-narrative.md` § What an agent is told through MCP: the `task_done` bullet (ADR-0160) says which flai checks when the manifest names one.
- `docs/users/flai.md` § Closing a task: the failed-check bullet and the MCP `task_done` paragraph say that in a project whose manifest names the flai its tree builds, the check runs that flai, the answer's `checked_by` names it, and a tree that does not build stops the close.

Waits for T-2110, whose behaviour it describes. Runs beside the issue close in its layer: they share no path.

## Done when

The three documents name the key, the command run, `checked_by`, and the failure rule, link the ADR, and `flai check --strict --story S-0495` and the markdown lint on the three report nothing.

## Notes

`design/conventions/work-management.md` is left alone: the close-out and the task cycle it describes do not change, only which flai runs the check.
