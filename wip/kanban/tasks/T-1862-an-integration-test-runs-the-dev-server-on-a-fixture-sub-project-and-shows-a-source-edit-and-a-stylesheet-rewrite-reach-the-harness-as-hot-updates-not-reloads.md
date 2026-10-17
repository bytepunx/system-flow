---
id: T-1862
type: task
nature: feature
title: An integration test runs the dev server on a fixture sub-project and shows a source edit and a stylesheet rewrite reach the harness as hot updates, not reloads
status: done
parent: S-0438
owner: alex
created: 2026-10-16T21:34:03Z
updated: 2026-10-17T00:39:38Z
transitions:
  - to: ready
    at: 2026-10-16T23:52:10Z
    by: system-flow
  - to: in-progress
    at: 2026-10-16T23:52:11Z
    by: system-flow
  - to: done
    at: 2026-10-17T00:39:38Z
    by: system-flow
stream: S-0438
tags: [cli]
touches: [flai/internal/components/integration_test.go, flai/internal/components/testdata, docs/users/flai-reference.md]
after: [T-1861]
usage:
  source: log
  seconds: 2845
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 151
      output: 600
      cache_read: 8152880
      cache_write: 423262
      cost: 3.9869
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
