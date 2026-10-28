---
id: T-2206
type: task
nature: remediation
title: mcphttp builds both SDK handlers with the localhost protection off, a test reproduces I-0213's 403 on a foreign Host, and flai mcp http's help says the token is the gate
status: backlog
parent: S-0515
owner: alex
created: 2026-10-28T13:02:26Z
updated: 2026-10-28T13:02:26Z
transitions: []
stream: S-0515
tags: [cli, mcp]
touches: [flai/internal/mcphttp/mcphttp.go, flai/internal/mcphttp/mcphttp_test.go, flai/cmd/mcp_http.go, docs/users/flai-reference.md]
---
# T-2206 mcphttp builds both SDK handlers with the localhost protection off, a test reproduces I-0213's 403 on a foreign Host, and flai mcp http's help says the token is the gate

## Work

In `flai/internal/mcphttp/mcphttp.go`, `Handler` builds `h.sessions` and `h.alone` with `mcp.StreamableHTTPOptions` that set `DisableLocalhostProtection: true`, with a comment saying why: flai refuses an `Origin` and checks the bearer token before the SDK sees a request, so the DNS-rebinding page the SDK's check guards against is refused already, and a container that dials the host under `host.docker.internal` or any other name is admitted once it carries the token (I-0213).

In `flai/internal/mcphttp/mcphttp_test.go`, add a test that reproduces the cause: `lab` listens on loopback, so a `POST` with the token and `Host: host.docker.internal:4243` (set `req.Host`, since `Host` is not a header Go sends from the map) answers the `initialize` with 200 only once the fix is in, and fails with 403 `invalid Host header` without it. In the same test, a request with that `Host` and an `Origin` still answers 403, and one with that `Host` and no token still answers 401, so the gates stand.

In `flai/cmd/mcp_http.go`, the `http` command's `Long` help says that the token is the gate: a request from a browser is refused, and any other `Host` the request names, such as a container's name for the host, is admitted with the token. Run `make flai-reference` so that `docs/users/flai-reference.md` carries the new help, since `reference_test` fails when it is stale.

Waits for nothing: it is the fix the ADR records and the documents describe.

## Done when

- `go test ./internal/mcphttp/` passes with the new test, and the test fails on the handler as it was before the change
- A `POST /mcp` with the token and a `Host` that is not loopback answers 200 on a loopback listener; one with an `Origin` answers 403 and one without the token 401, whatever the `Host`
- `flai mcp http --help` says the token is the gate and a foreign `Host` is admitted with it, and `go test ./cmd/ -run Reference` passes with the regenerated `docs/users/flai-reference.md`
- `flai test` on the paths changed passes

## Notes

The SDK's option is `StreamableHTTPOptions.DisableLocalhostProtection` (go-sdk v1.8.0, `mcp/streamable.go`); its check compares the listener's local address and `req.Host` with `util.IsLoopback`. The spike S-0505 got through only by sending `Host: 127.0.0.1:4243` from `mcp.json`, which Claude Code honours and Node's fetch does not.
