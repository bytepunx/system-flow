---
id: T-2005
type: task
nature: remediation
title: The design records that every work item write takes the items lock, and I-0175 is closed with what fixed it
status: backlog
parent: S-0479
owner: alex
created: 2026-10-21T12:11:28Z
updated: 2026-10-21T12:11:28Z
transitions: []
stream: S-0479
tags: [design, issues]
touches: [design/system/work-hierarchy.md, design/adrs, design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md, design/issues/summary.md]
after: [T-2001, T-2002]
---
# T-2005 The design records that every work item write takes the items lock, and I-0175 is closed with what fixed it

## Work

Record the rule the code now keeps, and close the issue. It waits for T-2001 and T-2002 because it describes, and closes the issue on, what they built.

- In `design/system/work-hierarchy.md`, under the paragraph on who changes what after an item is made, say that every flai write to a work item file holds the project's items lock under `.flai-cache`. Say that it reads the item again under that lock and changes only the fields it owns, so that `flai serve`, the MCP server, and a shell never lose each other's writes. Name the lock file and its wait.
- Write an ADR with `flai adr new` if the lock is a decision a reader would look for there: one project-wide lock rather than one per item, held across a transition's several files and not across `check.Run`. Add its row to `design/adrs/README.md` as `flai adr new` does. `design/adrs` is a folder touch because the ADR's number is not known yet.
- Close I-0175 with `flai issue close I-0175 --reason`, naming the story, the lock, and the tests that reproduce it. This updates the issue file and `design/issues/summary.md`.
- Tick each acceptance criterion through flai once verified.

## Done when

- `design/system/work-hierarchy.md` states the rule, and the ADR, if written, is accepted and indexed.
- I-0175 is closed with a reason that says what fixed it.
- `flai check --strict` and the markdown lint pass.

## Notes

Drafted by the planner. If the operator folds S-0480 into this story on the plan's thread, close I-0176 here too.
