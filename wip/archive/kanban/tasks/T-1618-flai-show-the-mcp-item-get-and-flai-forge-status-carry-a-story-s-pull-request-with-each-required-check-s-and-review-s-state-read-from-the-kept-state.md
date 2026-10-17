---
id: T-1618
type: task
nature: feature
title: flai show, the MCP item_get, and flai forge status carry a story's pull request with each required check's and review's state, read from the kept state
status: done
parent: S-0394
owner: alex
created: 2026-10-16T04:36:54Z
updated: 2026-10-17T16:30:37Z
transitions:
  - to: ready
    at: 2026-10-17T16:15:37Z
    by: agent-S-0394
  - to: in-progress
    at: 2026-10-17T16:15:38Z
    by: agent-S-0394
  - to: done
    at: 2026-10-17T16:30:37Z
    by: agent-S-0394
stream: S-0394
tags: [cli]
touches: [flai/cmd/items.go, flai/cmd/show_test.go, flai/cmd/forge.go, flai/cmd/forge_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go]
after: [T-1615]
usage:
  source: log
  seconds: 899
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 49
      output: 20698
      cache_read: 2963375
      cache_write: 89520
      cost: 1.5134
---
# T-1618 flai show, the MCP item_get, and flai forge status carry a story's pull request with each required check's and review's state, read from the kept state

## Work

- In `flai/cmd/items.go`, have `flai show S-nnnn` print a pull request section when the kept state has one:
  - the URL
  - one line per required check: its name, its state, and its log link
  - one line per review: its author and its state
- `--json` adds a `pull_request` object.
- In `flai/internal/mcpserver/server.go`, have `itemOut` carry the same `pull_request` for `item_get`. Add it to the tool's description in `folder.go`, so that an agent knows to look there.
- In `flai/cmd/forge.go`, add the checks and reviews to `flai forge status`.
- Read only T-1615's state, never the forge.
- Test with the state present and absent.
- Waits for T-1615: it reads the kept state.

## Done when

The three surfaces show the pull request, its checks, and its reviews from the kept state, show nothing when there is none, and the tests pass.

## Notes
