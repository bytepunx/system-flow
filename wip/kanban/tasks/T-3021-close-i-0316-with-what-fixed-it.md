---
id: T-3021
type: task
nature: remediation
title: Close I-0316 with what fixed it
status: backlog
parent: S-0679
owner: alex
created: 2026-11-10T10:32:08Z
updated: 2026-11-10T10:32:08Z
transitions: []
stream: S-0679
tags: [issues]
touches: [design/issues/I-0316-flai-check-finds-item-heading-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-3014, T-3017]
---
# T-3021 Close I-0316 with what fixed it

## Work

Close I-0316 with `flai issue close I-0316 --reason "<reason>"`, the reason naming what fixed it: the close-out's record step confirms a finding outside the story on a work item under the checked-write lock and records only what still holds, the close-out's container mounts the lock file, and the test that reproduces the cause. Name T-3014, T-3017, and the ADR T-3018 writes if it exists by then. `flai issue close` writes `design/issues/summary.md` again.

This task waits for T-3014 and T-3017: the reason names what they built. It runs beside T-3018, with no path in common.

## Done when

- I-0316 is closed with a reason naming the fix and its test.
- `design/issues/summary.md` no longer lists I-0316.

## Notes
