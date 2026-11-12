---
id: T-3144
type: task
nature: improvement
title: mcp.bind holds one address, several, or auto, and a configuration with one address reads as before
status: backlog
parent: S-0732
owner: alex
created: 2026-11-12T15:54:16Z
updated: 2026-11-12T15:54:16Z
transitions: []
stream: S-0732
tags: [cli]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go]
---
# T-3144 mcp.bind holds one address, several, or auto, and a configuration with one address reads as before

## Work

Let `MCPHost.Bind` in `flai/internal/config/config.go` hold one address or a list. Read a JSON string as a list of one and a JSON array as the list, and write a list of one back as a string, so a configuration written before this story reads and saves unchanged. Accept the value `auto`, alone, as the choice left to flai. Have `MCPBind` answer the list, with `DefaultMCPBind` when unset, or the default the plan's thread settles. Have `WithMCPBind` take several addresses, check each with `checkMCPBind`, refuse a repeated address, and refuse `auto` mixed with addresses. Keep `WithoutMCPBind` as it is.

This task waits for none: it is the first layer, with the gateway lookup.

## Done when

- `config_test.go` shows a string `bind` and an array `bind` both read, a list of one saved back as a string, and `auto` accepted alone and refused beside an address.
- A bad address in a list is refused, naming that address.
- `flai test flai/internal/config` passes.

## Notes
