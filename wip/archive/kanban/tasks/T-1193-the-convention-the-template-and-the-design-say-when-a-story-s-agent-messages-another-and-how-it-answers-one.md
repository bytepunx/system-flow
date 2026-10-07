---
id: T-1193
type: task
nature: feature
title: The convention, the template, and the design say when a story's agent messages another and how it answers one
status: done
parent: S-0331
owner: alex
created: 2026-10-07T20:14:49Z
updated: 2026-10-07T21:29:02Z
transitions:
  - to: ready
    at: 2026-10-07T21:24:49Z
    by: agent-S-0331
  - to: in-progress
    at: 2026-10-07T21:24:49Z
    by: agent-S-0331
  - to: done
    at: 2026-10-07T21:29:02Z
    by: agent-S-0331
stream: S-0331
tags: [flai, template]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, design/conventions/delegation.md, design/system/workflow.md, template/root/design/conventions/delegation.md]
after: [T-1191, T-1192]
usage:
  source: log
  seconds: 253
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 64
      output: 22923
      cache_read: 3907952
      cache_write: 139673
      cost: 2.085
---
# T-1193 The convention, the template, and the design say when a story's agent messages another and how it answers one

## Work

Tell agents and readers what T-1190 to T-1192 built. It waits for T-1191 and T-1192, so the words match the tools and the guard.

- `work-management.md`, in `design/conventions` and `template/root/design/conventions`: message the other story's agent before changing a path both claim, answer a message to your story before going on, and ask the operator on a thread only when the two do not agree. Replace the rules that say to coordinate on a thread.
- `template/CHANGELOG.md`: the convention change.
- `design/system/agent-narrative.md` § What an agent is told through MCP: the tools, the `messages` field, and the `message` event.
- `design/system/flai-cli.md` and `docs/users/flai.md`: the MCP tools.

## Done when

- Both copies of the convention say the same, and no rule left tells an agent to coordinate with another on a thread.
- The markdown lint passes on the changed files.
- `flai check --strict` passes.

## Notes
