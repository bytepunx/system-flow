---
id: T-0674
type: task
nature: improvement
title: claude-code runs each role's model, and the command harness is told the roles
status: done
parent: S-0189
owner: arobson
created: 2026-10-01T10:53:32Z
updated: 2026-10-01T11:05:23Z
transitions:
  - to: ready
    at: 2026-10-01T10:53:48Z
    by: agent-S-0189
  - to: in-progress
    at: 2026-10-01T11:03:25Z
    by: agent-S-0189
  - to: done
    at: 2026-10-01T11:05:23Z
    by: agent-S-0189
stream: S-0189
tags: []
touches: [flai/internal/harness]
usage:
  source: log
  seconds: 118
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 152
      cache_read: 3407186
      cache_write: 17971
      cost: 1.4828
---
# T-0674 claude-code runs each role's model, and the command harness is told the roles

## Work

When the story's agent has roles, start `claude` with `--agents`. Build it from the project's `.claude/agents/<definition>.md` for each role (`explore` is `explorer`, `verify` is `verifier`), with the role's model over the definition's. Refuse a role whose harness is not claude-code, which names config claude-code does not take, or whose definition is missing. The `command` harness gets `FLAI_AGENT_ROLES` as JSON.

## Done when

- Adapter tests cover the `--agents` JSON, the refusals, and the env.
- `flai-cli.md` and `docs/users/flai.md` say it.

## Notes
