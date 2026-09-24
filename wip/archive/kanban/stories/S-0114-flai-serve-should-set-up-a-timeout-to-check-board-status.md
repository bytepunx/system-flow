---
id: S-0114
type: story
nature: remediation
title: Flai serve should set up a timeout to check board status
status: done
parent: E-0008
owner: alex
created: 2026-09-24T08:16:14Z
updated: 2026-09-24T08:47:33Z
transitions:
  - to: ready
    at: 2026-09-24T08:16:35Z
    by: alex
  - to: in-progress
    at: 2026-09-24T08:28:03Z
    by: system-flow
  - to: review
    at: 2026-09-24T08:44:11Z
    by: agent-S-0114
  - to: done
    at: 2026-09-24T08:47:33Z
    by: alex
tags: [cli]
touches: [flai/cmd, flai/internal/serve, flai/internal/harness, design/adrs, design/system, docs/users, docs/operators]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
---
# S-0114 Flai serve should set up a timeout to check board status

## Goal

flai serve needs to perform periodic board state checks to determine if there are stories that are ready and can be moved into in process.

## Acceptance criteria
- [x] Items that exist in ready and don't get picked up immediately should get picked up when the timer elapses.

## Tasks
- T-0408 Someone attending holds a ready story back for attended minutes, then flai serve starts it
- T-0409 The claude-code harness names its agent to flai's MCP server itself, so an environment that sets FLAI_AGENT cannot merge every started agent into one name
- T-0410 Record the bounded hold in an ADR, the living design, and the operator and user docs

## Notes

- The launcher already looked every minute (S-0104); the serve log on 2026-09-24 shows those looks. What held ready stories back was the attended rule, with no bound, fed by a shared MCP cursor: `.claude/settings.local.json` sets `FLAI_AGENT=system-flow` in its `env`, which Claude Code applies to the MCP server it spawns, so every agent flai serve started was one name (I-0037) and none of their signs counted as flai's own.
- The "timer" delivered is a bound on that hold (ADR-0042): a story that only someone attending holds back is started after `attended_minutes`, and the reason says until when. The limit is reported ahead of attendance. `flai mcp --agent` and the claude-code adapter passing it fix the shared name.
- Verified by `make flai-test` (lint and all three tiers) on the story branch. The serving flai on this host is the installed 1.15.5, so the running launcher behaves as before until the story is published and `flai host upgrade` runs.
- Recommended to the operator: remove `FLAI_AGENT` from `.claude/settings.local.json`; flai serve sets the name itself, and shell commands in a started agent still see the settings' value.
