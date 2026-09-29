---
id: S-0165
type: story
nature: remediation
title: Tokens rate chart should express token use in minutes
status: backlog
parent: E-0013
owner: alex
created: 2026-09-29T22:37:02Z
updated: 2026-09-29T22:37:02Z
transitions: []
tags: [dashboard, cli]
topics: [client-side, server-side]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0165 Tokens rate chart should express token use in minutes

## Goal

The token rate chart is showing tokens / agent hour for stories but this is a bad metric - most stories are completed in minutes, not hours, the chosen unit of time measure makes token expenditure seem much higher. change this to agent minutes, not hours.

## Acceptance criteria
- [ ] token rate chart uses agent minutes, not agent hours

## Tasks

## Notes
