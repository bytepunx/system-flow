---
id: S-0599
type: story
nature: improvement
title: flai serve's replanner writes every forecast it moves in one pass with one check, so an acceptance's replan lands in seconds, not 80 s
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:39:17Z
updated: 2026-11-07T22:39:17Z
transitions: []
tags: [cli]
touches: [flai/internal/serve/replan.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0599 flai serve's replanner writes every forecast it moves in one pass with one check, so an acceptance's replan lands in seconds, not 80 s

## Goal

Cause 9 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): the replanner writes each moved forecast through `itemedit.Apply`, one after another (`flai/internal/serve/replan.go`), so each pays a checked edit's two whole-repository checks, 5 to 6.4 s a forecast. Over the week to 2026-11-07, 98 replans moved 1,200 forecasts, about 1.7 to 2 h of serve CPU, and an acceptance's replan commit landed p50 80 s and p90 127 s after the acceptance.

## Acceptance criteria

- [ ] A replan writes every forecast it moves in one pass under the items lock, with one check over the written items and one write of `wip/agents/index.md`.
- [ ] A replan that moves 20 forecasts in a scratch copy of this repository takes under 5 s, and server-performance.md records the before and after.
- [ ] A test shows a replan writes the same forecasts as before, and that a forecast another writer changed meanwhile is merged or refused as ADR-0184 says.

## Tasks

## Notes

- Proposed fix: an `itemedit` entry that applies many forecast changes at once. Forecasts are values flai owns, so the check can be the work-item rules over the written items.
- Independent of the story for cause 8, but it gains from it: with a scoped checked edit alone, 20 forecasts would still cost about 40 s.
