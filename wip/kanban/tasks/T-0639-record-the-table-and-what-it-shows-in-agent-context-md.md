---
id: T-0639
type: task
nature: research
title: Record the table and what it shows in agent-context.md
status: done
parent: S-0188
owner: arobson
created: 2026-10-01T08:31:38Z
updated: 2026-10-01T09:25:00Z
transitions:
  - to: ready
    at: 2026-10-01T08:31:52Z
    by: agent-S-0188
  - to: in-progress
    at: 2026-10-01T09:18:51Z
    by: agent-S-0188
  - to: done
    at: 2026-10-01T09:25:00Z
    by: agent-S-0188
stream: S-0188
tags: []
touches: [design/system/agent-context.md]
usage:
  source: log
  seconds: 369
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 101
      output: 28050
      cache_read: 4439269
      cache_write: 100338
      cost: 2.1417
---
# T-0639 Record the table and what it shows in agent-context.md

## Work

Add the table and the findings to `design/system/agent-context.md § Sub-agents › Measured`. If delegation costs more than it saves, file a story or open a question saying what to change: a cheaper sub-agent model, fewer delegations, or none.

## Done when

`agent-context.md` holds the table and its reading, `flai check --strict` passes, and any follow-up is filed.

## Notes
- `agent-context.md § Sub-agents › Measured` rewritten in 575cd25. A verifier checked every number against the logs and found five claims that overreached; they are corrected in 35ea8ac. `scripts/lint-md.sh` reports 0 errors, and `flai check --strict` reports 0 errors and 9 warnings, all of them main's.
- Delegation costs more than it saves, as far as these runs show, so S-0189 (backlog) is filed: a cheaper sub-agent model, and one verifier run in place of the agent's own full-suite runs.
