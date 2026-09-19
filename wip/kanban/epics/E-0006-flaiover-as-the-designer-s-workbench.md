---
id: E-0006
type: epic
nature: feature
title: flaiover as the designer's workbench
status: ready
owner: alex
created: 2026-09-17T19:46:28Z
updated: 2026-09-19T02:09:43Z
transitions:
  - to: ready
    at: 2026-09-18T16:25:54Z
    by: alex
tags: []
---

# E-0006 flaiover as the designer's workbench

## Outcome
flaiover is where the system's designer reads, edits, discusses, and accepts the documentation and work, and therefore where they direct the agents implementing it. Files under `design/`, `docs/`, and `wip/` stay the single record; `flai mcp` gives agents a low-latency view of the designer's input; the dashboard is authenticated, edits through flai, shows who is working on what, and is shaped so a multi-project hub reached over a tunnel can be added later without redesign.

## Stories
- S-0036 Dashboard authentication with a per-project token
- S-0037 Story branches with wip on main, stream sync, and the touches flag
- S-0038 Comment threads and answers as files under wip
- S-0039 flai mcp: inbox, reply, read, and move tools with the template .mcp.json
- S-0040 Document editing in the dashboard with commit through flai
- S-0041 Review and acceptance from the dashboard
- S-0042 Agent presence, activity, and the designer's inbox
- S-0043 MCP over HTTP and project identity in the API, hub-ready
- S-0046 Moving a story to done is acceptance
- S-0048 Parent epic indicator on story cards on the board
- S-0050 Accepting a story from the dashboard works when the story has a worktree
- S-0051 Accepting from the dashboard shows uncommitted changes outside wip and lets the designer include them
- S-0052 Work out whether and how a board move to done can push the acceptance
- S-0054 Board cards: larger detail text under a thin divider
- S-0055 Board cards are colour coded by type and nature in pastels
- S-0057 Cards can be dragged within a column to change the pull order
- S-0058 An agent learns through MCP that the designer moved work to ready

## Notes
Carved out 2026-09-17 from the operator's direction. Decisions: ADR-0018 (dashboard token), ADR-0019 (story branches, wip on main, touches), ADR-0020 (files plus MCP). Pull order is the story order above: authentication first, hub transport last. E-0003 stays open for dashboard remediation only.
