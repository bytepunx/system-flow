---
id: T-1417
type: task
nature: feature
title: AgentFields has a provider field, with the project's default as its placeholder, sent only when typed
status: done
parent: S-0359
owner: alex
created: 2026-10-08T08:53:20Z
updated: 2026-10-12T17:04:44Z
transitions:
  - to: ready
    at: 2026-10-12T16:56:01Z
    by: agent-S-0359
  - to: in-progress
    at: 2026-10-12T16:56:01Z
    by: agent-S-0359
  - to: done
    at: 2026-10-12T17:04:44Z
    by: agent-S-0359
stream: S-0359
tags: [dashboard]
touches: [flaiover/src/lib/components/AgentFields.svelte, flaiover/src/lib/agent.ts, flaiover/src/lib/agent.test.ts, flaiover/src/lib/components/ItemEditor.svelte, flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/components/ItemEditor.svelte.test.ts, flaiover/src/lib/components/NewItemForm.svelte.test.ts]
usage:
  source: log
  seconds: 523
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 23
      output: 7324
      cache_read: 1290126
      cache_write: 40247
      cost: 0.655
---
# T-1417 AgentFields has a provider field, with the project's default as its placeholder, sent only when typed

## Work

- `agent.ts`: `Agent` gains `provider`; the builder, the equality check, and the one-line summary carry it.
- `AgentFields.svelte`: a bindable provider field beside harness and model, with the project's default as its placeholder.
- `ItemEditor.svelte` and `NewItemForm.svelte` bind it and send it with the agent; the Settings page's default agent binds it in the Providers section task, which owns `SettingsPanel.svelte`.
- `agent.test.ts`: built with and without a provider, equality, and the summary.

It waits for no task: S-0349 already made flai's agent writes accept `provider`, so it runs together with T-1415, which touches only Go files.

## Done when

- `flai test flaiover/src/lib/agent.ts flaiover/src/lib/components/` passes.

## Notes

Layer 1 of S-0359.
