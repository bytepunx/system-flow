---
id: T-0368
type: task
nature: feature
title: "flai serve starts an agent for every story that enters ready, within the WIP limit, and knows each one's state: working, waiting, failed"
status: done
parent: S-0104
owner: alex
created: 2026-09-23T17:25:01Z
updated: 2026-09-23T17:33:41Z
transitions:
  - to: ready
    at: 2026-09-23T17:29:25Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T17:29:26Z
    by: system-flow
  - to: done
    at: 2026-09-23T17:33:41Z
    by: system-flow
stream: S-0104
tags: []
---

# T-0368 flai serve starts an agent for every story that enters ready, within the WIP limit, and knows each one's state: working, waiting, failed

## Work
The launcher starts an agent for each ready story that has none, in pull order, while stories in progress plus agents started for stories still in ready stay under the in-progress limit. A story created in ready counts as entering it. Each run gets its own `FLAI_AGENT` (`<name>-<story>`), so wait_for_work and cursors are its own. Its signs (its cursor, its story's narrative) do not count as someone attending. State per story: `running`, and `ended` with an exit code. An agent that ends with its story neither in review nor done has failed.

## Done when
- Tests: two stories entering ready start two agents under a limit of 2 and one under a limit of 1.
- The next story starts when an agent ends.
- An exit that leaves the story in progress is a failure.
- A restart of flai serve starts nothing.

## Notes
