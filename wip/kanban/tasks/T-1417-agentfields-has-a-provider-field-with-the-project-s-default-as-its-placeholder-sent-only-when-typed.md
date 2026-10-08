---
id: T-1417
type: task
nature: feature
title: AgentFields has a provider field, with the project's default as its placeholder, sent only when typed
status: backlog
parent: S-0359
owner: alex
created: 2026-10-08T08:53:20Z
updated: 2026-10-08T08:53:31Z
transitions: []
stream: S-0359
tags: [dashboard]
touches: [flaiover/src/lib/components/AgentFields.svelte, flaiover/src/lib/agent.ts, flaiover/src/lib/agent.test.ts, flaiover/src/lib/components/ItemEditor.svelte, flaiover/src/lib/components/NewItemForm.svelte]
after: [T-1415]
---
# T-1417 AgentFields has a provider field, with the project's default as its placeholder, sent only when typed

## Work

- `agent.ts`: `Agent` gains `provider`; the builder, the equality check, and the one-line summary carry it.
- `AgentFields.svelte`: a bindable provider field beside harness and model, with the project's default as its placeholder.
- `ItemEditor.svelte` and `NewItemForm.svelte` bind it and send it with the agent; the Settings page's default agent binds it in the Providers section task, which owns `SettingsPanel.svelte`.
- `agent.test.ts`: built with and without a provider, equality, and the summary.

It waits for T-1415, whose writes accept `provider`. It runs together with T-1416, which touches other files.

## Done when

- `flai test flaiover/src/lib/agent.ts flaiover/src/lib/components/` passes.

## Notes

Layer 2 of S-0359.
