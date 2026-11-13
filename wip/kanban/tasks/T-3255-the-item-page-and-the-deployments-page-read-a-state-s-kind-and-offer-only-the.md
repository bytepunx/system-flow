---
id: T-3255
type: task
nature: feature
title: The item page and the Deployments page read a state's kind and offer only the moves the operator may make
status: backlog
parent: S-0758
owner: alex
created: 2026-11-12T16:09:57Z
updated: 2026-11-12T16:09:57Z
transitions: []
stream: S-0758
tags: [dashboard]
touches: ["flaiover/src/routes/items/[id]/+page.svelte", flaiover/src/routes/ops/deployments/+page.svelte]
after: [T-3252]
---
# T-3255 The item page and the Deployments page read a state's kind and offer only the moves the operator may make

## Work

Waits for T-3252: both pages read the module that task adds. Replace the state names in the item page with the module's kind and moves: offer only the moves its transitions allow the operator, and ask for a reason where one is needed. Make the Deployments page tell a story's namespace by the kind of its state.

## Done when

Neither page compares a status with a state's name, and the flaiover tests pass.

## Notes
