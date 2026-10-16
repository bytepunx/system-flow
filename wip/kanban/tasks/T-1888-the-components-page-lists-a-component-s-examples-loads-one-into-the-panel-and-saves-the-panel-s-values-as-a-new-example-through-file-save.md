---
id: T-1888
type: task
nature: feature
title: The Components page lists a component's examples, loads one into the panel, and saves the panel's values as a new example through file.save
status: backlog
parent: S-0443
owner: alex
created: 2026-10-16T21:38:38Z
updated: 2026-10-16T21:48:35Z
transitions: []
stream: S-0443
tags: [dashboard]
touches: [flaiover/src/lib/components/explorer/ExampleList.svelte, flaiover/src/lib/components/explorer/ExampleList.svelte.test.ts, flaiover/src/routes/api/components/examples/+server.ts, flaiover/src/routes/api/components/examples/examples.test.ts, flaiover/src/routes/design/components/+page.svelte]
after: [T-1886, T-1887]
---
# T-1888 The Components page lists a component's examples, loads one into the panel, and saves the panel's values as a new example through file.save

## Work

- `ExampleList.svelte`: the picked component's examples. Choosing one loads its props, theme, and width into the page. "Save as example" asks for a name.
- `PUT /api/components/examples` `{component, example, hash, story?}`: reads the component's examples file, at the path `components.list` gives, with `repo().file`, adds or replaces the example, and saves with `file.save`, with no hash when the file is new. It answers 409 with the current content on a conflict and 422 with the findings on a refusal, as `/api/design/tokens` does.
- `+page.svelte`: the list beside the panel; on a refusal, the findings shown and the values kept.

Waits for T-1886 for the create and T-1887 for the examples and their paths in the catalogue.

## Done when

Tests cover choosing an example, a save that creates the file, one that adds to it, a conflict, and a refusal. `flai test flaiover/src/lib/components/explorer flaiover/src/routes/api/components flaiover/src/routes/design/components` passes.

## Notes

None.
