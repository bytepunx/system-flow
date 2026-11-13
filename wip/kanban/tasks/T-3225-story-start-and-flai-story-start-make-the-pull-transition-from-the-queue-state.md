---
id: T-3225
type: task
nature: feature
title: story_start and flai story start make the pull transition from the queue state the workflow names
status: backlog
parent: S-0746
owner: alex
created: 2026-11-12T16:07:45Z
updated: 2026-11-12T16:07:45Z
transitions: []
stream: S-0746
tags: [cli]
touches: [flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/story_start_test.go, flai/internal/storystart/start.go]
---
# T-3225 story_start and flai story start make the pull transition from the queue state the workflow names

## Work

Replace the references to a state's name in `story_start.go` and `storystart/start.go`: a story starts by the transition with the pull effect from the state it leaves, and one in any other state is refused with today's words and that state's name (ADR-0153). Test a start on a workflow with a state renamed.

## Done when

Both files name no state, and `story_start_test.go` passes on both workflows.

## Notes
