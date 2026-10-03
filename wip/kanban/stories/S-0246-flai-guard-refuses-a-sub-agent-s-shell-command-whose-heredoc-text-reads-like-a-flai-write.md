---
id: S-0246
type: story
nature: improvement
title: flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write
status: backlog
owner: alex
created: 2026-10-03T17:49:40Z
updated: 2026-10-03T17:49:40Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0246 flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write

## Goal

This story remediates [I-0058](../../../design/issues/I-0058-flai-guard-refuses-a-sub-agent-s-shell-command-whose-heredoc-text-reads-like-a-flai-write.md), "flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0058 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0058 is closed with `flai issue close I-0058 --reason` saying what fixed it

## Tasks

## Notes
