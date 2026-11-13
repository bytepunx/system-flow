---
id: T-3216
type: task
nature: feature
title: The strategic roles' starts and the orchestrator's reasons and hold read a state's kind and a transition's effect
status: backlog
parent: S-0742
owner: alex
created: 2026-11-12T16:07:05Z
updated: 2026-11-12T16:07:05Z
transitions: []
stream: S-0742
tags: [cli]
touches: [flai/internal/serve/role.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/analyze.go, flai/internal/serve/curate.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/formula_schedule.go]
---
# T-3216 The strategic roles' starts and the orchestrator's reasons and hold read a state's kind and a transition's effect

## Work

Replace the references to a state's name in the strategic roles' files: the orchestrator's reasons name a move to the wait state the acceptance leaves from where they name `review` today, the curator starts on a transition with the acceptance effect, and a formula's schedule makes its instance enter the state the pull takes from. Leave the role conditions in `orchestrate_reasons.go` to S-0748. Test the reasons on a workflow with a state renamed and one added.

## Done when

These files name no state, and `orchestrate_reasons_test.go` passes on both workflows.

## Notes
