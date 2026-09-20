---
id: T-0258
type: task
nature: feature
title: "flai: a JSON-RPC peer over a WebSocket that flai opens, with the handshake, pings, reconnection, and project.info"
status: done
parent: S-0072
owner: alex
created: 2026-09-20T07:28:27Z
updated: 2026-09-20T07:34:50Z
transitions:
  - to: ready
    at: 2026-09-20T07:29:03Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T07:29:03Z
    by: system-flow
  - to: done
    at: 2026-09-20T07:34:50Z
    by: system-flow
stream: S-0072
tags: []
---
# T-0258 flai: a JSON-RPC peer over a WebSocket that flai opens, with the handshake, pings, reconnection, and project.info

## Work
A package in flai/internal for the host end of the channel: dial with the agent credential, prove and check possession over a nonce in the first exchange, JSON-RPC 2.0 requests, answers, notifications, and cancel, a method table with project.info checked against the project key, pings both ways, a size cap, reconnect with backoff and jitter. coder/websocket, recorded in design/tech.

## Done when
- Tests against an in-process server cover the handshake (good, wrong credential, impostor server), a method call, an unknown method, a wrong project key, a lost connection and the reconnect, and the cap

## Notes
