---
id: S-0546
type: story
nature: improvement
title: The close-out runs as a flai serve job the story's agent waits on by event, so it never dies at the harness's tool timeout and is never re-run from the start
status: backlog
parent: E-0031
owner: alex
created: 2026-10-30T01:51:56Z
updated: 2026-10-30T01:51:56Z
transitions: []
tags: [flai]
topics: [testing, mcp, cli]
touches: [flai/internal/serve, flai/internal/mcpserver, flai/cmd, scripts/close-out.sh, design/system/flai-cli.md, design/conventions/testing.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0546 The close-out runs as a flai serve job the story's agent waits on by event, so it never dies at the harness's tool timeout and is never re-run from the start

## Goal

`scripts/close-out.sh` runs `flai verify` inside one Bash tool call of the agent's turn. Claude Code ends a tool call at 600 s, and 64 of the fortnight's 267 close-outs died there, 146 ended with no summary line, and only 88 reported every step passing; the agent then ran the whole close-out again from the start, every tier included, p50 109 s and p85 the full 600 s, for 18.5 hours of the fortnight's 113 session-hours. This story makes the close-out a job: the agent asks flai serve, through the `verify` MCP tool or `flai verify --job`, to run the close-out in the story's worktree, gets a job id at once, and waits for its `verified` event with `wait_for_events`, so no tool call outlives its timeout and a close-out is never cut short. The job's tiers, output, and verdict are recorded where `flai verify` records them today, and a close-out that is asked for again at the same commit answers from that record instead of running.

## Acceptance criteria
- [ ] `verify` (MCP) and `flai verify --job` start the close-out as a serve job in the story's worktree and return its id without waiting; `wait_for_events` reports `verified` with the job's verdict when it ends, and `flai verify` without `--job` behaves as today for a session without serve
- [ ] A close-out asked for again at a commit the record already holds a verdict for answers from the record, saying so, and runs nothing
- [ ] `scripts/close-out.sh`, the testing convention, `flai-cli.md`, and the user docs say the close-out is a job the agent waits on, and the story agent's prompt tells it to wait rather than call the script in a turn
- [ ] The serve journal records each job's start, end, duration, and verdict, and `flai stats` counts close-outs from it

## Tasks

The planner drafts them.

## Notes

Evidence: the fortnight's story agent logs, counted in the operator's conversation on 2026-10-29. The harness timeout is the harness's, not a setting flai can raise, so the fix is to leave the turn.
