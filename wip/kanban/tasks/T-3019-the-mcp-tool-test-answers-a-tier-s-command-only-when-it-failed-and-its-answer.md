---
id: T-3019
type: task
nature: remediation
title: The MCP tool test answers a tier's command only when it failed, and its answer fits 40,000 bytes
status: backlog
parent: S-0707
owner: alex
created: 2026-11-10T10:32:03Z
updated: 2026-11-10T10:32:03Z
transitions: []
stream: S-0707
tags: [flai, mcp]
touches: [flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go]
---
# T-3019 The MCP tool test answers a tier's command only when it failed, and its answer fits 40,000 bytes

## Work

In `flai/internal/mcpserver/test.go`, project the `verify.Result` that `verify.Test` returns before the handler answers it, rather than returning it as it is. Each tier keeps its name, state, duration, findings, omitted count, and log. The tier that failed keeps its `command`. A tier that passed, or was not reached, leaves `command` out and gives `files` in its place: the number of files the tier was given. When the encoded answer would still pass 40,000 bytes (`context.PartLimit`, measured as the MCP go-sdk encodes it, as `inbox.Fit` and `taskdone.Result.Fit` do), cut `paths` to its first entries and count the rest in `paths_omitted`. Leave `verify.Result`, `flai test --json`, and the verify record as they are. Rewrite the tool's description (`testDescription`) to say what the answer leaves out, and that `flai test --json` gives every tier's whole command.

In `flai/internal/mcpserver/test_test.go`, reproduce I-0260: a manifest tier whose command takes `{files}`, run on enough files that today's answer passes 40,000 bytes, every tier passing. Assert that the answer fits 40,000 bytes, that no passed tier has a `command`, and that each gives `files`. Add a failing run over the same files, and assert that the failed tier keeps its whole `command`. Keep the existing tests passing. This is the first layer: it waits for no task.

## Done when

- The reproducing test fails on the code before the change and passes after it.
- A passing run's answer has no tier `command` and fits 40,000 bytes.
- A failing run's answer keeps the failed tier's whole `command`.
- `flai test --json` prints what it printed before.
- `flai test` passes on the changed files.

## Notes
