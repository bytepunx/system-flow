---
id: T-3240
type: task
nature: feature
title: strategic-agents.md, the roles READMEs, and the operators' and users' guides list the move condition in place of to-review and to-backlog
status: done
parent: S-0748
owner: alex
created: 2026-11-12T16:08:54Z
updated: 2026-11-15T21:04:39Z
transitions:
  - to: ready
    at: 2026-11-15T20:40:52Z
    by: agent-S-0748
  - to: in-progress
    at: 2026-11-15T20:40:52Z
    by: agent-S-0748
  - to: done
    at: 2026-11-15T21:04:39Z
    by: agent-S-0748
stream: S-0748
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/operators/index.md, docs/users/flai.md, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go]
after: [T-3231]
usage:
  source: log
  seconds: 1427
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 7806
      cache_read: 2340893
      cache_write: 81396
      cost: 1.0987
---
# T-3240 strategic-agents.md, the roles READMEs, and the operators' and users' guides list the move condition in place of to-review and to-backlog

## Work

Waits for T-3231: it documents the condition that task adds and what became of the old names. Change each mention of `to-review` and `to-backlog` in these documents to the new condition, with what it takes and what happens to a file that still names the old ones.

## Done when

None of these documents names `to-review` or `to-backlog` but as the old names, and each lists the new condition.

## Notes
