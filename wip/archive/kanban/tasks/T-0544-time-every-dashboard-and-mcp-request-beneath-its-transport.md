---
id: T-0544
type: task
nature: research
title: Time every dashboard and MCP request beneath its transport
status: done
parent: S-0152
owner: alex
created: 2026-09-29T06:46:36Z
updated: 2026-09-29T06:50:41Z
transitions:
  - to: ready
    at: 2026-09-29T06:46:43Z
    by: agent-S-0152
  - to: in-progress
    at: 2026-09-29T06:46:43Z
    by: agent-S-0152
  - to: done
    at: 2026-09-29T06:50:41Z
    by: agent-S-0152
stream: S-0152
tags: []
touches: [flai/internal/perf, flai/internal/channel, flai/internal/mcpserver]
usage:
  source: log
  seconds: 238
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 41
      output: 14279
      cache_read: 3726080
      cache_write: 41344
      cost: 1.3617
---
# T-0544 Time every dashboard and MCP request beneath its transport

## Work

Add `flai/internal/perf`: a per-request recorder carried in the context, with named phases. `channel.Client.call` and a receiving middleware on the MCP server start one for each request and log one event when it is answered: method (the tool for MCP `tools/call`), transport, duration in milliseconds, answer size, error, and the phases with their time and count. Requests slower than a threshold log at `info`, the rest at `debug`. Waits by design (`wait_for_work`, `wait_for_events`) are logged at `debug` whatever their length.

## Done when

- A dashboard request and an MCP tool call each produce one `request answered` event with the fields above, measured inside flai, so that the dashboard's own request duration minus flai's is the transport.
- The threshold is a setting documented in `docs/operators/settings.md`.
- Behaviour tests cover the recorder, the channel event, and the MCP event.

## Notes
