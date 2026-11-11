---
id: T-3084
type: task
nature: improvement
title: Say in go-libraries.md's go-sdk row that flai shares one SchemaCache
status: done
parent: S-0721
owner: alex
created: 2026-11-11T05:28:25Z
updated: 2026-11-11T05:36:05Z
transitions:
  - to: ready
    at: 2026-11-11T05:30:02Z
    by: agent-S-0721
  - to: in-progress
    at: 2026-11-11T05:30:02Z
    by: agent-S-0721
  - to: done
    at: 2026-11-11T05:36:05Z
    by: agent-S-0721
stream: S-0721
tags: [mcp]
touches: [design/tech/go-libraries.md]
usage:
  source: log
  seconds: 363
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2175
      cache_read: 1251286
      cache_write: 31820
      cost: 0.5394
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 14
          cache_read: 94870
          cache_write: 22591
          cost: 0.0326
---
# T-3084 Say in go-libraries.md's go-sdk row that flai shares one SchemaCache

## Work

In `design/tech/go-libraries.md`, extend the `github.com/modelcontextprotocol/go-sdk` row's use column to say that every MCP server flai builds in one process shares one `mcp.SchemaCache`, so each tool's schema is resolved once per process, and name its trade-off: pre-defined schemas are cached by pointer, so none may be changed in place. Bump the document's `updated`. The version stays v1.8.0. Waits for nothing: it shares no path with the code task, so it runs in the first layer with it.

## Done when

- The go-sdk row names the shared `SchemaCache`; no dependency is added or upgraded.
- The markdown lint and `flai check --strict` are clean for the file.
- Criterion 4 of S-0721 is ticked.

## Notes
