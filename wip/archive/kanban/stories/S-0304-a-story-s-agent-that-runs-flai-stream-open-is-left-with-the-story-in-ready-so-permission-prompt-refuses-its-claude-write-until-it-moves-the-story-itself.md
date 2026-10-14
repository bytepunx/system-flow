---
id: S-0304
type: story
nature: improvement
title: A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself
status: done
owner: alex
created: 2026-10-07T01:07:09Z
updated: 2026-10-14T21:34:18Z
transitions:
  - to: ready
    at: 2026-10-14T20:54:54Z
    by: alex
  - to: in-progress
    at: 2026-10-14T20:55:01Z
    by: agent-S-0304
  - to: review
    at: 2026-10-14T21:33:29Z
    by: agent-S-0304
  - to: done
    at: 2026-10-14T21:34:18Z
    by: orchestrator
tags: []
touches: [design/adrs/0153-flai-stream-open-starts-a-ready-story-as-flai-story-start-does-moving-it-to-in.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/move_cancel_test.go, flai/cmd/stream.go, flai/cmd/stream_open_test.go, flai/internal/storystart/start.go, design/issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md, design/issues/summary.md, design/issues/I-0097-a-story-s-agent-that-runs-flai-stream-open-is-left-with-the-story-in-ready-so-permission-prompt-refuses-its-claude-write-until-it-moves-the-story-itself.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2322
  priced_by: harness
  turns:
    - day: 2026-10-14
      ceremony: 3
      hand_edits: 2
      work: 43
  models:
    - model: claude-opus-5-5
      input: 154
      output: 42587
      cache_read: 9365335
      cache_write: 490069
      cost: 5.7226
  strategic:
    - kind: orchestrator
      seconds: 144
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 139
          cache_read: 877387
          cache_write: 2463
          cost: 0.2206
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-10-07T01:07:09Z
finalized:
  by: alex
  at: 2026-10-07T02:18:11Z
---
# S-0304 A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself

## Goal

This story remediates [I-0097](../../../design/issues/I-0097-a-story-s-agent-that-runs-flai-stream-open-is-left-with-the-story-in-ready-so-permission-prompt-refuses-its-claude-write-until-it-moves-the-story-itself.md), "A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0097 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0097 is closed with `flai issue close I-0097 --reason` saying what fixed it

## Tasks
- T-1572 flai stream open starts a ready story as flai story start does, moving it to in-progress and refusing one the board holds
- T-1573 Close I-0097 naming the fix

## Notes

Cost of delay inputs set by flai from I-0097. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T23:28:35Z, 0.1 days before this story; under one cycle counts as one).

### Accepted by the orchestrator

- Verified: c30d92033a73f4ee0921375a888d451190bb59f5
- At: 2026-10-14T21:34:18Z

Verdict: meets all criteria; flai verify passed every step at c30d9203; every changed path is within the story's touches, none under .claude/.

- 1: flai/cmd/stream.go and flai/internal/storystart/start.go (flai stream open starts a ready story through storystart.Begin, refusing a held one), flai/cmd/stream_open_test.go, flai/cmd/move_cancel_test.go, design/adrs/0153-flai-stream-open-starts-a-ready-story-as-flai-story-start-does-moving-it-to-in.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md
- 2: design/issues/I-0097 (closed, naming storystart.Begin and ADR-0153), design/issues/summary.md
