---
id: T-0278
type: task
nature: feature
title: "flaiover: /mcp, the bridge, flai.ts, and their settings go; /mcp says where MCP lives now"
status: done
parent: S-0076
owner: alex
created: 2026-09-20T12:16:19Z
updated: 2026-09-20T12:27:43Z
transitions:
  - to: ready
    at: 2026-09-20T12:25:32Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T12:25:32Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:27:43Z
    by: system-flow
stream: S-0076
tags: []
---
# T-0278 flaiover: /mcp, the bridge, flai.ts, and their settings go; /mcp says where MCP lives now

## Work
Remove routes/mcp, mcpbridge.ts, flai.ts, FLAIOVER_MCP_* settings, and the bearer-only rule for /mcp in auth.ts. A request to /mcp, any method, answers 410 with a JSON-RPC error that names flai mcp and flai mcp start on the host. flai dashboard stops passing anything that existed for the bridge.

## Done when
- Route test for 410 with and without a token; no-project-files allow-list shrinks
- make flaiover-test and make flaiover-build pass

## Notes
