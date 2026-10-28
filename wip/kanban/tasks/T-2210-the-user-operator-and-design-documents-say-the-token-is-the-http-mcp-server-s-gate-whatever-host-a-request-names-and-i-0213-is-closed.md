---
id: T-2210
type: task
nature: remediation
title: The user, operator, and design documents say the token is the HTTP MCP server's gate whatever Host a request names, and I-0213 is closed
status: backlog
parent: S-0515
owner: alex
created: 2026-10-28T13:02:53Z
updated: 2026-10-28T13:02:53Z
transitions: []
stream: S-0515
tags: [docs, mcp]
touches: [docs/users/flai.md, docs/operators/index.md, design/system/flai-cli.md, design/issues/I-0213-flai-s-http-mcp-server-refuses-an-agent-in-a-container-with-403-invalid-host-header-since-the-mcp-sdk-s-localhost-protection-admits-only-a-loopback-host-on-a-loopback-listener.md, design/issues/summary.md]
after: [T-2206, T-2208]
---
# T-2210 The user, operator, and design documents say the token is the HTTP MCP server's gate whatever Host a request names, and I-0213 is closed

## Work

Say in each document what the ADR of T-2208 settles, linking it: in `docs/users/flai.md` § Serving agents over MCP, where the HTTP server is described, that a request is admitted whatever `Host` it names once it carries the token, so an agent in a container reaches the server through its runtime's name for the host, such as `host.docker.internal`; in `docs/operators/index.md` § MCP over HTTP, in the Authentication row of its table, that the `Origin` refusal and the token are the gates and the SDK's localhost `Host` check is off; and in `design/system/flai-cli.md`, in the `flai mcp http` row, beside the sentence on the `Origin` refusal, that the two SDK handlers are built with `DisableLocalhostProtection` and why. Bump each document's `updated`.

Then close the issue in the story's worktree with `flai issue close I-0213 --reason` naming the fix: the SDK's localhost protection off on both handlers, the token and the `Origin` refusal as the gates, and the ADR; `flai issue close` rewrites the issue's file and `design/issues/summary.md`.

Waits for T-2206, since the documents describe the behaviour its test proves, and for T-2208, since they link the ADR by its number.

## Done when

- The three documents name the two gates, the `Host` check off, and the ADR, and `make flai-test`'s markdown lint passes on them
- I-0213 is `closed` with a reason that names what fixed it, and `design/issues/summary.md` no longer lists it as open
- `flai check --strict` is clean

## Notes

The `settings_doc_test` is not touched: the story adds no setting, variable, or flag.
