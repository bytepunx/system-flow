---
id: E-0020
type: epic
nature: feature
title: A Design System Extension
status: in-progress
owner: alex
created: 2026-10-13T17:37:49Z
updated: 2026-10-13T17:53:28Z
transitions:
  - to: ready
    at: 2026-10-13T17:53:19Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-13T17:53:28Z
    by: agent-S-0380
tags: [dashboard, cli]
topics: [design]
usage:
  source: sum
  seconds: 2466
  priced_by: harness
  turns:
    - day: 2026-10-13
      ceremony: 2
      hand_edits: 2
      work: 41
    - day: 2026-10-14
      ceremony: 3
      hand_edits: 3
      work: 33
  models:
    - model: claude-haiku-5-5
      input: 58701
      output: 3160
      cache_read: 0
      cache_write: 0
      cost: 0.0175
    - model: claude-opus-5-5
      input: 516
      output: 209043
      cache_read: 31935964
      cache_write: 1311811
      cost: 19.2778
  strategic:
    - kind: planner
      seconds: 898
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 74
          output: 33814
          cache_read: 2503250
          cache_write: 118029
          cost: 0.1415
        - model: claude-opus-5-5
          input: 126
          output: 71152
          cache_read: 12098341
          cache_write: 253452
          cost: 5.8708
    - kind: orchestrator
      seconds: 1833
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 208
          output: 3419
          cache_read: 25510616
          cache_write: 73598
          cost: 6.4438
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
- 2026-10-13T17:53:19Z: moved to ready: follows S-0380, which moved to ready
- 2026-10-13T17:53:28Z: moved to in-progress: follows S-0380, which moved to in-progress
