---
id: S-0274
type: story
nature: improvement
title: "Opening a story is one call: flai story start moves it to in-progress, opens the stream, primes, and answers the first inbox together"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:32Z
updated: 2026-10-05T01:35:32Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0274 Opening a story is one call: flai story start moves it to in-progress, opens the stream, primes, and answers the first inbox together

## Goal

A story agent's first minute is four turns: `flai move S-nnnn in-progress`, `flai stream open`, `prime`, and `inbox`, and in S-0248 the prime's result overflowed the harness and was read back from a file (S-0261 fixes the size). `flai story start S-nnnn`, `story_start` over MCP, and `story.start` on the host channel do the four in one call and answer the worktree path, the branch, the prime pack within its budget, and the inbox. flai serve's prompt tells the agent to begin with it.

## Acceptance criteria
- [ ] `flai story start S-nnnn` moves the story to in-progress, opens the stream, and answers the worktree, the branch, the prime pack, and the inbox in one result, as text and `--json`, refusing a story that is not ready or is held
- [ ] The same is `story_start` over MCP and `story.start` on the host channel
- [ ] The harness prompt, `design/conventions/session-start.md`, the template's copies, `design/system/flai-cli.md`, and the user guide begin the loop with it

## Tasks

## Notes

Depends on S-0261 for the prime pack to fit the harness's tool result limit when returned inline.
