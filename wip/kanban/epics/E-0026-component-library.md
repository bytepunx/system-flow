---
id: E-0026
type: epic
nature: feature
title: Component Library
status: in-progress
owner: alex
created: 2026-10-16T21:20:59Z
updated: 2026-10-16T21:52:24Z
transitions:
  - to: ready
    at: 2026-10-16T21:21:01Z
    by: alex
  - to: in-progress
    at: 2026-10-16T21:52:24Z
    by: agent-S-0436
tags: []
usage:
  source: sum
  seconds: 5717
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-16
      ceremony: 4
      hand_edits: 4
      work: 94
  models:
    - model: claude-opus-5-5
      input: 728
      output: 169088
      cache_read: 33516085
      cache_write: 1333505
      cost: 18.1512
  strategic:
    - kind: planner
      seconds: 1821
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 30
          output: 31642
          cache_read: 1192291
          cache_write: 134508
          cost: 0.1491
        - model: claude-opus-5-5
          input: 2298
          output: 167660
          cache_read: 572807022
          cache_write: 650437
          cost: 257.0809
    - kind: orchestrator
      seconds: 3431
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 392
          output: 7063
          cache_read: 50097177
          cache_write: 106748
          cost: 12.5868
cost_of_delay:
  inputs:
    penalty_per_week: 600
    by: alex
    at: 2026-10-16T21:20:59Z
---
# E-0026 Component Library

## Outcome

Tools like Storybook (React, Svelte, Vue) and Histoire (Svelte, Vue) allow UI/UX developers to develop, interact with, and test UI components in the browser. They do have limitations, however - Storybook isn't Vite native so HML (hot module loading) isn't available. Both Storybook and Histoire push the developer towards an IDE. flai/flaiover seek to give the operator in the loop a single-pane-of-glass experience.

We want to build an interactive component library/explorer that:

 1. gives the operator a way to see and interact with every component in their system
 2. allow the operator to change properties to manipulate the component and see how it renders as its underlying data changes
 3. see how changes to the design system or the underlying source impact the component live (without a redeploy or full rebuild of the project)

## Stories
- S-0436 flai reads and saves a text file that is not markdown through file.get and file.save, and a token save writes the stylesheet with it
- S-0437 flai lists every component of a project's Svelte sub-projects and the props each takes, read with the project's own Svelte compiler
- S-0438 flai runs a Svelte sub-project's own Vite dev server on the host, behind the components host action, with a harness that renders any one of its components live
- S-0439 The dashboard's Components page lists every component of the project and renders the one the operator picks, live, from the project's own dev server
- S-0440 The operator changes a component's props on the Components page, a control for each from its type, and sees it render as they change
- S-0441 A token change on the Design page reaches the components on the Components page, as the operator tries it and once it is saved
- S-0442 The Components page renders a story's worktree, so the operator sees a component an agent is changing live while the story is in progress or in review
- S-0443 A component's examples are kept under design/components, chosen on the Components page, and kept in step with its props by flai check
- S-0444 Every component of flaiover renders on the Components page, with the setup its stylesheet needs and examples for the components that need data

## Notes
- 2026-10-16T21:52:24Z: moved to in-progress: follows S-0436, which moved to in-progress
