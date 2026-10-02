---
id: S-0209
type: story
nature: feature
title: The planner drafts an epic's stories into the backlog and revisits the children it already has
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:14Z
updated: 2026-10-02T11:54:40Z
transitions: []
tags: [flai]
touches: [flai/internal/harness, ".claude/agents/planner.md", template/, design/system/strategic-agents.md]
after: [S-0208]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0209 The planner drafts an epic's stories into the backlog and revisits the children it already has

## Goal

Given an epic, the planner writes the stories that deliver its outcome, as drafts in the backlog, and when the epic already has stories it revisits each open one: re-enriches it, proposes splits or merges, and reports what changed.

## Acceptance criteria
- [ ] For an epic with no stories, the planner writes stories with goal, acceptance criteria, nature, tags, topics, touches, and `after` between them, each `draft: true` in backlog, and opens one thread on the epic summarising the plan: the stories, their order, and the assumptions it made
- [ ] For an epic with stories, it revisits every story not done or cancelled: re-enriches it (the enrichment story), and proposes in the epic's thread any story it would split, merge, add, or drop, creating drafts for additions but never cancelling or rewriting a finalized story's words without asking
- [ ] Stories it writes pass `flai check --strict` and the markdown lint before they are kept
- [ ] Its log entry names the stories created and revisited and the run's cost
- [ ] Measured on E-0016 itself or a comparable epic: the designer judges the drafts on a thread and the result is recorded in `design/system/strategic-agents.md`

## Tasks

## Notes
