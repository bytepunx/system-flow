---
id: E-0026
type: epic
nature: feature
title: Component Library
status: done
owner: alex
created: 2026-10-16T21:20:59Z
updated: 2026-10-17T06:40:03Z
transitions:
  - to: ready
    at: 2026-10-16T21:21:01Z
    by: alex
  - to: in-progress
    at: 2026-10-16T21:52:24Z
    by: agent-S-0436
  - to: review
    at: 2026-10-17T05:44:22Z
    by: agent-S-0444
  - to: in-progress
    at: 2026-10-17T06:39:19Z
    by: alex
  - to: review
    at: 2026-10-17T06:39:50Z
    by: alex
  - to: done
    at: 2026-10-17T06:40:03Z
    by: alex
tags: []
usage:
  source: sum
  seconds: 34805
  priced_by: harness
  turns:
    - day: 2026-10-16
      ceremony: 7
      hand_edits: 7
      work: 159
    - day: 2026-10-17
      ceremony: 16
      hand_edits: 23
      work: 440
  models:
    - model: claude-opus-5-5
      input: 4170
      output: 1722880
      cache_read: 231635908
      cache_write: 8289222
      cost: 128.1515
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
      seconds: 34807
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1206
          output: 21915
          cache_read: 249667954
          cache_write: 416570
          cost: 62.6961
        - model: claude-sonnet-5-5
          input: 148
          output: 766
          cache_read: 2754575
          cache_write: 559936
          cost: 2.7127
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
- 2026-10-17T05:44:22Z: moved to review: follows S-0444, which moved to review
- 2026-10-17T06:39:19Z: moved to in-progress: follows S-0444, which moved to in-progress
- 2026-10-17T06:39:50Z: moved to review: follows S-0444, which moved to review
- 2026-10-17T06:40:03Z: moved to done: follows S-0444, which moved to done
