---
id: S-0103
type: story
nature: feature
title: Extend existing data structures, APIs, and UI to track agent configuration
status: backlog
parent: E-0008
owner: alex
created: 2026-09-23T16:46:55Z
updated: 2026-09-23T16:46:55Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0103 Extend existing data structures, APIs, and UI to track agent configuration

## Goal

system flow should extend the data it currently tracks for projects and stories so that flai can be extended to manage one or more agent sessions based on kanban activity.

## Acceptance criteria

- [ ] system flow captures the default harness and model as well as optional configuration to use for all stories
- [ ] stories carry front-matter set to the default harness and model as well as optional configuration flai will use to invoke the agent to work on that story
- [ ] the operator is able to specify overrides for the harness, model, and configuration that should work the story
- [ ] the story creation and editing CLI, MCP, and web interface all need to support these new fields/values

## Tasks

## Notes
