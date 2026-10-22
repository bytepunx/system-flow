---
id: T-2045
type: task
nature: improvement
title: The Components page sends the chosen example's answers to the harness, lets the operator edit them, and saves them with the example
status: backlog
parent: S-0476
owner: alex
created: 2026-10-22T21:00:10Z
updated: 2026-10-22T21:00:10Z
transitions: []
stream: S-0476
tags: [flaiover]
touches: [flaiover/src/lib/explorer/protocol.ts, flaiover/src/lib/explorer/protocol.test.ts, flaiover/src/lib/components/explorer/ComponentFrame.svelte, flaiover/src/lib/components/explorer/ComponentFrame.svelte.test.ts, flaiover/src/lib/components/explorer/ExampleList.svelte, flaiover/src/lib/components/explorer/ExampleList.svelte.test.ts, flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts]
after: [T-2040, T-2041]
---
# T-2045 The Components page sends the chosen example's answers to the harness, lets the operator edit them, and saves them with the example

## Work

- `flaiover/src/lib/explorer/protocol.ts`: `select` and `props` carry an optional `answers` array, typed as `ComponentAnswer[]` from `repo.ts`; `connect()`'s `select` and `props` senders take it. `protocol.test.ts` checks the field is sent whole and left out when there are none.
- `flaiover/src/lib/components/explorer/ComponentFrame.svelte`: takes an `answers` prop beside `props`, `functions`, and `snippets`, and sends it with every `select` and every `props` message, since a props message is whole. `ComponentFrame.svelte.test.ts` checks both messages carry it.
- `flaiover/src/routes/design/components/+page.svelte`: keep the chosen example's answers in page state. `chooseExample()` loads them with the props, theme, and width; `pick()` clears them with the props; `saveExample()` sends them in the PUT, and the saved example carries them. Below the examples list, an "Answers" panel shows them as a JSON text area with a short lead ("What the harness answers this component's requests to the app's server; unmatched requests get 501"); a parse error or a shape that is not an array of objects with a string `path` is shown beside it and keeps the last good value from reaching the frame. The event log already shows the harness's `fetch` events; show the third argument, the status, as the log shows any argument.
- `flaiover/src/lib/components/explorer/ExampleList.svelte`: `exampleFrom()` carries `answers` into what is saved, and the list marks an example that has answers (a small count, "2 answers") so the operator can tell which examples show the answered state.
- Tests in `components.svelte.test.ts` and `ExampleList.svelte.test.ts`: choosing an example with answers sends them in `select`; editing the text area sends them in `props`; saving keeps them; a malformed edit is shown and not sent.

This task waits for T-2040 (the harness's side of the protocol) and T-2041 (the route and the type), since it sends what the one reads and saves through the other.

## Done when

- [ ] Choosing an example with answers sends them with `select`, and each later `props` message carries them
- [ ] The operator can edit an example's answers on the page and save them with the example
- [ ] A malformed edit is shown as an error and does not reach the harness
- [ ] The tests named above pass, and vitest, prettier, eslint, and svelte-check are clean on the changed files

## Notes

Planned by planner-S-0476. A JSON text area is the smallest editor that fits; a form per answer can follow if the operator asks for it.
