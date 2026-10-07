---
id: T-1198
type: task
nature: improvement
title: The workflow, the convention, the template, and the guides say conflicts and grown overlaps arrive as messages
status: backlog
parent: S-0332
owner: alex
created: 2026-10-07T20:15:25Z
updated: 2026-10-07T20:15:25Z
transitions: []
stream: S-0332
tags: [flai, template]
touches: [design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1195, T-1196, T-1197]
---
# T-1198 The workflow, the convention, the template, and the guides say conflicts and grown overlaps arrive as messages

## Work

Describe what T-1195 to T-1197 built. It waits for all three, so the words match the behaviour.

- `design/system/workflow.md` § Branches and collisions: the trial merge and a grown claim open a conversation, and escalation; link the ADR.
- `design/system/agent-narrative.md` § What an agent is told through MCP: the conversation beside the `overlapped` change.
- `work-management.md` in both copies: answer a conflict or overlap conversation, settle who changes what, and escalate only when the two do not agree. `template/CHANGELOG.md` records it.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai message escalate`, `message_escalate`, and what `flai stream sync` reports; regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- No document says a sync opens a conflict thread.
- The markdown lint passes on the changed files.
- `flai check --strict` passes.

## Notes
