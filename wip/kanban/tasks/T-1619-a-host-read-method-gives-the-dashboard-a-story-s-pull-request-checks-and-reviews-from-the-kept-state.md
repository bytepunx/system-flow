---
id: T-1619
type: task
nature: feature
title: A host read method gives the dashboard a story's pull request, checks, and reviews from the kept state
status: backlog
parent: S-0394
owner: alex
created: 2026-10-16T04:36:58Z
updated: 2026-10-16T04:36:58Z
transitions: []
stream: S-0394
tags: [cli, dashboard]
touches: [flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go]
after: [T-1615]
---
# T-1619 A host read method gives the dashboard a story's pull request, checks, and reviews from the kept state

## Work

- In `flai/internal/hostapi/reads.go`, add the read method `forge.status`, taking a story ID. It answers with T-1615's state for the story: the pull request's URL, its checks, and its reviews, or none.
- It needs no host action: it reads only what flai serve kept.
- Add the method to the contract test if the package has one, `contract_test.go`.
- Waits for T-1615: it reads the kept state.

## Done when

`forge.status` answers for a story with state and one without, and the hostapi tests pass.

## Notes
