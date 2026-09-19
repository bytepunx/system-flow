---
id: T-0210
type: task
nature: feature
title: "flai order: place a story in the pull order, with the ordering rules in workitem and the same sequence in flai board and the MCP tools"
status: done
parent: S-0057
owner: alex
created: 2026-09-19T06:41:25Z
updated: 2026-09-19T06:47:45Z
transitions:
  - to: ready
    at: 2026-09-19T06:42:03Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T06:42:04Z
    by: system-flow
  - to: done
    at: 2026-09-19T06:47:45Z
    by: system-flow
stream: S-0057
tags: []
---

# T-0210 flai order: place a story in the pull order, with the ordering rules in workitem and the same sequence in flai board and the MCP tools

## Work
In `flai/internal/workitem`: the pull sequence of a column (stories named in `order` first, in that order, then the rest by ID), and `Board.Place`, which puts a ready or backlog story before or after another story of the same state, or at the top or bottom of its state. It refuses anything else with a message that says why: an epic or task, a story in another state, a reference in a different column, itself. On write the list is normalised: ready stories first, then backlog stories, IDs that are no longer ready or backlog stories dropped, and backlog stories that were never placed and still sort by ID at the end left unlisted so the file does not grow to name the whole backlog. `Move` to ready inserts the story at the end of the ready section rather than at the end of the list, so a placed backlog story that becomes ready does not sit behind backlog entries. `flai order <id> --before|--after <other>|--top|--bottom` with `--json`; `flai board`, and the MCP `board` and `inbox` tools, list backlog and ready in the pull sequence. Behaviour tests for the rules and the command.

## Done when
- `flai order` changes `board.md` and refuses what the rules refuse, with tests for each rule
- `flai board` prints backlog and ready in the pull sequence
- `make flai-test` passes

## Notes
