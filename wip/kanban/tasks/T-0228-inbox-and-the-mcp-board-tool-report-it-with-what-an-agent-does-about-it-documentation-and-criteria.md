---
id: T-0228
type: task
nature: feature
title: "inbox and the MCP board tool report it, with what an agent does about it; documentation and criteria"
status: in-progress
parent: S-0063
owner: alex
created: 2026-09-19T08:34:59Z
updated: 2026-09-19T08:42:21Z
transitions:
  - to: ready
    at: 2026-09-19T08:42:20Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T08:42:21Z
    by: system-flow
stream: S-0063
tags: []
---

# T-0228 inbox and the MCP board tool report it, with what an agent does about it; documentation and criteria

## Work
`inbox` reports `unpushed` on every call while it exists, like ready work, and the `board` tool carries it; the server gets a command runner for this and reports nothing when it has none. The server instructions and `design/conventions/work-management.md` (template baseline first, then here) say what an agent does: fetch, push the branch and those tags with `flai push --pending`, never force. This part edits the same struct and the same instructions string as S-0061, which is in review: do it last, after S-0061 is accepted and the branch is synced, so the board can still accept both. User docs, operators' docs, `flai-cli.md`, `flaiover-dashboard.md`; criteria.

## Done when
- `inbox` and `board` carry `unpushed`, tested
- The convention and the documents say what to do about it
- The criteria are ticked and `flai check --strict` is clean

## Notes
