---
id: T-2835
type: task
nature: remediation
title: Close I-0324 with what fixed it
status: backlog
parent: S-0649
owner: alex
created: 2026-11-10T06:58:39Z
updated: 2026-11-10T06:58:39Z
transitions: []
stream: S-0649
tags: [issues]
touches: [design/issues/I-0324-testeachrequestistimedbeneaththetransport-reads-the-request-answered-event.md, design/issues/summary.md]
after: [T-2834]
---
# T-2835 Close I-0324 with what fixed it

## Work

Second layer: it waits for T-2834, because the reason names the fix T-2834 makes and the runs that show it holds.

Run `flai issue close I-0324 --reason "<reason>"` in the story's worktree. The reason says the test read the log once, before the channel client logged the event after the answer's write, that the tests now poll the log for it with a deadline, and names the reproducing run and the passing runs from T-2834's Notes.

## Done when

- I-0324's status is closed, with the reason.
- `design/issues/summary.md` no longer lists I-0324.
- `flai check --strict` is clean.

## Notes
