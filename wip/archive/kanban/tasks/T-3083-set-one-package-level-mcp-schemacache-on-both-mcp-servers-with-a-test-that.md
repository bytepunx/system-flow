---
id: T-3083
type: task
nature: improvement
title: Set one package-level mcp.SchemaCache on both MCP servers, with a test that their tool schemas are unchanged
status: done
parent: S-0721
owner: alex
created: 2026-11-11T05:28:22Z
updated: 2026-11-11T05:35:40Z
transitions:
  - to: ready
    at: 2026-11-11T05:30:02Z
    by: agent-S-0721
  - to: in-progress
    at: 2026-11-11T05:30:02Z
    by: agent-S-0721
  - to: done
    at: 2026-11-11T05:35:40Z
    by: agent-S-0721
stream: S-0721
tags: [mcp, tests]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/schema_cache_test.go]
usage:
  source: log
  seconds: 338
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 6196
      cache_read: 3564683
      cache_write: 90650
      cost: 1.5365
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 14
          cache_read: 94870
          cache_write: 22591
          cost: 0.0326
---
# T-3083 Set one package-level mcp.SchemaCache on both MCP servers, with a test that their tool schemas are unchanged

## Work

Declare one package-level `schemaCache = mcp.NewSchemaCache()` in `flai/internal/mcpserver/server.go` and set `ServerOptions.SchemaCache` to it in both `mcp.NewServer` calls: `New` in `server.go` and `newFolderServer` in `folder.go`. Waits for nothing: it is the first layer.

Before relying on the cache, check every tool for a schema that depends on its options rather than its Go type alone. `permission_prompt` in `folder.go` takes a pre-defined `permissionInputSchema` map; confirm how the SDK caches a pre-defined schema (by pointer, per `mcp/schema_cache.go`) and that it is never changed in place. Keep any option-dependent tool out of the cache, or stop and say so on a thread on S-0721.

Add `flai/internal/mcpserver/schema_cache_test.go`. It builds two servers from the same options and checks that `tools/list` answers the same tools with the same input and output schemas from both, and the same as a server built without the cache. It also checks that servers built from differing options (single project and folder, each role) list the schemas they listed before the change.

## Done when

- Both `mcp.NewServer` calls set `SchemaCache` to the one package-level cache.
- The new test passes under `-race`, and `scripts/flai-test.sh` is clean for `./internal/mcpserver`.
- Criteria 1 and 2 of S-0721 are ticked.

## Notes
