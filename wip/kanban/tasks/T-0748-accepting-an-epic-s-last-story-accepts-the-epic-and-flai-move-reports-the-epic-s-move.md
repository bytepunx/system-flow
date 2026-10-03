---
id: T-0748
type: task
nature: improvement
title: Accepting an epic's last story accepts the epic, and flai move reports the epic's move
status: done
parent: S-0200
owner: alex
created: 2026-10-03T07:11:21Z
updated: 2026-10-03T07:32:00Z
transitions:
  - to: ready
    at: 2026-10-03T07:11:52Z
    by: agent-S-0200
  - to: in-progress
    at: 2026-10-03T07:19:14Z
    by: agent-S-0200
  - to: done
    at: 2026-10-03T07:32:00Z
    by: agent-S-0200
stream: S-0200
tags: []
touches: [flai/cmd/move.go, flai/cmd/move_cancel.go, flai/cmd/accept.go, flai/internal/preview, flai/cmd/release_pending_test.go, flai/cmd/hostapi_reads_test.go, flai/cmd/accept_epic_test.go]
after: [T-0745]
usage:
  source: log
  seconds: 766
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 70
      output: 24675
      cache_read: 4270683
      cache_write: 95696
      cost: 1.9327
---
# T-0748 Accepting an epic's last story accepts the epic, and flai move reports the epic's move

## Work

In `flai accept`, and `flai move <story> done`, which the dashboard runs, when the story accepted is its epic's last open one, walk the epic to done, roll its usage up, and archive it with the story in the one acceptance commit. `preview.Accept` reports it in a dry run. `flai move` prints the epic's move and returns it in `--json`.

It waits for T-0745, whose `Follow` it calls.

## Done when

- [x] Accepting the last open story moves its epic to done and archives it; accepting another leaves the epic as it is
- [x] `flai move` text and `--json` output name the epic's move
- [x] Tests in `flai/cmd` cover acceptance of the last story and the move output

## Notes
