---
id: S-0796
type: story
nature: remediation
title: The acceptance preview lists verify_behind_main while the merge queue batches, so the review page disables Accept and offers Rebase and resolve for a story the acceptance would enqueue
status: backlog
owner: alex
created: 2026-11-16T00:41:47Z
updated: 2026-11-16T00:41:47Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: flai
    at: 2026-11-16T00:41:47Z
---
# S-0796 The acceptance preview lists verify_behind_main while the merge queue batches, so the review page disables Accept and offers Rebase and resolve for a story the acceptance would enqueue

## Goal

This story remediates [I-0381](../../../design/issues/I-0381-the-acceptance-preview-lists-verify-behind-main-while-the-merge-queue-batches.md), "The acceptance preview lists verify_behind_main while the merge queue batches, so the review page disables Accept and offers Rebase and resolve for a story the acceptance would enqueue". The issue recommends this solution:

While the merge queue batches and is not held, the operator's acceptance preview does not list `verify_behind_main` as a blocker: `verifyBehindMainGate` in `flai/internal/preview/accept.go` skips when `opts.QueueBatches` is set and the acceptance would enqueue the story, as `queuesInstead` decides, and still reports the stale tiers so the preview can say the batch will run them. The review page then offers Accept, which enqueues, in place of Rebase and resolve. An acceptance that lands by hand (`--without-queue-tiers`, `--without-kluster`, `--keep-removal`, or a held queue) keeps the gate. Refine ADR-0291 with an ADR that says so.

## Acceptance criteria
- [ ] The cause I-0381 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0381 is closed with `flai issue close I-0381 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0381. time_lost_per_cycle 2h carried over from I-0381's Impact section. Its Impact time_lost_per_cycle was taken rather than the 1h derived from its cost and count.
