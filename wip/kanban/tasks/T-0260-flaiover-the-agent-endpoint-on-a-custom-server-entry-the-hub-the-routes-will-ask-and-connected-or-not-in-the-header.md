---
id: T-0260
type: task
nature: feature
title: "flaiover: the agent endpoint on a custom server entry, the hub the routes will ask, and connected or not in the header"
status: backlog
parent: S-0072
owner: alex
created: 2026-09-20T07:28:28Z
updated: 2026-09-20T07:28:28Z
transitions: []
stream: S-0072
tags: []
---
# T-0260 flaiover: the agent endpoint on a custom server entry, the hub the routes will ask, and connected or not in the header

## Work
A custom server entry around adapter-node's handler and the same hook in the development server; the agent endpoint accepts one WebSocket per credential, refuses a browser's Origin and a wrong credential, answers the handshake; a hub in server code pairs requests with answers and marks a silent flai gone; a status route and the header say connected or not with the command to start it. The image runs the custom entry.

## Done when
- Unit tests for the hub, the handshake, and the refusals
- make flaiover-test and make flaiover-build pass

## Notes
