---
id: T-1876
type: task
nature: feature
title: The harness lays a tokens message over the project's stylesheet in one style element, and an empty message removes it
status: backlog
parent: S-0441
owner: alex
created: 2026-10-16T21:36:27Z
updated: 2026-10-16T21:36:27Z
transitions: []
stream: S-0441
tags: [cli]
touches: [flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go]
---
# T-1876 The harness lays a tokens message over the project's stylesheet in one style element, and an empty message removes it

## Work

- `harness.mjs`: on `tokens`, write its CSS into one `<style data-flai-tokens>` appended last in `<head>`, so it wins over the project's stylesheet. Keep it last after a hot update adds styles, and remove it on an empty message.
- `integration_test.go`: the harness page serves the element's code. The cascade itself is checked by hand in T-1878.

Waits for nothing in this story.

## Done when

`flai test --all flai/internal/components` passes.

## Notes

None.
