---
id: T-2154
type: task
nature: feature
title: harness.Request gains the container fields and the claude-code adapter writes an HTTP MCP client configuration beside the stdio one
status: backlog
parent: S-0507
owner: alex
created: 2026-10-26T04:47:49Z
updated: 2026-10-26T04:47:49Z
transitions: []
stream: S-0507
tags: [cli, docker]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go]
---
# T-2154 harness.Request gains the container fields and the claude-code adapter writes an HTTP MCP client configuration beside the stdio one

## Work

In `flai/internal/harness/harness.go`, add to `Request` the fields a container run needs, each optional and left empty by `flai serve` today: `Inner` (the harness pilot runs), `MCP` with `URLEnv` and `TokenEnv` (the names of the variables holding the host's HTTP MCP server's URL and token), `Cache` (the mounted folder for sessions, configuration, and the MCP file), and `ConfigDir` (where the inner harness keeps its configuration, `CLAUDE_CONFIG_DIR` for Claude Code). In `adapters.go`, make the claude-code adapter's `mcpConfig` write, when `Request.MCP` is set, an HTTP server entry (`type: http`, the URL, and `headers` with `Authorization: Bearer` read from the variable and `X-Flai-Agent`) under `Request.Cache` in place of the stdio entry, pass `--mcp-config <file>` and `--strict-mcp-config`, and set `CLAUDE_CONFIG_DIR` in `Start.Env` when `ConfigDir` is set; the stdio shape and the host start stay as they are when the fields are empty. Test both shapes in `harness_test.go`, checking the token's value is never in the file or argv, only the variable's name. This task waits for nothing.

## Done when

`harness_test.go` covers a request with and without the container fields, the HTTP configuration names the variable and not a value, and `flai test flai/internal/harness` passes.

## Notes

Claude Code reads `${VAR}` in `.mcp.json` headers; check the installed version expands it, and if not, pilot writes the value into the file under the cache folder, which is mounted and never committed.
