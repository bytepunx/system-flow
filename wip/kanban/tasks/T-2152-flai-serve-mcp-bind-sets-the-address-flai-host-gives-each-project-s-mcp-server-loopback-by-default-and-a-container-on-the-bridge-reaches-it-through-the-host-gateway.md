---
id: T-2152
type: task
nature: feature
title: flai serve mcp bind sets the address flai host gives each project's MCP server, loopback by default, and a container on the bridge reaches it through the host gateway
status: backlog
parent: S-0506
owner: alex
created: 2026-10-26T04:47:30Z
updated: 2026-10-26T12:30:53Z
transitions: []
stream: S-0506
tags: [cli, mcp]
touches: [flai/cmd/host.go, flai/cmd/host_test.go, flai/cmd/serve_mcp.go, flai/cmd/serve_mcp_test.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go]
---
# T-2152 flai serve mcp bind sets the address flai host gives each project's MCP server, loopback by default, and a container on the bridge reaches it through the host gateway

## Work

Add a host setting `mcp.bind` to `flai/internal/config/config.go`, the address (without port) `flai host` binds each project's `flai mcp http` to, `127.0.0.1` when unset. Add `flai serve mcp bind <address>` to `cmd/serve_actions.go`, beside the other `flai serve` settings, and show it in `flai serve settings`. Make `hostLauncher.mcpAddr` in `cmd/host.go` use the setting in place of the loopback address when choosing each project's address, keeping the per-project port choice as it is. Take the address the spike's results found reached the listener from a container on this WSL2 host (Docker's bridge gateway, `0.0.0.0`, or the host network) as the documented example. The `Origin` refusal and the bearer check in `mcphttp` stand unchanged. This task waits for nothing: it changes no file the other first-layer tasks change.

## Done when

With `mcp.bind` set to the bridge gateway's address, `flai host` starts the project's MCP server there, `docker run --rm --add-host host.docker.internal:host-gateway curlimages/curl` reaches `/mcp` with the token and is refused without it, and `flai test flai/cmd flai/internal/config` passes.

## Notes

`docs/operators/settings.md` and `settings_doc_test` are the docs task's; this one adds the key and the command.
