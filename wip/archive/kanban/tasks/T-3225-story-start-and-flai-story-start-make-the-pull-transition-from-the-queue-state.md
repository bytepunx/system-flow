---
id: T-3225
type: task
nature: feature
title: story_start and flai story start make the pull transition from the queue state the workflow names
status: done
parent: S-0746
owner: alex
created: 2026-11-12T16:07:45Z
updated: 2026-11-15T00:40:17Z
transitions:
  - to: ready
    at: 2026-11-14T23:56:37Z
    by: agent-S-0746
  - to: in-progress
    at: 2026-11-15T00:25:24Z
    by: agent-S-0746
  - to: done
    at: 2026-11-15T00:40:17Z
    by: agent-S-0746
stream: S-0746
tags: [cli]
touches: [flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/story_start_test.go, flai/internal/storystart/start.go, flai/internal/storystart/start_test.go]
after: [T-3224]
usage:
  source: log
  seconds: 893
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 5522
      cache_read: 2004547
      cache_write: 34469
      cost: 0.7872
---
# T-3225 story_start and flai story start make the pull transition from the queue state the workflow names

## Work

Replace the references to a state's name in `story_start.go` and `storystart/start.go`: a story starts by the transition with the pull effect from the state it leaves, and one in any other state is refused with today's words and that state's name (ADR-0153). Test a start on a workflow with a state renamed.

## Done when

Both files name no state, and `story_start_test.go` passes on both workflows.

## Notes
