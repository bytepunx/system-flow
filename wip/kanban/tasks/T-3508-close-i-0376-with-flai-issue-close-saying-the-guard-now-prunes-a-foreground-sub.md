---
id: T-3508
type: task
nature: remediation
title: Close I-0376 with flai issue close, saying the guard now prunes a foreground sub-agent whose result is in the session transcript
status: backlog
parent: S-0819
owner: alex
created: 2026-11-17T03:35:06Z
updated: 2026-11-17T03:35:06Z
transitions: []
stream: S-0819
tags: [guard, issues]
touches: [design/issues/I-0376-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md, design/issues/summary.md]
after: [T-3506, T-3507]
---
# T-3508 Close I-0376 with flai issue close, saying the guard now prunes a foreground sub-agent whose result is in the session transcript

## Work

Run `flai issue close I-0376 --reason "<reason>"`, the reason naming the new ADR, the guard's prune at `wait_for_events`, and the test T-3506 added. It waits for T-3506 and T-3507 because the reason names the fix and the documents that describe it, which must be in place first.

## Done when

- I-0376 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- Criterion 2 of S-0819 is ticked.

## Notes
