---
id: T-1862
type: task
nature: feature
title: An integration test runs the dev server on a fixture sub-project and shows a source edit and a stylesheet rewrite reach the harness as hot updates, not reloads
status: backlog
parent: S-0438
owner: alex
created: 2026-10-16T21:34:03Z
updated: 2026-10-16T21:34:03Z
transitions: []
stream: S-0438
tags: [cli]
touches: [flai/internal/components/integration_test.go, flai/internal/components/testdata]
after: [T-1861]
---
# T-1862 An integration test runs the dev server on a fixture sub-project and shows a source edit and a stylesheet rewrite reach the harness as hot updates, not reloads

## Work

- `integration_test.go`, skipped under `-short` and when `node` is absent: build a fixture SvelteKit sub-project in a temporary folder from `testdata/`, linking flaiover's installed `node_modules`. Start the dev server as `flai components serve` does, and fetch the harness page and one component's module. Connect to Vite's HMR WebSocket, edit the component, and rewrite the stylesheet the setup imports. Assert an `update` payload for each, and no `full-reload`.
- Run `flai components serve` by hand on this repository, open the harness for `flaiover/src/lib/components/design/Gallery.svelte`, and record in the story's notes what rendered.

Waits for T-1861, whose command path the test drives.

## Done when

`flai test --all flai/internal/components` runs the test and it passes on this host, and the story's notes record the manual run.

## Notes

`testdata` is a folder touch: the fixture's files are the agent's to name. A headless browser is not used: Vite's HMR payloads show that the update is hot, and the frame's side is S-0439's.
