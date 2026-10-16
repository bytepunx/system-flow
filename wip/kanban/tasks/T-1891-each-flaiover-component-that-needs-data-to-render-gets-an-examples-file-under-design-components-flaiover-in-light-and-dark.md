---
id: T-1891
type: task
nature: feature
title: Each flaiover component that needs data to render gets an examples file under design/components/flaiover, in light and dark
status: backlog
parent: S-0444
owner: alex
created: 2026-10-16T21:40:04Z
updated: 2026-10-16T21:49:04Z
transitions: []
stream: S-0444
tags: [dashboard]
touches: [design/components/flaiover]
after: [T-1890]
---
# T-1891 Each flaiover component that needs data to render gets an examples file under design/components/flaiover, in light and dark

## Work

For each listed component with a required prop, or whose defaults render nothing useful, write its examples file under `design/components/flaiover/`, at the path S-0443's ADR gives. Build each on the Components page and save it with "Save as example" on this story's branch, or write it in the worktree. Give each at least one example that renders the component's main state, taking its data from the shapes in `flaiover/src/lib` (an item, a thread, a chart option).

Waits for T-1890, whose exclusions decide which components need examples.

## Done when

Each listed component renders without an `error` on the Components page with the story picked as its source, and `flai check --strict` reports no `components.examples`.

## Notes

`design/components/flaiover` is a folder touch: which components need examples, and so which files are added, is known only once T-1890 has opened each. No other open story touches the folder, so it holds nothing.
