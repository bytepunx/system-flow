---
id: T-0273
type: task
nature: feature
title: "flaiover: every route asks the channel; acceptance streams over it; one retry after a reconnect; the image drops flai, git, and ssh"
status: done
parent: S-0075
owner: alex
created: 2026-09-20T08:49:14Z
updated: 2026-09-20T09:03:22Z
transitions:
  - to: ready
    at: 2026-09-20T08:53:12Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T08:53:12Z
    by: system-flow
  - to: done
    at: 2026-09-20T09:03:22Z
    by: system-flow
stream: S-0075
tags: []
---
# T-0273 flaiover: every route asks the channel; acceptance streams over it; one retry after a reconnect; the image drops flai, git, and ssh

## Work
flai.ts stops spawning and asks the hub; each route maps its body to a method; acceptance relays progress notifications as NDJSON; a write that loses the connection is repeated once with the same request ID when flai returns; typed errors become 409, 422, 400; the Dockerfile drops the Go stage, git, and openssh; readiness and the board's writable flag follow the connection.

## Done when
- Route tests through flai hostapi and with a fake hub
- make flaiover-test and make flaiover-build pass; image size before and after recorded

## Notes
