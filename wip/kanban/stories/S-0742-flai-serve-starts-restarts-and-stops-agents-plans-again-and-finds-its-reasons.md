---
id: S-0742
type: story
nature: feature
title: flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:59:30Z
updated: 2026-11-12T15:59:30Z
transitions: []
tags: [cli]
topics: [cli, orchestration, planning]
touches: [flai/internal/serve/agents_launch.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_activity.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/serve/stop.go, flai/internal/serve/readyempty.go, flai/internal/serve/replan.go, flai/internal/serve/plan.go, flai/internal/serve/role.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/analyze.go, flai/internal/serve/curate.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/formula_schedule.go, flai/internal/serve/usage.go, flai/internal/serve/budget.go, flai/internal/serve/confirm.go, flai/internal/serve/remind.go, flai/internal/serve/commit.go, flai/internal/serve/claudecheck.go, flai/internal/serve/components.go, flai/internal/serve/checks.go, flai/internal/serve/verifyjob.go, flai/internal/serve/upgrade.go, flai/internal/serve/agents_test.go, flai/internal/serve/readyempty_test.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0742 flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name

## Goal

flai serve runs a story's agent when the story makes the transition whose effect is the pull, not when it enters `ready`, and keeps its other runs on the kind of the state each item is in. An operator's workflow with a state of its own between `in-progress` and `review` keeps the agent running there, starts nothing for it that the baseline does not, and is starved, reminded, and replanned as the baseline is.

## Acceptance criteria

- [ ] The files in `flai/internal/serve` that start, restart, stop, and measure a story's agent, start the strategic roles, find the orchestrator's reasons, open the starved-board thread, plan again, and remind name no state: each reads a state's kind or a transition's effect from the workflow
- [ ] A story's agent starts on the transition with the pull effect, restarts while its story is in a state of the active-work kind, and is stopped or left to end as today when it leaves that kind
- [ ] The starved-board thread (ADR-0285) and the `ready-empty` event count the stories in the queue state the pull takes from and the WIP limit of the active-work states
- [ ] Tests run the starts, the restarts, the ready-empty count, and the orchestrator's reasons against a workflow with a state renamed and one added between `in-progress` and `review`, with the outcome the kinds give
- [ ] With the baseline, every existing `serve` test passes unchanged

## Tasks

## Notes

- `serve` holds about 100 references to a state's name in 32 files. This story takes the agents' and the strategic roles' files; the kluster, forge, and merge queue files are a story of their own, and `replan.go`'s and `orchestrate_reasons.go`'s use of the role conditions changes with the role condition story.
