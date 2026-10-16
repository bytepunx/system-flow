---
id: T-1623
type: task
nature: feature
title: The acceptance preview lists each required check not passed, each standing requested change, missing approvals, a missing pull request, and a stale state as blockers under pull-request acceptance
status: backlog
parent: S-0395
owner: alex
created: 2026-10-16T04:38:05Z
updated: 2026-10-16T04:38:05Z
transitions: []
stream: S-0395
tags: [cli]
touches: [flai/internal/preview/accept.go, flai/internal/preview/accept_forge_test.go]
---
# T-1623 The acceptance preview lists each required check not passed, each standing requested change, missing approvals, a missing pull request, and a stale state as blockers under pull-request acceptance

## Work

- In `flai/internal/preview/accept.go`, when the manifest's `acceptance.mode` is `pull-request`, read S-0394's kept state for the story. Add one blocker, by name, for each of these:
  - every required check failing or pending
  - every standing review requesting changes
  - approvals fewer than the base branch requires
  - no pull request
  - a state recorded at a commit other than the branch's head
- Give each blocker a stable `kind`, so that `--json` and the dashboard can tell them apart.
- Under `local`, add nothing.
- Test each blocker, a story that passes, and a project in `local` mode, in `accept_forge_test.go`.
- First in its layer: it waits for no task.

## Done when

The preview names each forge blocker in the cases above, adds none under `local`, and the tests pass.

## Notes
