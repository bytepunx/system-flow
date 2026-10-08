---
id: S-0297
type: story
nature: improvement
title: Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite
status: done
owner: alex
created: 2026-10-06T19:46:46Z
updated: 2026-10-08T10:01:20Z
transitions:
  - to: ready
    at: 2026-10-08T08:50:34Z
    by: alex
  - to: in-progress
    at: 2026-10-08T09:42:11Z
    by: agent-S-0297
  - to: review
    at: 2026-10-08T10:00:38Z
    by: agent-S-0297
  - to: done
    at: 2026-10-08T10:01:20Z
    by: orchestrator
tags: [cli]
topics: [cli, git, continuous-improvement]
touches: [flai/cmd/stream_sync_test.go, design/issues/I-0092-two-story-branches-that-each-bump-the-same-issue-conflict-in-its-front-matter-whose-count-last-reported-and-updated-lines-both-rewrite.md, design/issues/summary.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1123
  turns:
    - day: 2026-10-08
      hand_edits: 2
      work: 34
  models:
    - model: claude-opus-5-5
      input: 96
      output: 23323
      cache_read: 5161692
      cache_write: 224774
      cost: 3.1143
  strategic:
    - kind: orchestrator
      seconds: 38
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 17
          output: 232
          cache_read: 5008682
          cache_write: 6369
          cost: 1.2376
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-06T19:46:46Z
  value: 12.5
  by: planner-S-0297
  at: 2026-10-07T02:18:55Z
finalized:
  by: alex
  at: 2026-10-07T02:18:22Z
---
# S-0297 Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite

## Goal

This story remediates [I-0092](../../../design/issues/I-0092-two-story-branches-that-each-bump-the-same-issue-conflict-in-its-front-matter-whose-count-last-reported-and-updated-lines-both-rewrite.md), "Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0092 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0092 is closed with `flai issue close I-0092 --reason` saying what fixed it

## Tasks
- T-1428 Tests reproduce I-0092's instances not yet covered: an acceptance whose rebase meets main's bump of the same issue, and a sync whose bump meets main's close of it
- T-1429 I-0092 is closed with flai issue close, naming S-0326's merge and the tests that reproduce its instances

## Notes

Cost of delay inputs set by flai from I-0092. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T19:46:45Z, 0 days before this story; under one cycle counts as one).

### Accepted by the orchestrator

- Verified: 45a0db97c318b218dda4110994ac8ce5e01fb2e4
- At: 2026-10-08T10:01:20Z

Verdict: accept. flai verify passed every step at the branch head 45a0db97, and the verifier matched both criteria to the diff. It is tests only: S-0326's issues.Merge is the fix, and these tests reproduce I-0092's remaining instances.
- 1: flai/cmd/stream_sync_test.go
- 2: design/issues/I-0092-two-story-branches-that-each-bump-the-same-issue-conflict-in-its-front-matter-whose-count-last-reported-and-updated-lines-both-rewrite.md, design/issues/summary.md
