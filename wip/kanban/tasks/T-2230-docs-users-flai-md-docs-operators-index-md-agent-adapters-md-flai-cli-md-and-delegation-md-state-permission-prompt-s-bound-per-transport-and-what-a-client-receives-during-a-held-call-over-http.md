---
id: T-2230
type: task
nature: improvement
title: docs/users/flai.md, docs/operators/index.md, agent-adapters.md, flai-cli.md, and delegation.md state permission_prompt's bound per transport and what a client receives during a held call over HTTP
status: backlog
parent: S-0519
owner: alex
created: 2026-10-28T13:45:08Z
updated: 2026-10-28T13:45:36Z
transitions: []
stream: S-0519
tags: [docs, mcp]
touches: [docs/users/flai.md, docs/operators/index.md, design/system/agent-adapters.md, design/system/flai-cli.md, design/conventions/delegation.md]
after: [T-2229]
---
# T-2230 docs/users/flai.md, docs/operators/index.md, agent-adapters.md, flai-cli.md, and delegation.md state permission_prompt's bound per transport and what a client receives during a held call over HTTP

## Work

Each place that states the four-minute bound says the bound per transport as the ADR settles it: `docs/users/flai.md` § Writes to protected paths ("held at most four minutes, Claude Code's bound") and the `permission_prompt` row of its MCP tools table; `docs/operators/index.md` § MCP over HTTP, the "Long requests" row, which says a progress notification goes every minute to a client that asked, where the server answers as JSON and a notification reaches the client only when the call ends, and which gains Claude Code's measured tool-call timeout over HTTP beside the 30-minute wait; `design/system/agent-adapters.md` § Asking the owner before a protected write, the Claude Code row of the bound table, and § A container harness of flai's own, the sentence of "What the spike settled" that names the hold outliving the timeout; `design/system/flai-cli.md`, the `flai mcp` row, where `permissionWait` is "below Claude Code's idle timeout for an MCP call, 30 minutes over stdio and 5 over HTTP"; and `design/conventions/delegation.md` § Project additions, the sentence on the hold. Bump `updated` on each design file. It waits for the code task: it describes what is built.

## Done when

No document states four minutes as the bound over HTTP, each named section gives the bound per transport and the measured timeout, and the markdown lint and `flai check --strict` pass.

## Notes

The template's copies of the conventions do not state the bound, so `template/root` is not touched.
