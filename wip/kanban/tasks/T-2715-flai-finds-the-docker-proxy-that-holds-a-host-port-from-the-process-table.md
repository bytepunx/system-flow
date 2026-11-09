---
id: T-2715
type: task
nature: remediation
title: flai finds the docker-proxy that holds a host port from the process table
status: done
parent: S-0615
owner: alex
created: 2026-11-09T09:51:29Z
updated: 2026-11-09T10:06:01Z
transitions:
  - to: ready
    at: 2026-11-09T09:57:53Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T09:57:53Z
    by: system-flow
  - to: done
    at: 2026-11-09T10:06:01Z
    by: agent-S-0615
stream: S-0615
tags: [flai, dashboard]
touches: [flai/cmd/dashboard_porthold.go, flai/cmd/dashboard_porthold_test.go]
usage:
  source: log
  seconds: 488
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 22
      output: 140
      cache_read: 2715872
      cache_write: 22961
      cost: 0.9964
---
# T-2715 flai finds the docker-proxy that holds a host port from the process table

## Work

Add `flai/cmd/dashboard_porthold.go` with a lookup of what holds a host port, for the watch and the status commands to name (I-0243). Given a `/proc` root and a port, read each `<root>/<pid>/cmdline` (NUL-separated, readable by any user) and return the first process whose program's base name is `docker-proxy` and whose arguments carry `-host-port <port>`: its PID, its `-host-ip`, and its `-container-ip` and `-container-port`. Return nothing, and no error, when the root does not exist (macOS, Windows, a container without `/proc`), when a process vanishes or cannot be read mid-scan, or when none matches. Add a one-line description of a holder for the messages, such as `docker-proxy (pid 188220) for 172.17.0.2:3000`. Keep it a function of the root, so tests give it a fake tree; the app passes `/proc`.

Add the probe the watch needs beside it: whether the host port the dashboard was recorded with (`dashboardRecord.Publish`, `<bind>:<port>:3000`) accepts a TCP connection, with a wildcard bind probed at `127.0.0.1` as `probeHost` does, and a short timeout. Give it a seam on `app`, as `healthProbe` has, so that the watch's tests need no real socket.

This task waits for nothing: it is a new file the watch task calls.

## Done when

- `flai/cmd/dashboard_porthold_test.go` builds a fake `/proc` tree with a `docker-proxy -proto tcp -host-ip 0.0.0.0 -host-port 4242 -container-ip 172.17.0.2 -container-port 3000`, another `docker-proxy` on 8080, an unrelated process, a PID folder with no `cmdline`, and a non-PID entry, and finds the 4242 holder with its PID and addresses; finds none for a port no process names; and finds none, without error, for a root that does not exist.
- The port probe parses a `publish` of `0.0.0.0:4242:3000`, `127.0.0.1:4242:3000`, and `[::]:4242:3000`, and a test probes a listener it opens on `127.0.0.1`.
- `flai test flai/cmd/dashboard_porthold.go flai/cmd/dashboard_porthold_test.go` passes.

## Notes
