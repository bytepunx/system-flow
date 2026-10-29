---
id: T-2331
type: task
nature: feature
title: Show the story agent's stream in the Agent panel
status: done
parent: S-0541
owner: alex
created: 2026-10-29T14:14:40Z
updated: 2026-10-29T14:16:36Z
transitions:
  - to: ready
    at: 2026-10-29T14:14:59Z
    by: agent-S-0541
  - to: in-progress
    at: 2026-10-29T14:15:00Z
    by: agent-S-0541
  - to: done
    at: 2026-10-29T14:16:36Z
    by: agent-S-0541
stream: S-0541
tags: []
touches: [flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts]
usage:
  source: log
  seconds: 96
  estimated: true
  priced_by: harness
  models:
    - model: claude-sonnet-5-5
      input: 9
      output: 1335
      cache_read: 438245
      cache_write: 19316
      cost: 0.1345
---
# T-2331 Show the story agent's stream in the Agent panel

## Work

StoryAgent.svelte mounts AgentStream below its buttons once the story's agent has started, open while the agent runs, as ActivityView shows it on the Activity page. StoryAgent.svelte.test.ts covers it.

## Done when

- The Agent panel of a story with a started agent holds the stream, open while the agent works or waits.
- A story with no started agent shows no stream.

## Notes
