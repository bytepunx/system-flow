---
id: T-1209
type: task
nature: improvement
title: The dashboard's story agent says an agent waits on another story's agent, naming the story
status: done
parent: S-0335
owner: alex
created: 2026-10-07T20:16:41Z
updated: 2026-10-07T21:59:15Z
transitions:
  - to: ready
    at: 2026-10-07T21:58:29Z
    by: agent-S-0335
  - to: in-progress
    at: 2026-10-07T21:58:29Z
    by: agent-S-0335
  - to: done
    at: 2026-10-07T21:59:15Z
    by: agent-S-0335
stream: S-0335
tags: [flaiover]
touches: [flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, flaiover/src/lib/activity.ts]
after: [T-1208]
usage:
  source: log
  seconds: 46
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 14
      output: 5778
      cache_read: 791971
      cache_write: 25577
      cost: 0.4313
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
