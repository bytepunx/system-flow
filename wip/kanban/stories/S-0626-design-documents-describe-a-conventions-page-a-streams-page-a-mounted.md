---
id: S-0626
type: story
nature: improvement
title: Design documents describe a /conventions page, a /streams page, a mounted-repository readiness check, a flai OTEL span, a dashboard scan of wip/kanban, and a single lockfile that are not there
status: backlog
owner: alex
created: 2026-11-07T23:40:21Z
updated: 2026-11-07T23:40:21Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
draft: true
---
# S-0626 Design documents describe a /conventions page, a /streams page, a mounted-repository readiness check, a flai OTEL span, a dashboard scan of wip/kanban, and a single lockfile that are not there

## Goal

This story remediates [I-0293](../../../design/issues/I-0293-design-documents-describe-a-conventions-page-a-streams-page-a-mounted.md), "Design documents describe a /conventions page, a /streams page, a mounted-repository readiness check, a flai OTEL span, a dashboard scan of wip/kanban, and a single lockfile that are not there", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0293 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0293 is closed with `flai issue close I-0293 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0293: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0293: analyzer 42s, 0.8156 USD.
