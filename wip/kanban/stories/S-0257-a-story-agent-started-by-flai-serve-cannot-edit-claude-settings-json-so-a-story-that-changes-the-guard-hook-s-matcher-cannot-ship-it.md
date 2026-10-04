---
id: S-0257
type: story
nature: remediation
title: A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it
status: ready
owner: alex
created: 2026-10-04T03:59:16Z
updated: 2026-10-04T21:43:07Z
transitions:
  - to: ready
    at: 2026-10-04T21:43:07Z
    by: alex
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0257 A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it

## Goal

This story remediates [I-0069](../../../design/issues/I-0069-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md), "A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0069 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0069 is closed with `flai issue close I-0069 --reason` saying what fixed it

## Tasks

## Notes
