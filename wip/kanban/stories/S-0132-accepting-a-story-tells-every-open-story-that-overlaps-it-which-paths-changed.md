---
id: S-0132
type: story
nature: feature
title: Accepting a story tells every open story that overlaps it which paths changed
status: backlog
parent: E-0009
owner: alex
created: 2026-09-26T07:59:23Z
updated: 2026-09-26T07:59:23Z
transitions: []
tags: [cli]
touches: [flai/cmd, flai/internal/mcpserver, docs, design/system]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0132 Accepting a story tells every open story that overlaps it which paths changed

## Goal

Build the notice at acceptance of [ADR-0046](../../../design/adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md): when a story is accepted, every story still in progress or in review whose claim overlaps what the accepted story changed is told which paths changed, so its agent syncs and checks its work against them before it is surprised at its own acceptance.

## Acceptance criteria

- [ ] `flai accept` (and acceptance from the board) records, for each open story whose claim overlaps the accepted change's paths, a change that the MCP `inbox` and `wait_for_events` report to that story's agent, naming the accepted story and the overlapping paths.
- [ ] A story whose claim does not overlap the accepted change gets no such notice.
- [ ] The agent convention in `design/conventions/work-management.md` (via `template/` first) says to run `flai stream sync` and re-run the tests on such a notice.
- [ ] Tests cover an overlapping and a non-overlapping open story; `make test` and lint pass; `design/system/workflow.md`, `flai-cli.md`, and `docs/users/flai.md` describe it.

## Tasks

## Notes
