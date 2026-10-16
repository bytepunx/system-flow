---
id: T-1633
type: task
nature: feature
title: An ADR records that under pull-request acceptance the merge goes through the remote, by flai serve when the token is not at hand, and that the acceptance commit reads as under local acceptance
status: backlog
parent: S-0397
owner: alex
created: 2026-10-16T04:39:54Z
updated: 2026-10-16T04:39:54Z
transitions: []
stream: S-0397
tags: [cli]
touches: [design/adrs]
---
# T-1633 An ADR records that under pull-request acceptance the merge goes through the remote, by flai serve when the token is not at hand, and that the acceptance commit reads as under local acceptance

## Work

- Write the ADR with `flai adr new`, refining ADR-0032, ADR-0067, and S-0395's ADR. It records:
  - the merge on the forge by `acceptance.merge_method`
  - the fetch and the rebase of unpublished `wip/` commits onto the remote's main
  - the hand-off to flai serve when the token's variable is unset
  - a request someone else merged
  - the refusals
  - that the acceptance commit and the archive are the same in both modes
  - that the acceptance commit still reaches the remote only when published, and what that means for a main branch that refuses direct pushes, as the plan's thread settles it
- First in its layer: it waits for no task.

## Done when

The ADR is accepted, and `flai check --strict` passes on it.

## Notes
