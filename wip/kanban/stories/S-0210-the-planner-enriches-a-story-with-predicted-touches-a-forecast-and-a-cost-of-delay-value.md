---
id: S-0210
type: story
nature: feature
title: The planner enriches a story with predicted touches, a forecast, and a cost of delay value
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:14Z
updated: 2026-10-03T20:33:46Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:46Z
    by: alex
tags: [flai]
touches: [flai/internal/harness, flai/internal/metrics, flai/cmd, flai/internal/storygit, ".claude/agents/planner.md", design/system/strategic-agents.md]
after: [S-0208]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0210 The planner enriches a story with predicted touches, a forecast, and a cost of delay value

## Goal

Enriching a story gives the orchestrator what it schedules by and the launcher what it holds by: an accurate list of files the story will touch, a wall-clock forecast and delivery date, and a cost of delay value from the operator's inputs.

## Acceptance criteria
- [ ] Touches: the planner predicts the files and folders from the story's goal and criteria, the design documents it links, the code layout, and co-change history (`flai touches suggest S-nnnn`, a deterministic command that lists files often changed with the ones named, from git history with flai's bookkeeping commits filtered), and writes them as `touches`, keeping any the story already declares; its Notes say which came from where
- [ ] Forecast: `flai forecast S-nnnn` computes `duration` and `delivery` deterministically from past stories' `usage.seconds` and cycle times by nature, model, and size (criteria and touches count), the story's position in the pull order, the in-progress limit, and the forecasts of the stories ahead; the planner reviews it, adjusts with a stated reason, and writes `forecast` with `by: planner`; without enough history the basis says so and uses the project's defaults (`planning.default_duration`)
- [ ] Cost of delay: `flai cod S-nnnn` derives `value` per week from the item's inputs (revenue and penalty as given; time lost per cycle × the manifest's hour rate × cycles per week), or apportions the epic's value over its open stories by forecast duration when the story has none; the planner writes `value`, `by`, `at`
- [ ] The three commands are usable without the planner; the planner's prompt uses them and records its reasoning in the story's Notes under a `### Planning` heading it owns
- [ ] `design/system/strategic-agents.md`, `flai-cli.md`, and the user guide describe the commands and the model; tests pin the forecast and cost of delay on fixtures

## Tasks

## Notes
