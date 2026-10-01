---
id: T-0641
type: task
nature: remediation
title: flai serve keeps no MCP servers for a host its config does not name, and says why
status: ready
parent: S-0183
owner: arobson
created: 2026-10-01T08:34:18Z
updated: 2026-10-01T08:34:25Z
transitions:
  - to: ready
    at: 2026-10-01T08:34:25Z
    by: agent-S-0183
stream: S-0183
tags: []
touches: [flai/cmd/serve.go, flai/internal/host]
---
# T-0641 flai serve keeps no MCP servers for a host its config does not name, and says why

## Work
Under a host (FLAI_HOST_URL and FLAI_HOST_TOKEN set), flai serve tells it which MCP servers to keep only when the host's folder beside its own config holds that host's token. Otherwise it keeps none, as outside a host, and logs why. Its registry and dashboard connection already follow its config.

## Done when
A test shows a serve with another --config does not take the host in its environment, with the reason logged; the tests pass.

## Notes
