---
id: T-2332
type: task
nature: improvement
title: Say in the dashboard's design and user guide that a story's Agent pane shows its agent's stream
status: done
parent: S-0541
owner: alex
created: 2026-10-29T14:15:26Z
updated: 2026-10-29T14:17:23Z
transitions:
  - to: ready
    at: 2026-10-29T14:17:14Z
    by: agent-S-0541
  - to: in-progress
    at: 2026-10-29T14:17:15Z
    by: agent-S-0541
  - to: done
    at: 2026-10-29T14:17:23Z
    by: agent-S-0541
stream: S-0541
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-2331]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-sonnet-5-5
      input: 3
      output: 376
      cache_read: 123309
      cache_write: 5435
      cost: 0.0378
---
# T-2332 Say in the dashboard's design and user guide that a story's Agent pane shows its agent's stream

## Work

Describe what T-2331 built, in the same story, as `documentation.md` asks when user-facing behaviour changes.

- `design/system/flaiover-dashboard.md`: in the Views table's `/items/:id` row, say the Agent pane shows the story agent's stream (`AgentStream.svelte`, S-0142) below its status once the agent has started, open while it works or waits, as the Activity page shows it (S-0541). In the bullet "An agent's stream (S-0142)" under the data section, say the story page's Agent pane follows the same stream.
- `docs/users/flaiover.md`: in the story page's sections table, the Agent row says the pane also shows what the agent is saying and doing, as on the Activity page.

It waits for T-2331 so that the words describe what was built, such as whether the stream opens on its own.

## Done when

- `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` say the Agent pane shows the story agent's stream, and when it is open.
- `updated` is bumped on both documents.
- `flai test` on the two files passes the markdown lint.

## Notes

Drafted by the planner for S-0541.
