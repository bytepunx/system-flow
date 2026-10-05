---
id: T-0867
type: task
nature: remediation
title: A write that grows an open story's claim into an in-progress story's says so and tells both stories as an overlapped change
status: done
parent: S-0244
owner: alex
created: 2026-10-05T03:15:27Z
updated: 2026-10-05T05:01:19Z
transitions:
  - to: ready
    at: 2026-10-05T04:41:44Z
    by: agent-S-0244
  - to: in-progress
    at: 2026-10-05T04:41:44Z
    by: agent-S-0244
  - to: done
    at: 2026-10-05T05:01:19Z
    by: agent-S-0244
stream: S-0244
tags: [flai, mcp]
touches: [flai/internal/workitem/hold.go, flai/internal/itemedit, flai/internal/itemnew, flai/cmd/items.go, flai/cmd/edit.go, flai/cmd/touches.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/internal/mcpserver/cursor.go, flai/cmd/touches_overlap_test.go]
usage:
  source: log
  seconds: 638
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 135
      output: 53735
      cache_read: 5851174
      cache_write: 215409
      cost: 3.4517
---
# T-0867 A write that grows an open story's claim into an in-progress story's says so and tells both stories as an overlapped change

## Work

This is I-0059's remediation 2. S-0242 and S-0198 passed the pull-time hold, then their tasks' touches widened each claim into the other's. Nothing said so until the close-out ran.

- After a write that changes the touches of an open story or of one of its tasks, compare the story's claim with every other in-progress story's claim. Use `workitem.Claim` and `PathsOverlap` from `flai/internal/workitem/hold.go`, as `cmd/accept_overlap.go` does. Compare only the paths the write added, so that an overlap the hold already allowed is not reported again. The writes are `flai task new`, `flai edit --touches`, `flai touches`, and the MCP `item_new` and `item_edit`. The dashboard's writes go through the CLI, so they are covered too.
- When the new paths overlap, print one line per other story naming it and the paths, and return them in `--json` and in the MCP result.
- Record a notice for each of the two stories in `.flai-cache/overlaps.jsonl` through `itemedit`. The `Overlap` record needs a cause other than an accepted story: add a field for the in-progress story whose claim was reached, and keep the reader accepting the old records. `inbox` and `wait_for_events` (`mcpserver/cursor.go`) then report it once to each story as an `overlapped` change. Its `cause` is the other story, its `to` is the paths, and its summary says the claims grew to overlap rather than that a story was accepted.
- The write itself always stands. The report is advisory, as `wip.overlap` is (ADR-0019, ADR-0046).

It waits for nothing. It shares no path with the long-poll task, so the two run together.

## Done when

- [ ] A Go test creates two in-progress stories with disjoint touches, gives a task of one a touch inside the other's claim through `flai task new`, and through the MCP `item_edit`, and sees the command name the other story and the paths
- [ ] The same test sees one `overlapped` change for each story in `inbox`, with `cause` the other story, reported once
- [ ] A write whose touches add no path inside another in-progress story's claim reports nothing
- [ ] Old `overlaps.jsonl` records written at acceptance still read and report as before
- [ ] `scripts/flai-test.sh` passes

## Notes
