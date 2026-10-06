---
id: T-0901
type: task
nature: feature
title: A fixture-board test runs the orchestrator's calls for each of its four permissions, on and off, through the guard and flai
status: done
parent: S-0219
owner: alex
created: 2026-10-05T04:47:33Z
updated: 2026-10-06T03:37:02Z
transitions:
  - to: ready
    at: 2026-10-06T03:16:31Z
    by: agent-S-0219
  - to: in-progress
    at: 2026-10-06T03:16:32Z
    by: agent-S-0219
  - to: done
    at: 2026-10-06T03:26:46Z
    by: agent-S-0219
stream: S-0219
tags: [flai]
touches: [flai/cmd/orchestrate_permissions_test.go]
after: [T-0893, T-0896]
usage:
  source: log
  seconds: 614
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 117
      output: 43660
      cache_read: 6623737
      cache_write: 152759
      cost: 3.0652
---
# T-0901 A fixture-board test runs the orchestrator's calls for each of its four permissions, on and off, through the guard and flai

## Work

The story's last criterion asks for tests of each permission on and off against a fixture board. The tasks before this one test their parts alone; this test runs them together, as the orchestrator would, so a gap between the guard and flai shows.

Build one fixture project in `flai/cmd/testdata/orchestrate` (or in code, with the helpers `cmd`'s tests already use): a backlog epic with no stories, an epic whose stories are all done, a complete draft and an incomplete one, two promotion candidates and a held backlog story, a ready column one short of its limit, and a ready story placed by hand an hour ago. With `FLAI_ROLE=orchestrate` and the manifest's `orchestration.permissions` set in turn, run each permission's calls through `flai guard` and then flai:

- each permission on: its call passes the guard and does what the criterion says, on this board
- each permission off: the guard refuses its call naming the permission, and the board is unchanged
- with `promote_to_ready` on, the first candidate is promoted, the second is refused at the limit, and the held story and the draft are never moved
- with `order_ready` on, the hand-placed story keeps its place

This task waits for T-0893 and T-0896, whose calls it runs. It runs with T-0899, whose paths it does not share.

## Done when

- the test covers the four permissions on and off, and the cases above, and changes nothing outside its temporary directory
- `go test ./cmd/` passes

## Notes
