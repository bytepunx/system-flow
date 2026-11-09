---
id: T-2808
type: task
nature: remediation
title: Close I-0317 naming the event-driven cancel that fixed it
status: backlog
parent: S-0638
owner: alex
created: 2026-11-09T22:08:12Z
updated: 2026-11-09T22:08:12Z
transitions: []
stream: S-0638
tags: [issues]
touches: [design/issues/I-0317-testcancelstopsthewaittoaskagain-times-out-against-its-local-release-server.md, design/issues/summary.md]
after: [T-2806]
---
# T-2808 Close I-0317 naming the event-driven cancel that fixed it

## Work

Run `flai issue close I-0317 --reason "<text>"` in the story's worktree. The reason names the cause, a 100 ms wall-clock deadline that raced the first request on a loaded host, and the fix: the test cancels from a seam called when the pause between attempts begins, and a second test holds the first answer past the old deadline. `flai issue close` writes the issue file and regenerates `design/issues/summary.md`.

Waits for T-2806: the reason names what T-2806 changed, so the fix must be in and passing first.

## Done when

- I-0317 is `closed`, with the reason in its file.
- `design/issues/summary.md` no longer lists I-0317 as open.
- `flai check --strict` is clean.

## Notes
