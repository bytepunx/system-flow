---
id: S-0711
type: story
nature: improvement
title: sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt
status: backlog
owner: alex
created: 2026-11-10T09:56:05Z
updated: 2026-11-11T09:03:41Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-11-11T09:03:41Z
---
# S-0711 sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt

## Goal

This story remediates [I-0340](../../../design/issues/I-0340-sendbackreason-puts-a-ci-check-name-and-reviewer-login-unquoted-into-the-agent.md), "sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0340 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0340 is closed with `flai issue close I-0340 --reason` saying what fixed it

## Tasks

## Notes
