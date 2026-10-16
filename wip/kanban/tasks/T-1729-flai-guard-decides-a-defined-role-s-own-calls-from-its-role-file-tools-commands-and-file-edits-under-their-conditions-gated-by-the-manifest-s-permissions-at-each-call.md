---
id: T-1729
type: task
nature: feature
title: "flai guard decides a defined role's own calls from its role file: tools, commands, and file edits under their conditions, gated by the manifest's permissions at each call"
status: backlog
parent: S-0416
owner: alex
created: 2026-10-16T05:44:01Z
updated: 2026-10-16T05:50:22Z
transitions: []
stream: S-0416
tags: [cli]
touches: [flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/guard/call.go, flai/cmd/guard.go, docs/users/flai-reference.md]
---
# T-1729 flai guard decides a defined role's own calls from its role file: tools, commands, and file edits under their conditions, gated by the manifest's permissions at each call

## Work

- Add `flai/internal/guard/roles.go`: given the role's file from `roles.Load` and the manifest's permissions read at the call, it decides an MCP tool, a flai command, or a file edit. It passes the reads every sub-agent has, then each listed call under its named conditions and its permission. It refuses the rest with the role's name, the rule, and the advice to ask with `thread_open`.
- Route `Decide` in `call.go` to it for any `FLAI_ROLE` that is not `plan`, `orchestrate`, or `analyze`.
- Have `flai/cmd/guard.go` load the project's roles.
- Fail closed: a missing or unparsable file refuses every write, naming the file and `flai check`.
- Test a fixture role with a permission on and off between two calls.

First layer: the built-in roles move onto it in T-1730 and T-1731.

## Done when

`go test ./internal/guard/... ./cmd/...` passes. The existing guard tests are unchanged.

## Notes
