---
id: S-0066
type: story
nature: feature
title: flai dashboard finds a free port in a block reserved for flai, and a project keeps its port
status: cancelled
parent: E-0007
owner: alex
created: 2026-09-20T01:18:58Z
updated: 2026-09-20T06:46:19Z
transitions:
  - to: ready
    at: 2026-09-20T03:31:57Z
    by: alex
  - to: cancelled
    at: 2026-09-20T06:46:19Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd]
---
# S-0066 flai dashboard finds a free port in a block reserved for flai, and a project keeps its port

## Goal
`flai dashboard` finds a free port for each project's dashboard by itself, from a block of ports reserved for flai that is unlikely to collide with common software, so the operator can run many projects on one host without choosing or remembering ports, and the same project comes back on the same port.

## Acceptance criteria
- [ ] flai has a default block of 50 to 100 ports, chosen with evidence that it avoids IANA-registered and commonly used open source ports, the Linux and IANA ephemeral ranges, and the ports this project already uses (4242 today); the block is a host setting (`dashboard.port_range` or the like) and the choice is documented with its reasons
- [ ] Choosing a port: an explicit `--port`, a manifest port, or a configured single port still wins and is used as given, failing clearly when it is taken; otherwise the port this project's record (S-0065) last held is reused when it is free; otherwise the lowest port in the block that no other recorded instance holds and that is actually free, which flai proves by binding it on the address it will publish on
- [ ] Two `flai dashboard` commands started at the same moment for different projects never get the same port: allocation and the write of the record happen under a lock in flai's home, and a test starts them together
- [ ] A port held by a record whose container is gone is reclaimable; a port in the block that something else on the host has taken is skipped, not fought over; when the block is exhausted flai says so, names what holds each port, and starts nothing
- [ ] The chosen port is written to the instance record before the container is started and corrected if the start fails, so a reader never sees a port that belongs to nothing
- [ ] flaiover's part is stated: the container always listens on its fixed internal port and only the published host port varies, so nothing changes in the image unless this story finds otherwise, and it says so
- [ ] Existing set-ups keep working: a project that has always used 4242, or sets a port in its manifest or config, is not moved; `docs/users/flai.md`, `docs/operators/index.md`, and `flai-cli.md` describe the block, the order of precedence, and how to pin a port
- [ ] Tests with the fake runner, a scratch home, and real sockets for the free-port probe: first start, reuse, collision with another record, collision with a foreign listener, exhaustion, explicit port taken, concurrent starts

## Tasks

## Notes
From E-0007, the operator's second bullet: "flai/flaiover will need to negotiate a port depending on what's available. this can be part of the metadata read/write pattern to a ~/.flai folder such that flai reserves a port range block of ~50-100 ports that are unlikely to collide with common OSS/system ports."

What exists. The port comes from `--port`, then `dashboard.port` in `system-flow.yaml`, then the host config, default 4242. A second project on the same host fails at `docker run` with the port already allocated, and the operator picks another by hand. The container listens on 3000 inside and `flai dashboard` publishes `<bind>:<port>:3000`.

Depends on S-0065: the record is where a project's port is remembered and where the other projects' ports are read from. Do that one first.

Ranges to stay out of: 0 to 1023 (system), the Linux default ephemeral range 32768 to 60999 and IANA's 49152 to 65535 (a free port there can be taken by any outgoing connection a moment later), and well-known development ports (3000, 5173, 8000, 8080, 8888, 9000 and their neighbours). Whether the block should start at 4242, so the first project stays where operators already look, is part of the choice.

The probe has to bind on the address the dashboard will be published on (0.0.0.0 by default since S-0035), because a port free on loopback can be taken on another interface.
- 2026-09-20T06:46:19Z: moved to cancelled: E-0007 cancelled: going a different direction
