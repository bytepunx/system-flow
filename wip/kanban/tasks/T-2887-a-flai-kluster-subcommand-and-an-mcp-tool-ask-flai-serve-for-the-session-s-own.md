---
id: T-2887
type: task
nature: feature
title: A flai kluster subcommand and an MCP tool ask flai serve for the session's own story's deploy or teardown and print the answer
status: backlog
parent: S-0660
owner: alex
created: 2026-11-10T08:27:58Z
updated: 2026-11-10T08:27:58Z
transitions: []
stream: S-0660
tags: [cli]
touches: [flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/internal/mcpserver/kluster.go, flai/internal/mcpserver/kluster_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-2876]
---
# T-2887 A flai kluster subcommand and an MCP tool ask flai serve for the session's own story's deploy or teardown and print the answer

## Work

Add the subcommand and the MCP tool S-0659's ADR names: a `flai kluster` subcommand in `flai/cmd/kluster.go`, beside `deploy` and `teardown`, and a tool in a new `flai/internal/mcpserver/kluster.go`, registered in `folder.go` as `verify` is. Both call T-2876's `Ask` for a deploy, or a teardown with a flag, and print its answer: the outcome, the commit, each step's result, and the log's path, or that the job is still queued or running and how to read it later, with `--json` as the other kluster subcommands give it.

Both refuse, with the reason and what to do, a story other than the session's own (`FLAI_STORY`, else the agent's name, as the inbox finds it). flai serve answers the other refusals.

Add the tool to `TestToolsAreAdvertised`'s list in `server_test.go`. Run `make flai-reference`, which writes `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md`, so that `TestReferenceIsCurrent` and `TestReferenceCoversEveryCommand` pass.

It waits for T-2876, whose `Ask` it calls. It runs beside T-2884, which touches none of its files.

## Done when

- `flai/cmd/kluster_test.go` and `flai/internal/mcpserver/kluster_test.go` cover the own story answered, another story refused, an answer that outlasts the wait, and the JSON output.
- `TestToolsAreAdvertised`, `TestReferenceIsCurrent`, and `TestReferenceCoversEveryCommand` pass.
- `go test ./cmd/... ./internal/mcpserver/...` passes.

## Notes
