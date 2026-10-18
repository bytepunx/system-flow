---
id: T-1633
type: task
nature: feature
title: An ADR records that under pull-request acceptance the merge goes through the remote, by flai serve when the token is not at hand, and that the acceptance commit reads as under local acceptance
status: done
parent: S-0397
owner: alex
created: 2026-10-16T04:39:54Z
updated: 2026-10-18T02:41:42Z
transitions:
  - to: ready
    at: 2026-10-18T02:28:29Z
    by: agent-S-0397
  - to: in-progress
    at: 2026-10-18T02:28:30Z
    by: agent-S-0397
  - to: done
    at: 2026-10-18T02:41:42Z
    by: agent-S-0397
stream: S-0397
tags: [cli]
touches: [design/adrs, flai/internal/forge/state_test.go]
usage:
  source: log
  seconds: 792
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 7854
      cache_read: 1114050
      cache_write: 43686
      cost: 0.626
---
# T-1633 An ADR records that under pull-request acceptance the merge goes through the remote, by flai serve when the token is not at hand, and that the acceptance commit reads as under local acceptance

## Work

- Write the ADR with `flai adr new`, refining ADR-0032, ADR-0067, ADR-0166, and ADR-0169. It records the design in the narrative's `## Decisions`:
  - the merge on the forge by `acceptance.merge_method`, at the head the gate passed at, through a new `PullRequest(number)` in the forge contract
  - the deletion of `origin`'s story branch with the host's git credentials, and the local branch removed with `git branch -D`
  - the fetch, and the rebase of unpublished wip-only commits onto `origin`'s main in a temporary worktree, with main moved by `git reset --keep`
  - the refusal of unpublished commits outside wip, before the merge
  - the hand-off to flai serve when the token's variable is unset: the request file under `.flai-cache/forge/merge/`, the merge pass, the 30s and 5m waits, and the forge host action, with no new listener or secret
  - the dashboard's acceptance merging in flai serve's process, refused while the forge host action is off
  - a pull request someone else merged: kept as merged, never reopened, and accepted with no second merge
  - a forge refusal that stops the acceptance before anything changes locally
  - an acceptance commit and archive that read the same in both modes
  - the acceptance commit still reaching the remote only when published (ADR-0067). A main branch that refuses direct pushes refuses that publish too, which this story leaves open, as TH-0464 says.
- **Layer.** First, beside T-1634: it waits for no task.

## Done when

The ADR is accepted, links what it refines, and `flai check --strict` passes on it.

## Notes
