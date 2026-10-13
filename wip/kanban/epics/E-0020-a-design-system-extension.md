---
id: E-0020
type: epic
nature: feature
title: A Design System Extension
status: backlog
owner: alex
created: 2026-10-13T17:37:49Z
updated: 2026-10-13T17:45:42Z
transitions: []
tags: [dashboard, cli]
topics: [design]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 132
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 81
          cache_read: 1103288
          cache_write: 3859
          cost: 0.2791
cost_of_delay:
  inputs:
    time_lost_per_cycle: 8h
    by: alex
    at: 2026-10-13T17:37:49Z
---
# E-0020 A Design System Extension

## Outcome

At the moment, many projects I work on use Svelte and TailwindCSS. There is not currently any way to settle on choices or guidance in a way that makes the impact on the visual language of a solution tangible until work is delivered and deployed.

This results in an expensive way to iterate towards desirable UI/UX as the operator must manually test and use the system to find areas where there are inconsistencies, usability issues, etc. and then create epics or stories to address the shortfall.

This epic should specify the stories necessary to create a new top-level heading for building out a design system tool at the project level that the agents can leverage and refer to. This top level `Design` menu option can then host necessary pages and tooling to allow the user to work with agent(s) in developing and establishing their project's design in such a way that they can see it respond to their decisions in the browser without having to build an entire working system and deploy it before they can view and interact with it.

## Stories
- S-0380 Research how a project keeps its design system in the repository, how the dashboard previews it live, and how agents work on it with the operator
- S-0381 A project's design system is a document whose front matter holds its tokens, and flai validates it, renders it as CSS, and serves it to the dashboard
- S-0382 The dashboard's Design group shows a project's design system and previews its tokens live on sample components
- S-0383 The operator changes the design system's tokens on the Design page, sees the preview follow before saving, and saves through flai
- S-0384 An agent's proposal for the design system is a story whose branch changes it, and the Design page previews that branch beside the main branch's before the operator accepts it
- S-0385 flai writes the design system's tokens into the stylesheet the manifest names, and flai check reports the stylesheet when it falls behind them
- S-0386 A story's agent working on a project's UI is primed with its design system, and a baseline convention tells it to build from the tokens and propose changes on its branch

## Notes
