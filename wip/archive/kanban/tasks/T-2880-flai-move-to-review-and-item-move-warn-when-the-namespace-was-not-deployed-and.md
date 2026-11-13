---
id: T-2880
type: task
nature: feature
title: flai move to review and item_move warn when the namespace was not deployed and checked at the head, and the unticked-criteria warning says which
status: done
parent: S-0662
owner: alex
created: 2026-11-10T08:27:15Z
updated: 2026-11-13T14:27:20Z
transitions:
  - to: ready
    at: 2026-11-13T13:56:31Z
    by: agent-S-0662
  - to: in-progress
    at: 2026-11-13T13:56:31Z
    by: agent-S-0662
  - to: done
    at: 2026-11-13T14:27:20Z
    by: agent-S-0662
stream: S-0662
tags: [cli]
touches: [flai/cmd/move.go, flai/cmd/move_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/review_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go]
after: [T-2875]
usage:
  source: log
  seconds: 1849
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 48
      output: 18467
      cache_read: 3548192
      cache_write: 151945
      cost: 1.9087
---
# T-2880 flai move to review and item_move warn when the namespace was not deployed and checked at the head, and the unticked-criteria warning says which

## Work

After `TransitionAll` moves a story to review, in `flai/cmd/move.go` and in `itemMove` in `flai/internal/mcpserver/server.go`, read T-2875's evidence. While the `kluster` host action is on and a `deploy` is set:

- When the namespace is not deployed at the branch's head, or its checks did not pass there, add a warning that says which and names the two commands that deploy and check. The move still goes through.
- When the move warns of unticked criteria, add to that warning whether the namespace was deployed and checked at the head.

`Repo.Move` in `flai/internal/workitem/rules.go` stays as it is: `workitem` cannot see kluster, so the callers add the sentence. The dashboard's `item.move` runs the CLI and gets it from there.

It waits for T-2875, whose evidence it reads.

## Done when

- `move_test.go` covers: the action off (no new warning), deployed and checked at the head (no new warning), not deployed, deployed at an older commit, checks failed, and an unticked criterion with each state.
- `review_test.go` covers the same warnings through `item_move`.
- No case refuses the move.
- `flai test` passes on the four files.

## Notes

- Criteria 1 and 2 of S-0662.
