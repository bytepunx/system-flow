---
id: S-0247
type: story
nature: remediation
title: The dashboard drops flai's archived flag from board cards, so the done lane shows archived items it should hide while the clone lags its remote's tags
status: done
owner: alex
created: 2026-10-03T17:49:41Z
updated: 2026-10-04T23:13:46Z
transitions:
  - to: ready
    at: 2026-10-04T21:42:34Z
    by: alex
  - to: in-progress
    at: 2026-10-04T21:42:51Z
    by: agent-S-0247
  - to: review
    at: 2026-10-04T21:56:33Z
    by: agent-S-0247
  - to: done
    at: 2026-10-04T23:13:46Z
    by: alex
tags: []
touches: [design/issues/I-0060-the-dashboard-drops-flai-s-archived-flag-from-board-cards-so-the-done-lane-shows-archived-items-it-should-hide-while-the-clone-lags-its-remote-s-tags.md, design/issues/summary.md, flaiover/src/lib/server/board.ts, flaiover/src/lib/server/repo-channel.test.ts, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 843
  models:
    - model: claude-opus-5-5
      input: 122
      output: 22609
      cache_read: 4740084
      cache_write: 149043
      cost: 2.4779
    - model: claude-sonnet-5
      input: 52
      output: 8110
      cache_read: 1034346
      cache_write: 63117
      cost: 0.4459
---
# S-0247 The dashboard drops flai's archived flag from board cards, so the done lane shows archived items it should hide while the clone lags its remote's tags

## Goal

This story remediates [I-0060](../../../design/issues/I-0060-the-dashboard-drops-flai-s-archived-flag-from-board-cards-so-the-done-lane-shows-archived-items-it-should-hide-while-the-clone-lags-its-remote-s-tags.md), "The dashboard drops flai's archived flag from board cards, so the done lane shows archived items it should hide while the clone lags its remote's tags". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0060 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0060 is closed with `flai issue close I-0060 --reason` saying what fixed it

## Tasks
- T-0829 board.ts carries flai's archived flag onto each card, with a test through board()

## Notes
