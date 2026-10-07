---
id: T-1201
type: task
nature: feature
title: The workflow, the convention, the template, and the guides say what an agent does with a task's change notice
status: done
parent: S-0333
owner: alex
created: 2026-10-07T20:15:43Z
updated: 2026-10-07T22:30:35Z
transitions:
  - to: ready
    at: 2026-10-07T22:27:28Z
    by: agent-S-0333
  - to: in-progress
    at: 2026-10-07T22:27:29Z
    by: agent-S-0333
  - to: done
    at: 2026-10-07T22:30:35Z
    by: agent-S-0333
stream: S-0333
tags: [flai, template]
touches: [design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/agent-coordination.md]
after: [T-1200]
usage:
  source: log
  seconds: 186
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 52
      output: 17530
      cache_read: 3167942
      cache_write: 84707
      cost: 1.5318
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
