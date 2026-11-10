---
id: T-2986
type: task
nature: remediation
title: Close I-0336 with what fixed it
status: backlog
parent: S-0710
owner: alex
created: 2026-11-10T10:29:38Z
updated: 2026-11-10T10:29:38Z
transitions: []
stream: S-0710
tags: [issues]
touches: [design/issues/I-0336-flai-check-finds-adr-index-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-2980]
---
# T-2986 Close I-0336 with what fixed it

## Work

Close I-0336 from the story's worktree with `flai issue close I-0336 --reason "<reason>"`, the reason naming the cause and T-2980's fix: an autocommitted ADR write committed an index rendered from the working tree's ADR files, so an ADR file left uncommitted put its row in main's index, and the commit now records the index HEAD's ADR files render. `flai issue close` writes the issue and `design/issues/summary.md`.

This task waits for T-2980, since the reason names the fix T-2980 makes and the issue closes only once it is in. It runs beside T-2983, which changes no path of this task.

## Done when

- I-0336 is closed with a reason naming the cause and the fix.
- `design/issues/summary.md` no longer lists I-0336, and `flai check --strict` reports no `issues.summary`.

## Notes
