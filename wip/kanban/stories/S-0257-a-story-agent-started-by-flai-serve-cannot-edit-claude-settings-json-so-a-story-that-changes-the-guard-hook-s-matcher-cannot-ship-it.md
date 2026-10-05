---
id: S-0257
type: story
nature: remediation
title: A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it
status: in-progress
owner: alex
created: 2026-10-04T03:59:16Z
updated: 2026-10-05T04:27:27Z
transitions:
  - to: ready
    at: 2026-10-04T21:43:07Z
    by: alex
  - to: in-progress
    at: 2026-10-05T04:16:31Z
    by: agent-S-0257
tags: [flai]
touches: [flai/internal/mcpserver, flai/internal/harness, flai/internal/hostapi, flai/cmd/mcp.go, flai/cmd/mcp_test.go, design/adrs, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/delegation.md, design/issues, flai/cmd/serve_actions_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1579
  estimated: true
  models:
    - model: claude-haiku-4-5-20251001
      input: 380
      output: 64
      cache_read: 2052012
      cache_write: 137451
      cost: 0.5043
    - model: claude-opus-5-5
      input: 332
      output: 1990
      cache_read: 13309706
      cache_write: 437274
      cost: 5.6758
---
# S-0257 A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it

## Goal

This story remediates [I-0069](../../../design/issues/I-0069-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md), "A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0069 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0069 is closed with `flai issue close I-0069 --reason` saying what fixed it

## Tasks
- T-0874 flai's MCP server has a permission_prompt tool that asks the operator on a thread before a .claude/ write
- T-0875 flai serve starts claude-code with --permission-prompt-tool mcp__flai__permission_prompt
- T-0876 A shell-only auto-approve host action lets permission_prompt allow .claude/ writes without asking
- T-0877 The ADR, design, and users' docs describe permission_prompt and auto-approve, and I-0069 is closed

## Notes
