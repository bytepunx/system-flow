---
id: S-0551
type: story
nature: improvement
title: A story starts on a trial merge instead of an overlap hold, and after holds a story for a dependency only, so the worktrees run stories in parallel as people on one remote do
status: backlog
parent: E-0031
owner: alex
created: 2026-10-30T01:52:25Z
updated: 2026-10-30T01:52:25Z
transitions: []
tags: [flai]
topics: [git, orchestration, planning, conventions]
touches: [flai/internal/serve, flai/internal/sync, flai/internal/claims, design/workflows/roles, design/system/workflow.md, design/system/strategic-agents.md, design/conventions/work-management.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0551 A story starts on a trial merge instead of an overlap hold, and after holds a story for a dependency only, so the worktrees run stories in parallel as people on one remote do

## Goal

The worktrees isolate every story, and the scheduler then serializes them: with an in-progress limit of 5, nothing was in progress for 42 percent of the fortnight, one story for 29 percent, two for 26 percent, three for 2 percent, and four or five never, mean 1.6 when anything ran. The journal records 124 overlap holds and 113 `after` holds against 62 real conflicts, the overlap holds on `check.go`, the strategic-agents design and convention, the manifest and dashboard designs, and the `after` chains planner-made, S-0417 to S-0422 each waiting on the one before, as I-0152 already noted. People sharing one remote merge optimistically and resolve the rare conflict; this story makes flai do the same. A ready story whose touches overlap an in-progress story's is started when a trial merge of the in-progress story's branch into a worktree at the ready story's base is clean, and held only when it is not, with the journal naming the conflicting file; `after` is honoured as today but the planner sets it only when a story reads what another writes, named in the story's notes, and `flai check` warns of an `after` with no reason; and the ready wait at p85, 5.4 hours in the fortnight, is a figure `flai stats` reports.

## Acceptance criteria
- [ ] `wait_for_work` and the orchestrator's pull start a story overlapping an in-progress one when a trial merge of the in-progress branch at the ready story's base is clean, and hold it with the conflicting path in the journal otherwise, with a test of each
- [ ] The planner's role file and the work-management convention say `after` names a story this one reads the output of, with the reason in the story's notes; `flai check` warns of an `after` without one, and the open chains S-0417 to S-0422 are reviewed against it
- [ ] `flai stats` reports stories in progress at once as a share of the window and the ready wait at p50 and p85, and the dashboard's flow page shows them
- [ ] workflow.md and strategic-agents.md describe the trial merge in place of the overlap hold, and ADR-0046's `after` is refined by a new ADR where the rule changes

## Tasks

The planner drafts them.

## Notes

Evidence: the board's transitions and the serve journal's `held` entries 2026-10-15 to 2026-10-29, counted in the operator's conversation; the 2026-10-15 bottlenecks analysis and I-0152 for the `after` chains.
