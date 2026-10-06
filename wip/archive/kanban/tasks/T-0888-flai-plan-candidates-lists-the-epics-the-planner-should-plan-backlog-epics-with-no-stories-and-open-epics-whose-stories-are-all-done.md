---
id: T-0888
type: task
nature: feature
title: "flai plan --candidates lists the epics the planner should plan: backlog epics with no stories, and open epics whose stories are all done"
status: done
parent: S-0219
owner: alex
created: 2026-10-05T04:46:28Z
updated: 2026-10-06T03:03:30Z
transitions:
  - to: ready
    at: 2026-10-06T02:58:11Z
    by: agent-S-0219
  - to: in-progress
    at: 2026-10-06T02:58:12Z
    by: agent-S-0219
  - to: done
    at: 2026-10-06T03:03:30Z
    by: agent-S-0219
stream: S-0219
tags: [flai]
touches: [flai/internal/workitem/plancandidates.go, flai/internal/workitem/plancandidates_test.go, flai/cmd/plan.go, flai/cmd/plan_candidates_test.go]
usage:
  source: log
  seconds: 318
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 37
      output: 13928
      cache_read: 2113063
      cache_write: 48732
      cost: 0.9779
---
# T-0888 flai plan --candidates lists the epics the planner should plan: backlog epics with no stories, and open epics whose stories are all done

## Work

The first criterion, with `plan_backlog_epics`, needs the list of epics to plan worked out by flai, not by the orchestrator, as the designer chose for S-0217's commands.

- Add `flai/internal/workitem/plancandidates.go`: an epic is a candidate when it is in the backlog with no story, or when it is not done or cancelled and every story under it is done or cancelled (at least one done). Leave out an epic a planner runs for now, when the caller says which (`serve`'s runs under `plans`), and an epic whose newest planner run ended `asked` with its thread still awaiting the operator.
- Add `flai plan --candidates [--json]` to `flai/cmd/plan.go`: it lists the candidates, each with the reason it is one, and writes nothing. `flai plan <id>` is unchanged.

This task waits for no other task of this story. It runs with the other first-layer tasks, whose paths it does not share.

## Done when

- fixture tests pin a backlog epic with no stories and an epic whose stories are all done as candidates, and leave out an epic with an open story, a done epic, and a cancelled one
- `flai plan --candidates` changes no file, and `--json` gives each epic's ID and reason
- `go test ./internal/workitem/ ./cmd/` passes

## Notes
