---
id: T-1052
type: task
nature: feature
title: The close-out's issue for findings outside the story quotes a finding's code spans so the markdown lint passes
status: done
parent: S-0227
owner: alex
created: 2026-10-06T22:06:48Z
updated: 2026-10-06T22:07:20Z
transitions:
  - to: ready
    at: 2026-10-06T22:07:19Z
    by: agent-S-0227
  - to: in-progress
    at: 2026-10-06T22:07:20Z
    by: agent-S-0227
  - to: done
    at: 2026-10-06T22:07:20Z
    by: agent-S-0227
stream: S-0227
tags: []
touches: [flai/cmd/check.go, flai/cmd/check_test.go, design/issues]
usage:
  source: log
  seconds: 0
  models: []
---

# T-1052 The close-out's issue for findings outside the story quotes a finding's code spans so the markdown lint passes

## Work

S-0227's first close-out stopped at the markdown lint on an issue flai wrote during it. `recordOutside` (`flai/cmd/check.go`) copied a markdown finding's context, a code span with a trailing space, into the issue that records findings outside the story, and that made an MD038 finding of its own. Write each finding's message with its backticks as quotes (`findingText`). It waits for no task: it is a fix found at close-out.

## Done when

- A test pins a finding that quotes a code span written without backticks.
- The close-out passes the markdown lint over the issues it records.

## Notes

- Moved to done before its body was written, so the body was edited by hand: `flai edit` refuses a closed item.
