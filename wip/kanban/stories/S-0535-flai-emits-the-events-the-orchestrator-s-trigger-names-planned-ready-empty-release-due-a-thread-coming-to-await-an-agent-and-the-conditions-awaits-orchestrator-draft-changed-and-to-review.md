---
id: S-0535
type: story
nature: feature
title: "flai emits the events the orchestrator's trigger names: planned, ready-empty, release-due, a thread coming to await an agent, and the conditions awaits-orchestrator, draft-changed, and to-review"
status: backlog
parent: E-0030
owner: alex
created: 2026-10-28T21:02:24Z
updated: 2026-10-28T21:02:24Z
transitions: []
tags: [flai, serve, events]
topics: [orchestration, cli]
touches: [flai/internal/serve, flai/internal/mcpserver, flai/internal/workitem, design/system/flai-cli.md, design/workflows/roles/README.md, design/adrs]
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
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 5
          cache_read: 186500
          cache_write: 1871
          cost: 0.0415
---
# S-0535 flai emits the events the orchestrator's trigger names: planned, ready-empty, release-due, a thread coming to await an agent, and the conditions awaits-orchestrator, draft-changed, and to-review

## Goal

ADR-0209 §1 lists the reasons the orchestrator is started for in ADR-0207's trigger vocabulary. Some are events flai emits today with a condition to name (`moved` to review, `edited` on a draft, `reordered`, `accepted`); others are not events at all: a planner's run ending with its items written, the ready column draining below the in-progress limit, the release policy's figures changing with an acceptance, a thread coming to await the orchestrator or the operator (ADR-0181). This story defines each as an event or a named condition flai implements (ADR-0207 §3), emits it from where it happens, and reports it through `wait_for_events` and `inbox` like the others, so that the following story can name them in `orchestrator.md`'s `trigger` and flai serve can queue them.

## Acceptance criteria
- [ ] `planned` is an event, with the item planned and the IDs written, emitted when a planner's run ends with its items written, and reported by `wait_for_events` and `inbox`
- [ ] `ready-empty` is an event emitted when the ready column drains below the in-progress limit's room, once per draining
- [ ] `release-due` is an event emitted when an acceptance changes what `flai release --evaluate` would answer
- [ ] `awaits-orchestrator`, `draft-changed`, and `to-review` are conditions `flai check` accepts on `thread`, `edited`, and `moved` triggers, each implemented and tested
- [ ] A `thread` event exists for a thread coming to await an agent, with the thread and whom it awaits, under ADR-0181's rule
- [ ] The design (`flai-cli.md`'s event table, `design/workflows/roles/README.md`'s trigger table) and the MCP descriptions name each event and condition

## Tasks

## Notes

A named condition is one flai implements, not one a role file defines (ADR-0207 §3). The names are ADR-0209's proposals; the story may rename them with the ADR's text if a better name fits, before the ADR is accepted.
