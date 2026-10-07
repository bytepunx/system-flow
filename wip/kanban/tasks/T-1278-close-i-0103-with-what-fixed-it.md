---
id: T-1278
type: task
nature: improvement
title: Close I-0103 with what fixed it
status: backlog
parent: S-0309
owner: alex
created: 2026-10-07T23:27:40Z
updated: 2026-10-07T23:27:40Z
transitions: []
stream: S-0309
tags: [issues]
touches: [design/issues/I-0103-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md, design/issues/summary.md]
after: [T-1275, T-1276, T-1277]
---
# T-1278 Close I-0103 with what fixed it

## Work

Run `flai issue close I-0103 --reason "<what fixed it>" --commit` in the story's worktree, naming the ADR, the bounded wait, the open thread, and the retry. It waits for T-1275, T-1276, and T-1277 because the issue closes only once the fix and its words are in place.

## Done when

- I-0103's status is closed, its body ends with the reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
