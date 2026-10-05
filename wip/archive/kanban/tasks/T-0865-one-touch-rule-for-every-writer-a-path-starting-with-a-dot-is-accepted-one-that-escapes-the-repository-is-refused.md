---
id: T-0865
type: task
nature: remediation
title: "One touch rule for every writer: a path starting with a dot is accepted, one that escapes the repository is refused"
status: done
parent: S-0260
owner: alex
created: 2026-10-05T03:14:15Z
updated: 2026-10-05T03:21:04Z
transitions:
  - to: ready
    at: 2026-10-05T03:14:33Z
    by: agent-S-0260
  - to: in-progress
    at: 2026-10-05T03:14:33Z
    by: agent-S-0260
  - to: done
    at: 2026-10-05T03:21:04Z
    by: agent-S-0260
stream: S-0260
tags: []
touches: [flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/internal/itemedit/itemedit.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/edit_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go]
usage:
  source: log
  seconds: 391
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 79
      output: 19759
      cache_read: 3485246
      cache_write: 106338
      cost: 1.7989
---
# T-0865 One touch rule for every writer: a path starting with a dot is accepted, one that escapes the repository is refused

## Work

I-0071's cause: `itemedit.cleanList` (flai edit, MCP `item_edit`) and `hostapi`'s `listValue` (the dashboard's new and edit) check touches with the tags' pattern, which wants a letter or digit first, so `.claude/agents/planner.md` is refused; `flai touches` and `flai story new --touches` check nothing. Add one rule in `workitem`, `CleanTouches`, used by every writer: trim, drop a trailing slash, empties, and repeats; accept a leading dot; refuse a comma, a leading dash, an absolute path, a `..` segment, and a control character, each with what to do. Tags keep their own rule. Waits for nothing.

## Done when

- [x] `flai edit --touches`, MCP `item_edit`, `flai touches`, `flai story new` and `flai task new --touches`, and the host API's new and edit accept `.claude/agents/planner.md` and refuse `../x`, `/etc`, `a,b`, and `-x`
- [x] A test fails without the change for `flai edit` and `item_edit` with a dot path, and the package tests pass

## Notes
