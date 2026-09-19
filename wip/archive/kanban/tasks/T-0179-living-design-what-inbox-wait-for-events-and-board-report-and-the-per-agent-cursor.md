---
id: T-0179
type: task
nature: remediation
title: "Living design: what inbox, wait_for_events, and board report, and the per-agent cursor"
status: done
parent: S-0058
owner: alex
created: 2026-09-19T02:46:20Z
updated: 2026-09-19T02:47:56Z
transitions:
  - to: ready
    at: 2026-09-19T02:47:12Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:47:12Z
    by: alex
  - to: done
    at: 2026-09-19T02:47:56Z
    by: alex
stream: S-0058
tags: []
touches: [design/system]
---

# T-0179 Living design: what inbox, wait_for_events, and board report, and the per-agent cursor

## Work
Refine the living design, not a new ADR (ADR-0020 stands: files are the record, MCP is the agent's view). In `design/system/agent-narrative.md` and the `flai mcp` row of `flai-cli.md`, say: `inbox` reports threads, the stories ready to pull in pull order with whether the WIP limit allows a pull, and the changes others made since this agent last looked; changes are derived from the files (`transitions` with `at` and `by`, `blocked` intervals, the board's `order`), so nothing new is recorded in git; the cursor is a read marker per agent under `.flai-cache/mcp/`, a cache in the sense of the overview's principle: losing it only repeats or skips reports of changes, never state, because ready work is always listed; `wait_for_events` returns at once when the cursor is behind, returns events as well as paths, and advances the cursor; a `board` tool mirrors `flai board --json`.

## Done when
- Both documents describe the tools as they will be built
- `flai check --strict` is clean

## Notes
