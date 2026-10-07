---
id: T-1148
type: task
nature: remediation
title: The threads left open on archived items are resolved, flai check finds no threads.archived, and I-0073 is closed
status: backlog
parent: S-0277
owner: alex
created: 2026-10-07T01:20:56Z
updated: 2026-10-07T01:20:56Z
transitions: []
stream: S-0277
tags: [flai]
touches: [design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1147]
---
# T-1148 The threads left open on archived items are resolved, flai check finds no threads.archived, and I-0073 is closed

## Work

The fix stops new findings, but the threads already left open on archived stories stay until someone resolves them. On 2026-10-07 those were TH-0203 on S-0229 and TH-0232 on S-0272.

- Run `flai check` on main. Resolve each thread it still reports as `threads.archived` with `flai thread resolve TH-nnnn --reason "<S-nnnn> was accepted; S-0277"`. A thread anchored on a story belongs to the operator's conversation, so resolve only those the check names. If a thread's last entry is a question still waiting for the operator, ask on S-0277's plan thread first.
- Run `flai check --strict` and confirm it reports no `threads.archived`.
- Close the issue from the story's worktree: `flai issue close I-0073 --reason "..."`. The reason names the ADR from T-1147 and says that `flai accept` and `flai archive` now resolve the threads still open on what they archive. `design/issues/summary.md` regenerates with it.
- Tick the story's criteria with `flai criteria tick` once each is verified.

This task waits for T-1147: the issue is closed once the fix, its tests, and its ADR are all in place.

## Done when

- `flai check --strict` on main reports no `threads.archived`.
- I-0073 is closed with a reason that names the fix, and `design/issues/summary.md` no longer lists it.
- Both of S-0277's criteria are ticked.

## Notes
