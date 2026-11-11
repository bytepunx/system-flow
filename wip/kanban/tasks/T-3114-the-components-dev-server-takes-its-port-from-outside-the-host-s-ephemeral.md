---
id: T-3114
type: task
nature: remediation
title: The components dev server takes its port from outside the host's ephemeral range, so no other process's port 0 listener can be given it
status: done
parent: S-0714
owner: alex
created: 2026-11-11T09:05:21Z
updated: 2026-11-11T18:27:03Z
transitions:
  - to: ready
    at: 2026-11-11T18:18:02Z
    by: agent-S-0714
  - to: in-progress
    at: 2026-11-11T18:18:02Z
    by: agent-S-0714
  - to: done
    at: 2026-11-11T18:27:03Z
    by: agent-S-0714
stream: S-0714
tags: [flai, components]
touches: [flai/internal/components/devserver.go, flai/internal/components/devserver_test.go]
usage:
  source: log
  seconds: 541
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 4166
      cache_read: 1013710
      cache_write: 52428
      cost: 0.619
---
# T-3114 The components dev server takes its port from outside the host's ephemeral range, so no other process's port 0 listener can be given it

## Work

I-0342's 404 body, `404 page not found`, is Go's `net/http` default, not Vite's: a Go server answered on the port vite had answered the harness on a moment before. `freePort` in `flai/internal/components/devserver.go` finds the port by listening on `127.0.0.1:0` and closing it, so the port comes from the kernel's ephemeral range, the same range every `httptest` server and every outgoing connection of the other close-outs running at the same time is given. Once vite stops listening, by exiting or by restarting its server, any of those can take the port, and with `--strictPort` vite cannot have it back.

Make `freePort` pick a port below the host's ephemeral range: on Linux the range is in `/proc/sys/net/ipv4/ip_local_port_range`, elsewhere take IANA's 49152 to 65535. Pick at random within a band below it, check the port is free by listening and closing as now, and keep `portTries` and the restart on a taken port. The kernel never gives a port outside its ephemeral range to a listen on port 0 or to an outgoing connection, so the other tests on the host can no longer take vite's port.

Waits for no task: it is the first layer.

## Done when

- `freePort` answers a port outside the host's ephemeral range, and a test in `flai/internal/components/devserver_test.go` checks that against the range it reads.
- The tests of `pickPort` and the taken-port restart in `devserver_test.go` still pass.
- `flai test flai/internal/components/devserver.go flai/internal/components/devserver_test.go` passes.

## Notes
