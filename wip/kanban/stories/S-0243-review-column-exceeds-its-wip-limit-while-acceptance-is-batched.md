---
id: S-0243
type: story
nature: improvement
title: Review column exceeds its WIP limit while acceptance is batched
status: ready
owner: arobson
created: 2026-10-03T01:46:35Z
updated: 2026-10-03T01:52:11Z
transitions:
  - to: ready
    at: 2026-10-03T01:52:11Z
    by: alex
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0243 Review column exceeds its WIP limit while acceptance is batched

## Goal

This story remediates [I-0007](../../../design/issues/I-0007-review-queue-exceeds-limit.md), "Review column exceeds its WIP limit while acceptance is batched". The issue recommends this solution:

Either raise the review limit in `wip/kanban/board.md` to match the operator's acceptance cadence, or accept stories before pulling the next one. Operator's call; see the S-0026 open questions.

## Acceptance criteria
- [ ] The cause I-0007 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0007 is closed with `flai issue close I-0007 --reason` saying what fixed it

## Tasks

## Notes
