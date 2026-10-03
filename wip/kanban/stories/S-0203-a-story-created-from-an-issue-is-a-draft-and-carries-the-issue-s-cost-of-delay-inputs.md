---
id: S-0203
type: story
nature: improvement
title: A story created from an issue is a draft and carries the issue's cost of delay inputs
status: review
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:12Z
updated: 2026-10-03T18:29:46Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:14Z
    by: alex
  - to: in-progress
    at: 2026-10-03T18:04:17Z
    by: agent-S-0203
  - to: review
    at: 2026-10-03T18:29:46Z
    by: agent-S-0203
tags: [flai]
touches: [flai/internal/issues, flai/cmd/issue.go, flai/internal/mcpserver, flai/internal/itemnew, flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/cmd/issue_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/continuous-improvement.md, design/system/flai-cli.md, design/issues/I-0064-flai-stream-sync-s-trial-merge-blames-a-story-for-conflicts-between-main-and-another-story-s-stale-branch.md, design/issues/I-0065-flai-issue-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-issue-number.md, design/issues/summary.md]
after: [S-0199, S-0198]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1555
  models:
    - model: claude-haiku-4-5-20251001
      input: 162
      output: 8987
      cache_read: 929807
      cache_write: 74864
      cost: 0.2317
    - model: claude-opus-5-5
      input: 308
      output: 117048
      cache_read: 12681949
      cache_write: 456633
      cost: 7.6804
    - model: claude-sonnet-5
      input: 60
      output: 15436
      cache_read: 2129383
      cache_write: 121641
      cost: 0.8845
---
# S-0203 A story created from an issue is a draft and carries the issue's cost of delay inputs

## Goal

S-0198 gives flai a step that creates a remediation or improvement story from an issue. Such a story is an agent's draft until the operator finalizes it, and when the issue records an impact (the analyzer's, or a `cost`), the story should carry it as cost of delay inputs so the planner and orchestrator can rank it.

## Acceptance criteria
- [x] The story `flai issue story I-nnnn` (and the MCP equivalent, and S-0198's warning flow) creates has `draft: true`
- [x] When the issue carries `cost` (time per occurrence) and `count`, the story's `cost_of_delay.inputs.time_lost_per_cycle` is set from them, with `by: flai` and the derivation in the story's Notes; an issue the analyzer wrote with an `impact` section carries its revenue or penalty figures over
- [x] The story links the issue and the issue's Remediation section links the story
- [x] Tests cover the draft flag and the inputs carried over; the user guide and `design/system/continuous-improvement.md` say so

## Tasks
- T-0760 issues drafts a story with cost of delay inputs and links it from the Remediation section
- T-0761 workitem.NewOptions takes a cost of delay and itemnew commits what goes with the new item
- T-0762 flai issue story and MCP issue_story make a draft story with the issue's cost of delay inputs and link it from the issue
- T-0763 User guide and continuous-improvement design say a story from an issue is a draft carrying its cost of delay inputs

## Notes

- The dashboard's issue-to-story flow at acceptance runs `flai issue story --autocommit`, so it makes the same draft and commits the issue with it.
- No `## Impact` format existed: this story reads `- revenue_per_week: <amount>`, `- penalty_per_week: <amount>`, and `- time_lost_per_cycle: <duration>` lines, the keys of `cost_of_delay.inputs`, for S-0224's analyzer to write.
- The issue names the story by ID in its Remediation section, not by a path link, since the story's file moves on acceptance.
