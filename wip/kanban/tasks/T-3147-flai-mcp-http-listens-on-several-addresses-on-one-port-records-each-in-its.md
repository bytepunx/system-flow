---
id: T-3147
type: task
nature: improvement
title: flai mcp http listens on several addresses on one port, records each in its state, and warns only of an address beyond loopback and the container gateways
status: done
parent: S-0732
owner: alex
created: 2026-11-12T15:54:39Z
updated: 2026-11-12T17:55:55Z
transitions:
  - to: ready
    at: 2026-11-12T17:42:10Z
    by: agent-S-0732
  - to: in-progress
    at: 2026-11-12T17:42:10Z
    by: agent-S-0732
  - to: done
    at: 2026-11-12T17:55:55Z
    by: agent-S-0732
stream: S-0732
tags: [cli]
touches: [flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go]
after: [T-3145]
usage:
  source: log
  seconds: 825
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 69
      output: 27724
      cache_read: 4098312
      cache_write: 150246
      cost: 2.2635
---
# T-3147 flai mcp http listens on several addresses on one port, records each in its state, and warns only of an address beyond loopback and the container gateways

## Work

Let `flai mcp http` in `flai/cmd/mcp_http.go` take `--addr` more than once. Every address uses the same port, so that the container's `host.docker.internal:<port>` and the host's `127.0.0.1:<port>` name one server. Open one listener per address in `serveMCPHTTP`, serve the one handler on each, and close them all on a stop. If any address cannot be listened on, close the ones already open and fail, naming that address. Keep `url` and `addr` in the state file (`mcpState`) on the first address, loopback when it is among them, and add `addrs` with every address. Write the first address to `mcp-http.addr` as now, so a remembered port reads as before. Log the warning for listening beyond this machine for each address that is neither loopback nor a container network's gateway, using T-3145's check.

It waits for T-3145, whose gateway check the warning uses.

## Done when

- `mcp_http_test.go` starts a server on 127.0.0.1 and a second local address on one port and reaches it on both with the token. This is the criterion 3 test's host side.
- The state lists both addresses.
- The warning is logged for a LAN address and not for a gateway or loopback.
- A failed second listen leaves no listener open.
- `flai test` on the changed files passes.

## Notes
