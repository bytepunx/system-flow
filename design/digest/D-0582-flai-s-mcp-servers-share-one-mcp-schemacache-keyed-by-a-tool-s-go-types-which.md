---
id: D-0582
title: "flai's MCP servers share one mcp.SchemaCache keyed by a tool's Go types, which holds only while every cached tool's schema is its handler's types': a tool with a schema of its own is added with srv.AddTool, as permission_prompt is"
kind: decision
sources: [S-0721, T-3083, T-3085]
topics: [mcp, go, code, testing]
paths: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/schema_cache_test.go, design/tech/go-libraries.md]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0582 flai's MCP servers share one mcp.SchemaCache keyed by a tool's Go types, which holds only while every cached tool's schema is its handler's types': a tool with a schema of its own is added with srv.AddTool, as permission_prompt is

S-0721 made every MCP server flai builds share one package-level `mcp.SchemaCache` (`schemaCache` in `flai/internal/mcpserver/server.go`, passed by `Options.schemas()` to both servers). `design/tech/go-libraries.md`'s go-sdk row records the cache. What it leaves out is the condition the cache rests on, which the story's agent checked before relying on it (`wip/archive/agents/S-0721.md`, Decisions):

```text
Schemas depend on the Go type alone: every flai tool but permission_prompt registers with mcp.AddTool and no InputSchema or OutputSchema, so the SDK derives both from the handler's types (setSchema caches by reflect.Type). permission_prompt is added with the low-level srv.AddTool and a map schema, which never reaches the cache. No flai code changes a tool's schema after it is added, so no tool needs to be kept out of the cache.
```

So a new tool registered with `mcp.AddTool` gets its schemas from the cache. A tool whose schema is not what its handler's types give, or that changes its schema after it is added, must not go through the cache: add it with `srv.AddTool` and a schema of its own, as `permission_prompt` is.

A test that needs an uncached server sets the unexported `Options.noSchemaCache` field. It is a field, not a package variable, so the tests stay parallel and set nothing process-wide, as D-0571 requires.

T-3085 measured `./internal/mcpserver` with `-short -count=1` at the story's base and after T-3083, on a 24-core host at load 41 to 82. Under `-race` the user CPU fell from about 544 to 566 s to about 105 s, 19% of the before figure. Without `-race` it fell from 46.8 s to about 24 s, 52%.

Sources: S-0721 (T-3083, T-3085).
