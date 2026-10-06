---
id: T-1090
type: task
nature: improvement
title: The story agent's start prompt says to write the narrative's state and end on an open question, and to end when wait_for_events answers end
status: backlog
parent: S-0272
owner: alex
created: 2026-10-06T22:52:57Z
updated: 2026-10-06T22:52:57Z
transitions: []
stream: S-0272
tags: [cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
---
# T-1090 The story agent's start prompt says to write the narrative's state and end on an open question, and to end when wait_for_events answers end

## Work

Criterion 1, the prompt. `flai/internal/harness/harness.go` tells the story's agent to hold `wait_for_events` until the thread is answered, in two places:

- the rules for a story run (about line 503)
- the commit run (about line 238)

Rewrite both. After `thread_open`, the agent writes the narrative's `## Current state` and `## Next steps` and ends. flai serve starts it again in its session when the thread is answered, with the answer in its first `inbox`. If it has a task in progress it goes on with that first. When `wait_for_events` answers `end: true`, it ends.

Leave alone:

- The planner's and the analyzer's prompts (about lines 311 and 416). flai serve does not restart them on an answer, so they keep holding.
- The rule against waiting for a sub-agent with `wait_for_events`.

Update the prompt assertions in `harness_test.go` to match.

It waits for no task. The field's name, `end`, is the story's own, so the prompt need not wait for the MCP task.

## Done when

- Neither story prompt tells the agent to hold `wait_for_events` for an answer. Both say to write the narrative's state and end, and to end on `end: true`.
- `harness_test.go` asserts the new wording, and passes with `scripts/flai-test.sh`.

## Notes

Drafted by the planner.
