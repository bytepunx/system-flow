---
id: T-1611
type: task
nature: feature
title: An ADR records that flai serve pushes a story's branch and opens its pull request on the move to review, behind the forge host action, with the force-with-lease exception for story branches
status: done
parent: S-0393
owner: alex
created: 2026-10-16T04:35:29Z
updated: 2026-10-17T07:57:42Z
transitions:
  - to: ready
    at: 2026-10-17T07:44:47Z
    by: agent-S-0393
  - to: in-progress
    at: 2026-10-17T07:44:51Z
    by: agent-S-0393
  - to: done
    at: 2026-10-17T07:57:42Z
    by: agent-S-0393
stream: S-0393
tags: [cli]
touches: [design/adrs]
usage:
  source: log
  seconds: 771
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 4725
      cache_read: 608325
      cache_write: 17923
      cost: 0.3185
---
# T-1611 An ADR records that flai serve pushes a story's branch and opens its pull request on the move to review, behind the forge host action, with the force-with-lease exception for story branches

## Work

- Write the ADR with `flai adr new`, refining S-0392's ADR and ADR-0032. It records:
  - that flai serve pushes and opens or updates the pull request for a story in review under pull-request acceptance, and why there rather than in the move itself: the host's credentials, and an agent's move running in the agent's environment
  - that the pull request is found by its head branch
  - the `forge` host action
  - the token's absence from agents' environments
  - the story-branch exception to never force pushing, as TH-0463 settles it
- First in its layer: it waits for no task.

## Done when

The ADR is accepted, records the operator's answer on TH-0463 question 4, and `flai check --strict` passes on it.

## Notes
