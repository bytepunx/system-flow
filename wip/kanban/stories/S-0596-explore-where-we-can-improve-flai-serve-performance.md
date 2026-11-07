---
id: S-0596
type: story
nature: research
title: Explore where we can improve flai serve performance
status: in-progress
owner: alex
created: 2026-11-07T22:06:13Z
updated: 2026-11-07T22:13:32Z
transitions:
  - to: ready
    at: 2026-11-07T22:06:33Z
    by: alex
  - to: in-progress
    at: 2026-11-07T22:06:54Z
    by: agent-S-0596
tags: [cli]
topics: [server-side, back-end]
touches: [design/system/server-performance.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 1238
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 239
          output: 33832
          cache_read: 21955989
          cache_write: 452280
          cost: 7.2937
    - kind: orchestrator
      seconds: 59
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 164
          cache_read: 561858
          cache_write: 116030
          cost: 0.1718
forecast:
  duration: 1h
  delivery: 2026-11-07T23:10:00Z
  basis: "S-0152 measured the board's and MCP's reads alone in 1101 s of agent time; this story covers four groups of operations, the charts over several windows and the close-out and acceptance among them, and a review of serial I/O, about three times that."
  by: planner-S-0596
  at: 2026-11-07T22:10:58Z
---
# S-0596 Explore where we can improve flai serve performance

## Goal

flai is snappy when it's first being used on a project, but nearly 600 stories in and the way we're scanning and operating on files in the project directly appear to be taking a real toll on its responsiveness.

Noted slow operations:
- board load times
- creation and editing of epics/stories
- all chart load times (especially when exceeding 7 days of history)
- steps around close-out and review acceptance

## Acceptance criteria
- [ ] Conduct and document a thorough exploration of the file operations
- [ ] Find areas that may be slow due to serialization of I/O where it could be conducted in parallel
- [ ] Write up a performance improvement epic with stories to address any actionable findings (no epic or stories if anything would require a real change to how we track state in flai)

## Tasks
- T-2607 Explore and measure the dashboard's reads: agent.status, board.get, settings.get, threads, messages, and item.get, and whether flai serve answers channel requests one at a time
- T-2608 Explore and measure the writes: item.new, item.edit, item.finalize, item.move, issue.story, accept.run, and the close-out job
- T-2609 Explore and measure the chart data: stats.get over windows of 7, 30, 90 days and all history
- T-2610 Record the exploration's findings in server-performance.md: the cost of each slow read and write at 596 stories, the serial I/O that could run in parallel, and what would need a change to how flai keeps state
- T-2611 Draft a performance epic with one draft story per actionable cause, each with the cause's numbers and a proposed fix, and link them from server-performance.md

## Notes

### Planning

Planned by planner-S-0596 after the story's agent had started it and written its tasks. No earlier plan and no thread on the story existed, so this is its first plan.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `design/system/server-performance.md` | design: S-0152's finding on server-side response times, with the tools it measures by (request timing, `flai hostapi --timing`, `FLAI_PPROF_ADDR`) and the Stories table E-0012 hung from. T-2610 and T-2611 write this story's finding there |

The operator declared `flai/cmd`. The story's agent narrowed it to the file its tasks name, which fits: no task changes a command. `flai touches suggest` lists the paths that change with `flai/cmd`: `docs/users/flai-reference.md` 32%, `docs/users/flai.md` 20%, `design/system/flai-cli.md` 19%, `docs/operators/settings.md` 15%. None is taken: a research story changes no command and documents no new behaviour. The epic and stories T-2611 drafts are written through flai on main, so they are no touch.

No `after`: the story builds on S-0152, S-0156 to S-0162, S-0471, S-0500, and S-0531, all done and on main.

Topics `server-side` and `back-end` are added: they are `server-performance.md`'s, which T-2610 and T-2611 write.

Tasks: the story's agent wrote all five; none is added. The three explorations, T-2607, T-2608, and T-2609, write nothing in the repository and take no touches, so they run at once. T-2610 waits for all three and is the one task that writes the section; T-2611 waits for it, since it drafts a story per cause the section names and fills in its table. Together they cover the four slow operations the goal names and the three criteria.

Figures:

- Forecast 1h, adjusted from flai's 36m (531 s per unit of size over 3 done research stories, times size 4). S-0152 measured the board's and MCP's reads alone in 1101 s of agent time; this story covers four groups of operations, the charts over several windows and the close-out and acceptance among them, which run git and the test tiers, and a review of serial I/O, about three times that.
- Delivery 2026-11-07T23:10:00Z: the forecast from its start at 22:07.
- Cost of delay: none yet. The story has no inputs and no epic, so `flai cod` refuses. The inputs are the operator's, asked on the plan's thread.
