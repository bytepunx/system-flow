---
id: S-0257
type: story
nature: remediation
title: A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it
status: done
owner: alex
created: 2026-10-04T03:59:16Z
updated: 2026-10-05T04:40:46Z
transitions:
  - to: ready
    at: 2026-10-04T21:43:07Z
    by: alex
  - to: in-progress
    at: 2026-10-05T04:16:31Z
    by: agent-S-0257
  - to: review
    at: 2026-10-05T04:40:21Z
    by: agent-S-0257
  - to: done
    at: 2026-10-05T04:40:46Z
    by: alex
tags: [flai]
touches: [flai/internal/mcpserver, flai/internal/harness, flai/internal/hostapi, flai/cmd/mcp.go, flai/cmd/mcp_test.go, design/adrs, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/delegation.md, design/issues, flai/cmd/serve_actions_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2168
  models:
    - model: claude-haiku-4-5-20251001
      input: 388
      output: 14153
      cache_read: 2123691
      cache_write: 140631
      cost: 0.4593
    - model: claude-opus-5-5
      input: 392
      output: 116840
      cache_read: 18545279
      cache_write: 484657
      cost: 9.049
    - model: claude-sonnet-5
      input: 70
      output: 16978
      cache_read: 1916960
      cache_write: 121671
      cost: 0.8575
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
