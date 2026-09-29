---
id: S-0144
type: story
nature: remediation
title: Automatic publish/push should be a setting
status: backlog
owner: alex
created: 2026-09-29T01:04:15Z
updated: 2026-09-29T01:04:15Z
transitions: []
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0144 Automatic publish/push should be a setting

## Goal

At some point moving stories from review to done started automatically publishing without user involvement, breaking the operator's ability to batch stories for release. This needs to be an explicit setting that the operator enables, not something it defaults to.

## Acceptance criteria
- [ ] By default, enabling push in the server settings does _not_ mean that stories are automatically published/pushed
- [ ] A new setting that enables automatic publishing should be added so that controlling that behavior is explicit

## Tasks

## Notes
