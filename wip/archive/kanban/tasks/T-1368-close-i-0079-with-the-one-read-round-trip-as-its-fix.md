---
id: T-1368
type: task
nature: remediation
title: Close I-0079 with the one-read round trip as its fix
status: done
parent: S-0290
owner: alex
created: 2026-10-08T08:42:22Z
updated: 2026-10-08T09:02:47Z
transitions:
  - to: ready
    at: 2026-10-08T09:02:40Z
    by: agent-S-0290
  - to: in-progress
    at: 2026-10-08T09:02:40Z
    by: agent-S-0290
  - to: done
    at: 2026-10-08T09:02:47Z
    by: agent-S-0290
stream: S-0290
tags: [issues]
touches: [design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/summary.md]
after: [T-1367]
usage:
  source: log
  seconds: 7
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 2
      output: 384
      cache_read: 150723
      cache_write: 4731
      cost: 0.0757
---
# T-1368 Close I-0079 with the one-read round trip as its fix

## Work

Run `flai issue close I-0079 --reason "<reason>"` in the story worktree. The reason says what fixed it. `TestRoundTripRepositoryItems` now parses each item from one read of its file and compares the marshal with that read. It skips a file gone since the list. A test rewrites and removes items between the list and the check. Name the commit or T-1367.

Closing the issue rewrites `design/issues/summary.md`. Both paths are in the manifest's `claims.shared` (`design/issues`), and S-0321 also touches the I-0079 file. Sync before closing so the close lands on the issue's latest instances.

It waits for T-1367: the reason names T-1367's fix, and the issue closes only once that fix passes.

## Done when

- I-0079 is closed, and its reason names the one-read round trip and the reproduction test.
- `design/issues/summary.md` no longer lists I-0079.
- `flai check --strict` is clean on the two files.

## Notes

Drafted by planner-S-0290.
