---
id: T-1209
type: task
nature: improvement
title: The dashboard's story agent says an agent waits on another story's agent, naming the story
status: backlog
parent: S-0335
owner: alex
created: 2026-10-07T20:16:41Z
updated: 2026-10-07T20:16:41Z
transitions: []
stream: S-0335
tags: [flaiover]
touches: [flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts]
after: [T-1208]
---
# T-1209 The dashboard's story agent says an agent waits on another story's agent, naming the story

## Work

Show the new waiting reason. It waits for T-1208, which adds the field `agent.status` carries.

- When the agent waits on a story's agent, `StoryAgent.svelte` says so and links the story it waits on, instead of saying it waits on the operator.
- A wait on the operator reads as before.

## Done when

- A component test covers both reasons.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
