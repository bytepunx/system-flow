---
id: T-0538
type: task
nature: feature
title: The activity page's agent panes show a live stream window
status: done
parent: S-0142
owner: alex
created: 2026-09-29T05:44:46Z
updated: 2026-09-29T06:00:24Z
transitions:
  - to: ready
    at: 2026-09-29T05:44:52Z
    by: agent-S-0142
  - to: in-progress
    at: 2026-09-29T05:51:27Z
    by: agent-S-0142
  - to: done
    at: 2026-09-29T06:00:24Z
    by: agent-S-0142
stream: S-0142
tags: []
touches: [flaiover/src, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
---
# T-0538 The activity page's agent panes show a live stream window

## Work

Add `/api/agent-stream/[story]` asking flai `agent.stream`, and a stream window in each activity card whose story has an agent run: the latest entries, scrolled to the end, polled while the agent runs and appended from the next offset. The activity page also asks `/api/host-agent` to know which stories have runs. Document it in `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`.

## Done when

- Component and route tests cover entries shown, appending from the offset, an agent that ended, and a flai too old to answer.
- `pnpm check`, lint, and the vitest suite pass.

## Notes
