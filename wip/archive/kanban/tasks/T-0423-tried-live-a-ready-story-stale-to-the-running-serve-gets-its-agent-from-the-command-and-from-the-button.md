---
id: T-0423
type: task
nature: feature
title: "Tried live: a ready story stale to the running serve gets its agent from the command and from the button"
status: done
parent: S-0115
owner: alex
created: 2026-09-24T09:07:10Z
updated: 2026-09-24T09:32:07Z
transitions:
  - to: ready
    at: 2026-09-24T09:20:30Z
    by: agent-S-0115
  - to: in-progress
    at: 2026-09-24T09:20:31Z
    by: agent-S-0115
  - to: done
    at: 2026-09-24T09:32:07Z
    by: agent-S-0115
stream: S-0115
tags: []
touches: [flai/cmd]
---
# T-0423 Tried live: a ready story stale to the running serve gets its agent from the command and from the button

## Work

Against a scratch project served by a scratch flai serve built from the story branch, with a harmless host command for the agent, never the operator's serve: put a story in ready before serve starts so the launcher does not start it, then start it with `flai serve agent start`, and another from the dashboard's button. Watch serve track each run to its outcome.

## Done when

- Both starts are observed, with the run in `serve/agents.json`, the journal entry, and the outcome once each ends, and the observations are in the story notes.

## Notes
