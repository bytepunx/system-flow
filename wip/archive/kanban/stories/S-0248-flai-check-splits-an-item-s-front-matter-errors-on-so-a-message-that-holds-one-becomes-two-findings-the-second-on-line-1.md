---
id: S-0248
type: story
nature: remediation
title: "flai check splits an item's front-matter errors on \"; \", so a message that holds one becomes two findings, the second on line 1"
status: done
owner: alex
created: 2026-10-03T17:49:41Z
updated: 2026-10-04T23:56:22Z
transitions:
  - to: ready
    at: 2026-10-04T21:42:42Z
    by: alex
  - to: in-progress
    at: 2026-10-04T23:14:13Z
    by: agent-S-0248
  - to: review
    at: 2026-10-04T23:28:14Z
    by: agent-S-0248
  - to: done
    at: 2026-10-04T23:56:22Z
    by: alex
tags: []
touches: [design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/I-0061-flai-check-splits-an-item-s-front-matter-errors-on-so-a-message-that-holds-one-becomes-two-findings-the-second-on-line-1.md, design/issues/summary.md, design/system/flai-cli.md, flai/internal/check/check.go, flai/internal/check/planning_test.go, flai/internal/manifest/agent.go, flai/internal/workitem/item.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 876
  models:
    - model: claude-opus-5-5
      input: 96
      output: 20566
      cache_read: 4874648
      cache_write: 132118
      cost: 2.4436
    - model: claude-sonnet-5
      input: 68
      output: 18825
      cache_read: 1748696
      cache_write: 117658
      cost: 0.8323
---
# S-0248 flai check splits an item's front-matter errors on "; ", so a message that holds one becomes two findings, the second on line 1

## Goal

This story remediates [I-0061](../../../design/issues/I-0061-flai-check-splits-an-item-s-front-matter-errors-on-so-a-message-that-holds-one-becomes-two-findings-the-second-on-line-1.md), "flai check splits an item's front-matter errors on "; ", so a message that holds one becomes two findings, the second on line 1". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0061 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0061 is closed with `flai issue close I-0061 --reason` saying what fixed it

## Tasks
- T-0830 Item.Validate returns its problems as a list that flai check ranges over

## Notes

Remedy, proposed from I-0061's instance and built in T-0830: `Item.Validate` returns `workitem.Problems`, a list of phrases, and `flai check` ranges over it, so a phrase that holds "; " stays one finding on its key's line.
