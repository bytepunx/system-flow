---
id: S-0556
type: story
nature: improvement
title: flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree
status: backlog
owner: alex
created: 2026-11-01T10:44:38Z
updated: 2026-11-01T10:44:38Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0556 flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree

## Goal

This story remediates [I-0240](../../../design/issues/I-0240-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md), "flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0240 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0240 is closed with `flai issue close I-0240 --reason` saying what fixed it

## Tasks

## Notes

I-0240's cost "0m" is not a duration longer than zero, so no time_lost_per_cycle was derived from it.
