---
id: T-1662
type: task
nature: feature
title: A batch lands while the main checkout is dirty, stops only the story whose change would overwrite a dirty file, and records what the checkout is behind
status: backlog
parent: S-0402
owner: alex
created: 2026-10-16T05:02:19Z
updated: 2026-10-16T05:02:19Z
transitions: []
stream: S-0402
tags: [cli]
touches: [flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go]
after: [T-1647, T-1648]
---
# T-1662 A batch lands while the main checkout is dirty, stops only the story whose change would overwrite a dirty file, and records what the checkout is behind

## Work

- In `land.go`, land as T-1647's ADR decides, rather than refusing on any uncommitted file in the main checkout.
- Before each story lands, ask T-1648's function which dirty files its change would overwrite. With any, record the files as the story's outcome, skip it, and land the rest.
- Record in the state the commits and paths the checkout is behind, when the landing leaves its files behind main.
- Test landing with a dirty analysis report and issue files, as on TH-0462, and with a dirty file a story changes.
- Waits for T-1647, whose decision it implements, and T-1648, whose function it calls.

## Done when

Those landings behave as above, and the tests pass.

## Notes
