---
id: T-0902
type: task
nature: feature
title: flai thread reply, the MCP tool thread_reply, and the host API post a recommendation with its source, and flai thread confirm makes it the answer
status: backlog
parent: S-0220
owner: alex
created: 2026-10-05T04:47:37Z
updated: 2026-10-05T04:47:54Z
transitions: []
stream: S-0220
tags: [flai]
touches: [flai/cmd/thread.go, flai/cmd/thread_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-0894]
---
# T-0902 flai thread reply, the MCP tool thread_reply, and the host API post a recommendation with its source, and flai thread confirm makes it the answer

## Work

Expose T-0894's marks and `Confirm` wherever a thread is written:

- `flai thread reply` (`flai/cmd/thread.go`, line 74) takes `--recommend` and `--source <path>[#<heading>]`. A source must name a file in the repository, and a heading must exist in it; refuse one that does not, naming what was wrong.
- A new `flai thread confirm <TH-n>` runs `threads.Confirm` as the caller, with `--by` as `thread reply` takes it.
- The MCP tool `thread_reply` (`threadReply`, `flai/internal/mcpserver/server.go` line 359) takes `recommendation` and `source` with the same checks. A new tool, `thread_confirm`, is not added: confirming is the operator's, from the CLI or the dashboard.
- When the caller's role is `orchestrate`, `thread_reply` records the reply in the orchestrator's decision log, which S-0218 adds: the thread, whether it was a recommendation or an answer, and the source. The log then never depends on the model remembering to call `activity_log`.
- The host API's `thread.reply` (`flai/internal/hostapi/writes.go`, line 1123) passes `recommend` and `source` through. A new `thread.confirm` runs `flai thread confirm` as the manifest's owner, as `thread.reply` does.
- Describe the flags, the command, the tool's parameters, and the method in `design/system/flai-cli.md`, `docs/users/flai.md` (`### Threads`), and `docs/users/flai-reference.md`.

It waits for T-0894, whose marks and `Confirm` it calls. It runs with T-0906, whose paths it does not share.

## Done when

- a CLI test posts a recommendation with a source, refuses a source that names no file or no heading, and confirms it
- an MCP test posts a recommendation through `thread_reply` and finds the thread's status unchanged, and an orchestrate-role reply in the decision log with its source
- a host API test confirms a thread through `thread.confirm` as the owner
- the three documents describe the flags, the command, the parameters, and the method
- `go test ./cmd/ ./internal/mcpserver/ ./internal/hostapi/` passes

## Notes
