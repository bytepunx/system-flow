---
id: T-1198
type: task
nature: improvement
title: The workflow, the convention, the template, and the guides say conflicts and grown overlaps arrive as messages
status: done
parent: S-0332
owner: alex
created: 2026-10-07T20:15:25Z
updated: 2026-10-07T23:12:36Z
transitions:
  - to: ready
    at: 2026-10-07T23:08:58Z
    by: agent-S-0332
  - to: in-progress
    at: 2026-10-07T23:08:59Z
    by: agent-S-0332
  - to: done
    at: 2026-10-07T23:12:36Z
    by: agent-S-0332
stream: S-0332
tags: [flai, template]
touches: [design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/agent-coordination.md, design/system/continuous-improvement.md]
after: [T-1195, T-1196, T-1197]
usage:
  source: log
  seconds: 217
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 62
      output: 24393
      cache_read: 3732899
      cache_write: 118350
      cost: 1.9444
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
