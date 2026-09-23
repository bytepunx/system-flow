---
id: S-0103
type: story
nature: feature
title: Extend existing data structures, APIs, and UI to track agent configuration
status: done
parent: E-0008
owner: alex
created: 2026-09-23T16:46:55Z
updated: 2026-09-23T20:18:12Z
transitions:
  - to: ready
    at: 2026-09-23T16:52:02Z
    by: alex
  - to: in-progress
    at: 2026-09-23T16:56:54Z
    by: system-flow
  - to: review
    at: 2026-09-23T17:22:24Z
    by: system-flow
  - to: done
    at: 2026-09-23T20:18:12Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0103 Extend existing data structures, APIs, and UI to track agent configuration

## Goal

system flow should extend the data it currently tracks for projects and stories so that flai can be extended to manage one or more agent sessions based on kanban activity.

## Acceptance criteria

- [x] system flow captures the default harness and model as well as optional configuration to use for all stories
- [x] stories carry front-matter set to the default harness and model as well as optional configuration flai will use to invoke the agent to work on that story
- [x] the operator is able to specify overrides for the harness, model, and configuration that should work the story
- [x] the story creation and editing CLI, MCP, and web interface all need to support these new fields/values

## Tasks
- T-0362 Projects and stories carry an agent (harness, model, configuration); flai agent sets the project's defaults, and new stories get them
- T-0363 flai story new, flai edit, flai show, and hostapi carry the agent fields
- T-0364 MCP: item_get shows the agent; item_new and item_edit create and edit stories with it
- T-0365 The dashboard's new-story form, story editor, and story page show and set the agent
- T-0366 Tried end to end; docs

## Notes
