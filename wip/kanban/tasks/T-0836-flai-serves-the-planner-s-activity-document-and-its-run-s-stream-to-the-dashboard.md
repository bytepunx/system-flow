---
id: T-0836
type: task
nature: feature
title: flai serves the planner's activity document and its run's stream to the dashboard
status: done
parent: S-0259
owner: alex
created: 2026-10-05T00:07:40Z
updated: 2026-10-05T00:13:54Z
transitions:
  - to: ready
    at: 2026-10-05T00:08:19Z
    by: agent-S-0259
  - to: in-progress
    at: 2026-10-05T00:08:19Z
    by: agent-S-0259
  - to: done
    at: 2026-10-05T00:13:54Z
    by: agent-S-0259
stream: S-0259
tags: []
touches: [flai/internal/hostapi, flai/internal/serve/stream.go, flai/internal/serve/stream_test.go, flai/cmd/serve_actions.go, flaiover/src/lib/server/agent.ts]
usage:
  source: log
  seconds: 335
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 62
      output: 22911
      cache_read: 3203930
      cache_write: 94450
      cost: 1.6967
---
# T-0836 flai serves the planner's activity document and its run's stream to the dashboard

## Work

Two reads on the channel flai serve offers a dashboard (`flai/internal/hostapi`), so that the Planner page can show what the planner did and what its current run is doing:

- `activity.document {kind}`: a strategic agent's activity document, `wip/agents/<kind>.md`, as `workitem.Repo.Activity` reads it: its front matter totals and its entries, oldest first, with `path` repository-relative. A document that does not exist yet answers empty with zero totals; a kind that is not a strategic agent is a bad request naming the kinds.
- `agent.stream {plan}`: given `plan`, an epic's or a story's ID, in place of `story`, the stream of the newest planner run flai serve started for that item (`AgentState.Plans`), read as a story's agent's is; `item` names it in the answer. Exactly one of `story` and `plan` is given. No planner run for the item is not found, as for a story.

It waits for no task: nothing of the dashboard is needed to build or test it.

## Done when

- [x] `activity.document` answers the planner's document, an empty one when it is missing, and refuses an unknown kind, with tests in `flai/internal/hostapi`
- [x] `agent.stream` with `plan` reads the item's newest planner run, refuses both or neither of `story` and `plan`, and answers not found with no run, with tests in `flai/internal/serve` and `flai/internal/hostapi`
- [x] `go test -race -short` passes for the packages it changed

## Notes
