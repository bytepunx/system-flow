---
id: T-1573
type: task
nature: improvement
title: Close I-0097 naming the fix
status: ready
parent: S-0304
owner: alex
created: 2026-10-14T20:56:21Z
updated: 2026-10-14T20:56:27Z
transitions:
  - to: ready
    at: 2026-10-14T20:56:27Z
    by: agent-S-0304
stream: S-0304
tags: []
touches: [design/issues/I-0097-a-story-s-agent-that-runs-flai-stream-open-is-left-with-the-story-in-ready-so-permission-prompt-refuses-its-claude-write-until-it-moves-the-story-itself.md, design/issues/summary.md]
after: [T-1572]
---
# T-1573 Close I-0097 naming the fix

## Work

Close I-0097 with `flai issue close I-0097 --reason` (MCP `issue_close` with commit), saying that `flai stream open` now starts a ready story, so the agent is never left with it in ready.

Waits for the fix task, since the reason names what it did.

## Done when

- I-0097 is closed with the reason, and `design/issues/summary.md` no longer lists it.

## Notes
