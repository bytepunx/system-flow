---
id: T-0715
type: task
nature: improvement
title: The design and the user guide describe roles as ADR-0068 means them
status: done
parent: S-0196
owner: arobson
created: 2026-10-02T23:29:19Z
updated: 2026-10-02T23:41:37Z
transitions:
  - to: ready
    at: 2026-10-02T23:29:45Z
    by: agent-S-0196
  - to: in-progress
    at: 2026-10-02T23:37:19Z
    by: agent-S-0196
  - to: done
    at: 2026-10-02T23:41:37Z
    by: agent-S-0196
stream: S-0196
tags: []
touches: [design/system/conventions.md, docs/users/conventions.md, docs/users/flai.md, design/system/agent-context.md, design/system/flai-cli.md]
usage:
  source: log
  seconds: 258
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 61
      output: 14952
      cache_read: 2771847
      cache_write: 63234
      cost: 1.3061
---
# T-0715 The design and the user guide describe roles as ADR-0068 means them

## Work

Rewrite the Roles section and its table in `design/system/conventions.md`, and the role passages in `docs/users/conventions.md`, `docs/users/flai.md` (~798), and `design/system/agent-context.md` (~250, the role pack), so they give the new meaning. `roles` lists every agent that reads a convention: `story`, `explore`, `verify`, plus `orchestrator`, `planner`, and `analyzer`, which check accepts. No roles means every agent reads it. The story's pack holds the conventions with no roles or with `story`, and a role pack those with no roles or with the role. The table lists the baseline's values as the designer chose them (TH-0071). Waits for nothing in code; its table needs the designer's answer on TH-0071.

## Done when

- None of the four files describes the old meaning (no roles meaning the story's agent alone, or the story's agent reading everything)
- The table in `design/system/conventions.md` matches the values the baseline task sets

## Notes
