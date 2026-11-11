---
id: T-3122
type: task
nature: remediation
title: Close I-0345 with what fixed it
status: backlog
parent: S-0717
owner: alex
created: 2026-11-11T09:07:34Z
updated: 2026-11-11T09:07:34Z
transitions: []
stream: S-0717
tags: [flai]
touches: [design/issues/I-0345-wait-for-events-fails-with-no-front-matter-while-another-story-s-checked.md, design/issues/summary.md]
after: [T-3120, T-3121]
---
# T-3122 Close I-0345 with what fixed it

## Work

Close I-0345 with `flai issue close I-0345 --reason`, which writes the issue file and `design/issues/summary.md`. The reason names three things:

- S-0716's loader change, under which `Repo.List` skips a file that does not parse.
- T-3121's tests, which pin `wait_for_events` past a `PLACEHOLDER` task file in one project and across a folder of projects.
- What T-3120 found wrote the file, and the task or thread that carries a fix when flai wrote it.

Name I-0344 as the issue that records the same event.

This task waits for T-3120, which writes the same issue file and finds what the reason names, and for T-3121, whose tests the reason cites.

## Done when

- I-0345 is closed, and its reason names the fix, the tests, and the writer.
- `design/issues/summary.md` no longer lists I-0345 as open.
- `flai check --strict` is clean.

## Notes
