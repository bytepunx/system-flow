---
id: T-1201
type: task
nature: feature
title: The workflow, the convention, the template, and the guides say what an agent does with a task's change notice
status: backlog
parent: S-0333
owner: alex
created: 2026-10-07T20:15:43Z
updated: 2026-10-07T20:15:43Z
transitions: []
stream: S-0333
tags: [flai, template]
touches: [design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1200]
---
# T-1201 The workflow, the convention, the template, and the guides say what an agent does with a task's change notice

## Work

Describe the notice T-1200 sends. It waits for T-1200, so the words match what is sent.

- `design/system/workflow.md` § Branches and collisions: the notice at a task's close, beside the notice at acceptance (S-0132).
- `design/system/agent-narrative.md` § What an agent is told through MCP: the message and what it carries.
- `work-management.md` in both copies: on such a message, check whether the change breaks the story's work, answer when it does, and sync before relying on it. `template/CHANGELOG.md` records it.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `told` in `flai task done`; regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- Each document names the notice and what the agent does with it.
- The markdown lint passes on the changed files.
- `flai check --strict` passes.

## Notes
