---
id: S-0132
type: story
nature: feature
title: Accepting a story tells every open story that overlaps it which paths changed
status: done
parent: E-0009
owner: alex
created: 2026-09-26T07:59:23Z
updated: 2026-09-26T21:09:40Z
transitions:
  - to: ready
    at: 2026-09-26T17:50:41Z
    by: alex
  - to: in-progress
    at: 2026-09-26T18:15:50Z
    by: agent-S-0132
  - to: review
    at: 2026-09-26T18:26:29Z
    by: agent-S-0132
  - to: done
    at: 2026-09-26T21:09:40Z
    by: alex
tags: [cli]
touches: [flai/cmd, flai/internal/mcpserver, flai/internal/itemedit, docs, design/system, design/conventions, template/root/design/conventions]
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

- [x] `flai accept` (and acceptance from the board) records, for each open story whose claim overlaps the accepted change's paths, a change that the MCP `inbox` and `wait_for_events` report to that story's agent, naming the accepted story and the overlapping paths.
- [x] A story whose claim does not overlap the accepted change gets no such notice.
- [x] The agent convention in `design/conventions/work-management.md` (via `template/` first) says to run `flai stream sync` and re-run the tests on such a notice.
- [x] Tests cover an overlapping and a non-overlapping open story; `make test` and lint pass; `design/system/workflow.md`, `flai-cli.md`, and `docs/users/flai.md` describe it.

## Tasks
- T-0481 Acceptance records an overlap notice for each open story whose claim covers a changed path
- T-0482 MCP inbox and wait_for_events report an overlap notice as a change on the open story
- T-0483 The work-management convention says to sync and re-run the tests on an overlap notice
- T-0484 Docs describe the notice at acceptance; tests and lint pass

## Notes

- The notices live in `.flai-cache/overlaps.jsonl`, apart from `edits.jsonl`, so the installed flai 1.18.5 that runs the MCP server does not read them as edits. Agents see them once their MCP server is a flai with S-0132.
- Acceptance from the board runs `flai accept` through `accept.run`, so it records the same notices; no dashboard test covers it.
- `make smoke`'s template render, repository check, and installer steps pass. Its markdown lint fails on six files under `wip/` that this story does not touch (TH-0017, I-0027).
