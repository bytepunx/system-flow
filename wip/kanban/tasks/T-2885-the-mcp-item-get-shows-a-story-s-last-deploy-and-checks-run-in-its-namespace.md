---
id: T-2885
type: task
nature: feature
title: The MCP item_get shows a story's last deploy and checks run in its namespace, with their commits and outcomes
status: backlog
parent: S-0662
owner: alex
created: 2026-11-10T08:27:43Z
updated: 2026-11-10T08:27:43Z
transitions: []
stream: S-0662
tags: [cli]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go]
after: [T-2880, T-2881]
---
# T-2885 The MCP item_get shows a story's last deploy and checks run in its namespace, with their commits and outcomes

## Work

Add a `kluster` field to `ItemOut` in `flai/internal/mcpserver/server.go`, beside `pull_request`, filled for a story in the same block from the part T-2881 gives the verify report, or from T-2875's evidence when the story has no report yet. Leave it out while the action is off or no `deploy` is set. Say it in the tool's description.

It waits for T-2880, which changes `server.go` too, and for T-2881, whose record it shows.

## Done when

- `server_test.go` shows `item_get` with the field for a story deployed and checked at its head, and without it with the action off, as `TestItemGetCarriesTheKeptPullRequest` does for `pull_request`.
- `flai test` passes on the two files.

## Notes

- Criterion 3 of S-0662, its MCP half.
