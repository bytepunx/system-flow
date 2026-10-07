---
id: T-1146
type: task
nature: improvement
title: flai archive resolves the threads still open on the items it archives, and its dry run lists them
status: in-progress
parent: S-0277
owner: alex
created: 2026-10-07T01:20:27Z
updated: 2026-10-07T02:10:11Z
transitions:
  - to: ready
    at: 2026-10-07T02:10:10Z
    by: agent-S-0277
  - to: in-progress
    at: 2026-10-07T02:10:11Z
    by: agent-S-0277
stream: S-0277
tags: [flai, cli]
touches: [flai/cmd/archive.go, flai/cmd/archive_test.go]
after: [T-1142]
usage:
  source: log
  seconds: 624
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 8712
      cache_read: 886009
      cache_write: 53858
      cost: 0.6836
---
# T-1146 flai archive resolves the threads still open on the items it archives, and its dry run lists them

## Work

`flai archive` is the other way an item reaches `wip/archive`. It is used for a cancelled story, or for an item moved to done by an older flai. It leaves the item's threads open just as acceptance does.

- In `flai/cmd/archive.go`, after `repo.Archive(plan)`, call `threads.ResolveOnItems` with the IDs in `plan.Items`. The author is the one flai records for any write (`FLAI_AGENT`, else the config author). The reason is `<ID> was archived`.
- With `--dry-run`, list the threads `threads.OnItems` finds as ones it would resolve. Add `resolved_threads` to the `--json` output, and a `resolved TH-nnnn` line, or `would resolve`, to the text output.
- Update the command's `Long` help to say so.
- Test in a new `flai/cmd/archive_test.go`: archiving a cancelled story with an open thread on it resolves the thread with the reason. A dry run lists it and changes nothing. A resolved thread on the story is left as it is.

This task waits for T-1142, whose functions it calls. It touches no file of T-1144, so the two run together.

## Done when

- `flai archive` leaves no thread on an item it archived `open` or `answered`, and `--dry-run` and `--json` name the threads.
- The new test passes, and `flai test flai/cmd` passes.

## Notes
