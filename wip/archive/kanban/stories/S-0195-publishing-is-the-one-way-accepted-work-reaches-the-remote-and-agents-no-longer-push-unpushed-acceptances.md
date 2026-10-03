---
id: S-0195
type: story
nature: improvement
title: Publishing is the one way accepted work reaches the remote, and agents no longer push unpushed acceptances
status: done
owner: arobson
created: 2026-10-02T09:46:13Z
updated: 2026-10-03T00:38:25Z
transitions:
  - to: ready
    at: 2026-10-02T23:24:51Z
    by: alex
  - to: in-progress
    at: 2026-10-02T23:25:20Z
    by: agent-S-0195
  - to: review
    at: 2026-10-03T00:34:42Z
    by: agent-S-0195
  - to: done
    at: 2026-10-03T00:38:25Z
    by: alex
tags: [flai, dashboard]
topics: [release]
touches: [flai/internal/mcpserver, flai/cmd, flai/internal/harness, flai/internal/hostapi, flaiover/src, flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/release/remote.go, flai/internal/release/remote_test.go, flai/internal/preview, design/system, docs, design/conventions/git.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4195
  models:
    - model: claude-haiku-4-5-20251001
      input: 202
      output: 10545
      cache_read: 1456053
      cache_write: 89560
      cost: 0.3105
    - model: claude-opus-5-5
      input: 804
      output: 222024
      cache_read: 42334654
      cache_write: 847473
      cost: 17.8465
    - model: claude-sonnet-5-5
      input: 28
      output: 6278
      cache_read: 573401
      cache_write: 57904
      cost: 0.3223
---
# S-0195 Publishing is the one way accepted work reaches the remote, and agents no longer push unpushed acceptances

## Goal

[ADR-0067](../../../design/adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md) (proposed 2026-10-02): accepted work reaches the remote only when it is published (`git fetch`, then `flai release --pending`, or the board's Publish), and agents publish only when the operator asks. The work-management convention already says so. Today `inbox` and `flai board` report `unpushed`, and the MCP server's instructions (`flai/internal/mcpserver/folder.go` ~195 and the single-project server's) tell agents to run `flai push --pending` before anything else (S-0063); `flai push --pending` and the `auto-publish` host action (ADR-0048) are a second way to the remote.

## Acceptance criteria
- [x] The designer accepts ADR-0067, or the story records the changes they ask for in a new ADR, before anything else is built
- [x] The MCP servers' instructions, `inbox`, `flai board`, and `harness.Prompt` no longer tell an agent to push an unpushed acceptance; the board and the dashboard still show what is accepted and not yet published
- [x] `flai push --pending` and the `auto-publish` host action remains but is not integrated into the dashboard
- [x] the push banner over the top of the board is removed entirely
- [x] `flai release --pending` and Publish fetch first, or refuse with what to run when the clone is behind, as S-0174 does for tags
- [x] `design/system/pushing-from-the-board.md`, `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `design/system/workflow.md`, and the operator and user guides describe publishing as the one way to the remote
- [x] Tests cover the instructions and inbox without the push duty, and publishing after a fetch

## Tasks
- T-0717 ADR-0067 is accepted, or the designer's changes are a new ADR
- T-0718 Inbox, flai board, the MCP instructions, and the harness prompt no longer tell an agent to push
- T-0719 flai release --pending and Publish refuse a clone behind its remote, naming the fetch and merge to run
- T-0720 The host API no longer offers the dashboard pushing or auto-publish
- T-0721 The dashboard has no push banner, and Publish shows a clone behind its remote
- T-0722 The design and the guides describe publishing as the one way to the remote

## Notes

"Fetch first, or refuse": flai refuses, naming the fetch and the merge or rebase, as S-0174 does for tags; flai still never fetches. `flai push --pending` and `auto-publish` stay as the operator's shell tools; the dashboard neither offers pushing nor shows auto-publish. The designer's refinements on TH-0070 (drop tags made by a publish refused mid-way, rebase or merge, re-tag; conflicts through tasks and threads) are S-0242.

ADR-0067 is proposed because its details are the drafting agent's reading of the designer's instruction "update publishing and releasing".
