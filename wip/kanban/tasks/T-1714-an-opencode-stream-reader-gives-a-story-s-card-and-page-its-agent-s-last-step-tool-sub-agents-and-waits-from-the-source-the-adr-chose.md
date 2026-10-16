---
id: T-1714
type: task
nature: feature
title: An OpenCode stream reader gives a story's card and page its agent's last step, tool, sub-agents, and waits from the source the ADR chose
status: backlog
parent: S-0411
owner: alex
created: 2026-10-16T05:19:48Z
updated: 2026-10-16T05:19:48Z
transitions: []
stream: S-0411
tags: [cli]
touches: [flai/internal/serve/stream_opencode.go, flai/internal/serve/stream_opencode_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-1712, T-1713]
---
# T-1714 An OpenCode stream reader gives a story's card and page its agent's last step, tool, sub-agents, and waits from the source the ADR chose

## Work

- In `flai/internal/serve/stream_opencode.go`, read the source T-1712's ADR chose into the entries T-1713 keeps: the last step's text, the tool in use, a sub-agent's start and end, and a `wait_for_events` in progress.
- Register it for `opencode`, and check that `StoryActivity` and `Activity` in `agents.go` give a story on `opencode` its card's status.
- Test on S-0410's recorded output, and on a recorded event stream if the ADR chose `opencode serve`.
- Waits for T-1712 (the source) and T-1713 (the per-harness reader it registers in).

## Done when

A story on `opencode` has its card's and page's status from the reader, in the entry shape the dashboard reads, and the tests pass.

## Notes
