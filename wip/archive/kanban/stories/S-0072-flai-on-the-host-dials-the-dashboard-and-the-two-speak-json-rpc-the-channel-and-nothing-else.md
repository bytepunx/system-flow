---
id: S-0072
type: story
nature: feature
title: "flai on the host dials the dashboard and the two speak JSON-RPC: the channel and nothing else"
status: done
parent: E-0003
owner: alex
created: 2026-09-20T07:26:49Z
updated: 2026-09-20T07:52:32Z
transitions:
  - to: ready
    at: 2026-09-20T07:27:39Z
    by: alex
  - to: in-progress
    at: 2026-09-20T07:27:52Z
    by: system-flow
  - to: review
    at: 2026-09-20T07:50:17Z
    by: system-flow
  - to: done
    at: 2026-09-20T07:52:32Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, flaiover/Dockerfile, design/system]
---
# S-0072 flai on the host dials the dashboard and the two speak JSON-RPC: the channel and nothing else

## Goal
One flai process per user on the host (`flai serve`, the name to be settled here) opens a WebSocket to the dashboard's agent endpoint and keeps it open; the two speak JSON-RPC 2.0. This story builds the channel and proves it with one method; the clone stays mounted and the dashboard works exactly as before, with or without a flai connected.

## Acceptance criteria
- [x] `flai serve` dials the dashboard of each project registered with it, authenticates with a credential of its own (generated on the host, handed to the container as a secret by `flai dashboard`, accepted only on the agent endpoint, never the browser's token, never logged), and both sides prove they hold it in the first exchange
- [x] flaiover has the agent endpoint: a custom server entry around adapter-node's handler, the same hook in the development server, one connection per credential, requests whose `Origin` is a site refused; the image runs the custom entry and `make flaiover-test` covers the endpoint
- [x] JSON-RPC 2.0 framing with requests, answers, and notifications in both directions, ids pairing answers, a cancel notification, a cap on message size where `maxBuffer` is today, and one method, `project.info`, answered from the manifest and checked against the project key on every request (ADR-0024)
- [x] Both sides ping; a flai that is frozen or gone is marked so within ten seconds; flai reconnects with backoff and jitter from a quarter of a second to a few seconds, and a replaced container is reconnected within a second of listening, as the trials saw
- [x] `flai dashboard` starts `flai serve` detached when it is not running, without root, and registers the project with it (a list in flai's home, honouring an override for tests so the real one is never touched); `flai dashboard stop` leaves it running for other projects; `flai serve` in the foreground works for watching it
- [x] `flai dashboard status` says whether the host flai runs, since when, which projects it serves, and whether the dashboard has it connected; the dashboard shows connected, or not connected with the command to start it, in its header
- [x] Works on Linux and WSL2 and is tried there; macOS and Windows with Docker Desktop are reasoned about and marked tried or not
- [x] `design/system/flai-cli.md`, `flaiover-dashboard.md`, `design/tech`, and the users' and operators' documentation describe the channel; the WebSocket libraries on both sides are recorded in `design/tech`

## Tasks
- T-0258 flai: a JSON-RPC peer over a WebSocket that flai opens, with the handshake, pings, reconnection, and project.info
- T-0259 flai serve: one process per user, the projects registered with it, started detached, and its state readable
- T-0260 flaiover: the agent endpoint on a custom server entry, the hub the routes will ask, and connected or not in the header
- T-0261 flai dashboard hands the container the agent credential, starts flai serve, registers the project, and reports it in status
- T-0262 Try it end to end on this host, reason about the other platforms, and document the channel

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. The trials used `coder/websocket` in Go and `ws` in Node. SSE down with POST up carries the same methods and is the fallback if the custom server entry proves fragile. No dependency on other stories.
