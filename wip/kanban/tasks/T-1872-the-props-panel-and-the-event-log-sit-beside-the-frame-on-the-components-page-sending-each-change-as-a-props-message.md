---
id: T-1872
type: task
nature: feature
title: The props panel and the event log sit beside the frame on the Components page, sending each change as a props message
status: backlog
parent: S-0440
owner: alex
created: 2026-10-16T21:35:39Z
updated: 2026-10-16T21:35:39Z
transitions: []
stream: S-0440
tags: [dashboard]
touches: [flaiover/src/lib/components/explorer/PropsPanel.svelte, flaiover/src/lib/components/explorer/PropsPanel.svelte.test.ts, flaiover/src/lib/components/explorer/EventLog.svelte, flaiover/src/routes/design/components/+page.svelte]
after: [T-1870]
---
# T-1872 The props panel and the event log sit beside the frame on the Components page, sending each change as a props message

## Work

- `PropsPanel.svelte`: a row per prop with its name, its type as written, its default, a required mark, and T-1870's control. It shows invalid JSON inline and does not send it, and has a reset.
- `EventLog.svelte`: the last 50 `event` messages with their time, prop, and arguments.
- `+page.svelte`: the panel and the log beside the frame. It sends `props` on each change and keeps them in the URL with T-1870's `toQuery`.

Waits for T-1870, whose model the panel renders.

## Done when

Component tests cover a control per kind, invalid JSON, reset, the URL restored on load, and events listed. `flai test flaiover/src/lib/components/explorer flaiover/src/routes/design/components` passes.

## Notes

None.
