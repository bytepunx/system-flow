---
id: T-1210
type: task
nature: improvement
title: The workflow, the narrative design, the convention, the template, and the guides say an agent ends on another agent's reply
status: backlog
parent: S-0335
owner: alex
created: 2026-10-07T20:16:45Z
updated: 2026-10-07T20:16:45Z
transitions: []
stream: S-0335
tags: [flai, template]
touches: [design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md]
after: [T-1207, T-1209]
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
