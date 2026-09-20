---
id: T-0268
type: task
nature: feature
title: "flai: narratives with activity, and the designer's inbox, as methods"
status: done
parent: S-0074
owner: alex
created: 2026-09-20T08:29:59Z
updated: 2026-09-20T08:35:09Z
transitions:
  - to: ready
    at: 2026-09-20T08:32:28Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T08:32:28Z
    by: system-flow
  - to: done
    at: 2026-09-20T08:35:09Z
    by: system-flow
stream: S-0074
tags: []
---
# T-0268 flai: narratives with activity, and the designer's inbox, as methods

## Work
activity.get (each narrative's stream, agent, session, updated, the story's state, blocked, the task in progress, the last log entry) and inbox.designer (threads awaiting the designer, open questions outside the generated block, stories in review, blocked items, overlapping touches from the check's own rule), with the keys the dashboard already uses so that nothing becomes new by the change. Named apart from the agent's MCP inbox.

## Done when
- Tests for each kind of entry and for key stability against the TypeScript hash

## Notes
