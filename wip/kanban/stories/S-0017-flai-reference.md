---
id: S-0017
type: story
nature: feature
title: flai reference
status: in-progress
parent: E-0004
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-24T08:13:01Z
transitions:
  - to: ready
    at: 2026-09-24T07:58:25Z
    by: alex
  - to: in-progress
    at: 2026-09-24T08:12:05Z
    by: system-flow
tags: [cli]
touches: [flai/cmd, docs/users, design/system, scripts, Makefile]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0017 flai reference

## Goal
docs/users/flai.md documents every command, flag, and the config file.

## Acceptance criteria
- [ ] Generated from cobra help and reviewed

## Tasks
- T-0400 Generate the command reference from the cobra tree
- T-0401 Review the generated reference and fix the help text
- T-0402 Document the config file and link the reference

## Notes
