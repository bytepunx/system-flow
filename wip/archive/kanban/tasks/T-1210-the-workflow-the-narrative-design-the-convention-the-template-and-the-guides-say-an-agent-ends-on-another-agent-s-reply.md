---
id: T-1210
type: task
nature: improvement
title: The workflow, the narrative design, the convention, the template, and the guides say an agent ends on another agent's reply
status: done
parent: S-0335
owner: alex
created: 2026-10-07T20:16:45Z
updated: 2026-10-07T22:02:09Z
transitions:
  - to: ready
    at: 2026-10-07T21:59:20Z
    by: agent-S-0335
  - to: in-progress
    at: 2026-10-07T21:59:21Z
    by: agent-S-0335
  - to: done
    at: 2026-10-07T22:02:09Z
    by: agent-S-0335
stream: S-0335
tags: [flai, template]
touches: [design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md]
after: [T-1207, T-1209]
usage:
  source: log
  seconds: 168
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 52
      output: 21207
      cache_read: 2906780
      cache_write: 93877
      cost: 1.5832
---
# T-1210 The workflow, the narrative design, the convention, the template, and the guides say an agent ends on another agent's reply

## Work

Describe what T-1207 to T-1209 built. It waits for T-1207 and T-1209, and T-1209 for T-1208, so the words match the behaviour.

- `design/system/workflow.md` § Branches and collisions, **An answered agent**: a reply or a new message starts the agent again.
- `design/system/agent-narrative.md` § What an agent is told through MCP: `end: true` on a conversation.
- `work-management.md` in both copies: when nothing is left but another agent's reply, write Current state and Next steps and end, as for the operator's answer. `template/CHANGELOG.md` records it.
- `design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/users/flaiover.md`: the waiting reason.

## Done when

- Each document says an agent ends on another agent's reply and is started again.
- The markdown lint passes on the changed files.
- `flai check --strict` passes.

## Notes
