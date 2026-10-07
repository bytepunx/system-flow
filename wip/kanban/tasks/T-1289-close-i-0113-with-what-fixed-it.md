---
id: T-1289
type: task
nature: improvement
title: Close I-0113 with what fixed it
status: backlog
parent: S-0313
owner: alex
created: 2026-10-07T23:39:40Z
updated: 2026-10-07T23:39:40Z
transitions: []
stream: S-0313
tags: [issues]
touches: [design/issues/I-0113-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md, design/issues/summary.md]
after: [T-1286, T-1287]
---
# T-1289 Close I-0113 with what fixed it

## Work

Run `flai issue close I-0113 --reason "<what fixed it>"` in the story's worktree, so the close lands on the story branch. The reason names both fixes: the integration tier reads `go test -json` (T-1287), and a plain finding keeps the failure lines go test printed (T-1286).

It waits for T-1286 and T-1287: the reason states what they did, so both must be done.

## Done when

- I-0113 is closed with that reason.
- `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean.

## Notes

Drafted by the planner for S-0313.
