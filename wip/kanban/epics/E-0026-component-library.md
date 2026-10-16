---
id: E-0026
type: epic
nature: feature
title: Component Library
status: backlog
owner: alex
created: 2026-10-16T21:20:59Z
updated: 2026-10-16T21:20:59Z
transitions: []
tags: []
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

## Notes
