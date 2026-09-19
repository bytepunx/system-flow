---
id: T-0200
type: task
nature: feature
title: An official MCP client test against a running dashboard
status: done
parent: S-0043
owner: alex
created: 2026-09-19T05:17:07Z
updated: 2026-09-19T05:24:02Z
transitions:
  - to: ready
    at: 2026-09-19T05:22:31Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T05:22:31Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:24:02Z
    by: system-flow
stream: S-0043
tags: []
touches: [flai/tests]
---

# T-0200 An official MCP client test against a running dashboard

## Work
Add `flai/tests/integration/mcp_http_test.go`: with the official Go SDK's `StreamableClientTransport` and a bearer token, connect to the dashboard named by `FLAIOVER_MCP_URL` and `FLAIOVER_TOKEN`, list the tools, call `inbox`, and call `board`; skip with a stated reason when the variables are not set, so the ordinary tiers are unaffected. Say in the file how to run it against `flai dashboard`.

## Done when
- The test compiles and skips in `make integration`
- It passes against a running dashboard in T-0202

## Notes
