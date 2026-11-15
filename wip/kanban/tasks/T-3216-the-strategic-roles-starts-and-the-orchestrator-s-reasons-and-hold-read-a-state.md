---
id: T-3216
type: task
nature: feature
title: The strategic roles' starts and the orchestrator's reasons and hold read a state's kind and a transition's effect
status: done
parent: S-0742
owner: alex
created: 2026-11-12T16:07:05Z
updated: 2026-11-15T01:55:09Z
transitions:
  - to: ready
    at: 2026-11-15T01:27:02Z
    by: agent-S-0742
  - to: in-progress
    at: 2026-11-15T01:27:02Z
    by: agent-S-0742
  - to: done
    at: 2026-11-15T01:55:09Z
    by: agent-S-0742
stream: S-0742
tags: [cli]
touches: [flai/internal/serve/role.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/role_events.go]
after: [T-3370]
usage:
  source: log
  seconds: 1687
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 12949
      cache_read: 4391561
      cache_write: 100919
      cost: 1.7197
---
# T-3216 The strategic roles' starts and the orchestrator's reasons and hold read a state's kind and a transition's effect

## Work

Replace the references to a state's name in the strategic roles' files: the orchestrator's reasons name a move to the wait state the acceptance leaves from where they name `review` today, the curator starts on a transition with the acceptance effect, and a formula's schedule makes its instance enter the state the pull takes from. Leave the role conditions in `orchestrate_reasons.go` to S-0748. Test the reasons on a workflow with a state renamed and one added.

## Done when

These files name no state, and `orchestrate_reasons_test.go` passes on both workflows.

## Notes
