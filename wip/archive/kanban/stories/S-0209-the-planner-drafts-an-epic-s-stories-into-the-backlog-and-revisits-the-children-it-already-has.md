---
id: S-0209
type: story
nature: feature
title: The planner drafts an epic's stories into the backlog and revisits the children it already has
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:14Z
updated: 2026-10-04T06:29:48Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:40Z
    by: alex
  - to: in-progress
    at: 2026-10-04T03:59:46Z
    by: system-flow
  - to: review
    at: 2026-10-04T05:37:24Z
    by: agent-S-0209
  - to: done
    at: 2026-10-04T06:29:48Z
    by: alex
tags: [flai]
touches: [flai/internal/harness, flai/internal/guard, flai/internal/mcpserver, flai/internal/serve, flai/cmd/guard.go, ".claude/agents/planner.md", template, design/conventions/strategic-agents.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users, design/issues/I-0058-flai-guard-refuses-a-sub-agent-s-shell-command-whose-heredoc-text-reads-like-a-flai-write.md, design/issues/summary.md, design/issues/I-0069-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md, design/issues/I-0071-item-edit-refuses-a-touch-that-starts-with-a-dot-though-flai-touches-accepts-it.md, design/issues/I-0065-flai-issue-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-issue-number.md]
after: [S-0208]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5905
  models:
    - model: claude-opus-5-5
      input: 610
      output: 137678
      cache_read: 52594115
      cache_write: 597631
      cost: 17.2286
    - model: claude-sonnet-5
      input: 78
      output: 21500
      cache_read: 2207456
      cache_write: 176697
      cost: 1.0984
---
# S-0209 The planner drafts an epic's stories into the backlog and revisits the children it already has

## Goal

Given an epic, the planner writes the stories that deliver its outcome, as drafts in the backlog, and when the epic already has stories it revisits each open one: re-enriches it, proposes splits or merges, and reports what changed.

## Acceptance criteria
- [x] For an epic with no stories, the planner writes stories with goal, acceptance criteria, nature, tags, topics, touches, and `after` between them, each `draft: true` in backlog, and opens one thread on the epic summarising the plan: the stories, their order, and the assumptions it made
- [x] For an epic with stories, it revisits every story not done or cancelled: re-enriches it (the enrichment story), and proposes in the epic's thread any story it would split, merge, add, or drop, creating drafts for additions but never cancelling or rewriting a finalized story's words without asking
- [x] Stories it writes pass `flai check --strict` and the markdown lint before they are kept
- [x] Its log entry names the stories created and revisited and the run's cost
- [x] Measured on E-0016 itself or a comparable epic: the designer judges the drafts on a thread and the result is recorded in `design/system/strategic-agents.md`

## Tasks
- T-0791 The planner's prompt asks for a plan thread, revisit proposals, and a summary naming the stories, and the guard holds its stories to drafts
- T-0792 item_new with a body is checked as flai story new --body-stdin is, so a planner's story passes flai check --strict and the lint before it is kept
- T-0793 A planner run's activity entry names the items it created and changed
- T-0794 The convention, design, docs, and template say how the planner drafts and revisits an epic's stories
- T-0795 planner.md says how the planner drafts and revisits an epic's stories
- T-0796 The designer judges a planner run on an epic, and the result is recorded in strategic-agents.md

## Notes

The planner's run on E-0016 (TH-0107) was made with S-0209's prompt over the host's flai 1.30.0, which predates the guard's draft rule, the checked `item_new`, and the entry's items; those are covered by tests. The designer judged it on TH-0103, recorded in `design/system/strategic-agents.md`. The run measured the revisit; drafting an epic with no stories is covered by the prompt, the guard, and their tests, since no open epic was empty. Template 1.0.43.
