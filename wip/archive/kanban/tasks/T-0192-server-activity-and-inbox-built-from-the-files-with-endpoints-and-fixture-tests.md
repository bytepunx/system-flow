---
id: T-0192
type: task
nature: feature
title: "Server: activity and inbox built from the files, with endpoints and fixture tests"
status: done
parent: S-0042
owner: alex
created: 2026-09-19T04:46:17Z
updated: 2026-09-19T04:48:45Z
transitions:
  - to: ready
    at: 2026-09-19T04:46:41Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:46:41Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:48:45Z
    by: system-flow
stream: S-0042
tags: []
touches: [flaiover/src/lib/server, flaiover/src/routes/api]
---

# T-0192 Server: activity and inbox built from the files, with endpoints and fixture tests

## Work
`src/lib/server/activity.ts`: for each narrative under `wip/agents` except `index.md`, the stream, title, agent, session, `updated` and its age, the last log entry (its time and text), the story's status and blocked flag, and the task in progress under it, if any; newest first; an agent that stopped logging simply shows an old age. `src/lib/server/inbox.ts`: entries of five kinds, each with a stable key, a title, a link target, and a time: threads awaiting the designer (open threads whose last entry is not the designer's); open questions from narratives (the bullets under `## Open questions` outside the generated threads block); stories in review; blocked items with the reason; and `wip.overlap` warnings, taken from `flai check --json` so the rule has one implementation, cached until the repository changes, and left out with a note when flai is not available. `GET /api/activity` and `GET /api/inbox` (entries, counts by kind, total). Tests on a temp copy of the fixture extended with a question, a blocked item, a story in review, and overlapping touches.

## Done when
- The fixture tests pass
- flaiover lint and svelte-check are clean

## Notes
