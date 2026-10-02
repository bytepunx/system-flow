---
id: S-0195
type: story
nature: improvement
title: Publishing is the one way accepted work reaches the remote, and agents no longer push unpushed acceptances
status: backlog
owner: arobson
created: 2026-10-02T09:46:13Z
updated: 2026-10-02T16:17:53Z
transitions: []
tags: [flai, dashboard]
topics: [release]
touches: [flai/internal/mcpserver, flai/cmd/push.go, flai/cmd/release.go, flai/internal/harness, flai/internal/hostapi, flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0195 Publishing is the one way accepted work reaches the remote, and agents no longer push unpushed acceptances

## Goal

[ADR-0067](../../../design/adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md) (proposed 2026-10-02): accepted work reaches the remote only when it is published (`git fetch`, then `flai release --pending`, or the board's Publish), and agents publish only when the operator asks. The work-management convention already says so. Today `inbox` and `flai board` report `unpushed`, and the MCP server's instructions (`flai/internal/mcpserver/folder.go` ~195 and the single-project server's) tell agents to run `flai push --pending` before anything else (S-0063); `flai push --pending` and the `auto-publish` host action (ADR-0048) are a second way to the remote.

## Acceptance criteria
- [ ] The designer accepts ADR-0067, or the story records the changes they ask for in a new ADR, before anything else is built
- [ ] The MCP servers' instructions, `inbox`, `flai board`, and `harness.Prompt` no longer tell an agent to push an unpushed acceptance; the board and the dashboard still show what is accepted and not yet published
- [ ] `flai push --pending` and the `auto-publish` host action remains but is not integrated into the dashboard
- [ ] the push banner over the top of the board is removed entirely
- [ ] `flai release --pending` and Publish fetch first, or refuse with what to run when the clone is behind, as S-0174 does for tags
- [ ] `design/system/pushing-from-the-board.md`, `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `design/system/workflow.md`, and the operator and user guides describe publishing as the one way to the remote
- [ ] Tests cover the instructions and inbox without the push duty, and publishing after a fetch

## Tasks

## Notes

ADR-0067 is proposed because its details are the drafting agent's reading of the designer's instruction "update publishing and releasing".
