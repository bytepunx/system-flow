---
id: T-0313
type: task
nature: feature
title: "flai accept merges and moves to done; it computes no bump, creates no tag, pushes nothing"
status: done
parent: S-0087
owner: alex
created: 2026-09-21T03:32:54Z
updated: 2026-09-21T03:46:24Z
transitions:
  - to: ready
    at: 2026-09-21T03:33:21Z
    by: system-flow
  - to: in-progress
    at: 2026-09-21T03:33:21Z
    by: system-flow
  - to: done
    at: 2026-09-21T03:46:24Z
    by: system-flow
stream: S-0087
tags: []
---
# T-0313 flai accept merges and moves to done; it computes no bump, creates no tag, pushes nothing

## Work
`acceptItem` (`flai/cmd/accept.go`) stops calling `a.planReleaseAt`, `release.Apply`, the version-file commit, and `release.Tag`/push once the branch is merged and the item moved to done. Research (cuts no release, ADR-0025) and experiment (refused) behave exactly as before, since they are the same merge/move path with no release step added or removed for them. `flai release <id>` is untouched by this task; it still computes and tags one item on request.

## Done when
- Accepting a story merges its branch, moves it to done, and archives it; `git log` after acceptance shows no tag and no version-file commit
- Existing accept tests updated for the narrower behavior; no test still asserts a tag or push happens inline
- flai check --strict and the markdown lint pass

## Notes
First task: the rest of the story (batch computation, publish, the board) builds on accept doing less, so this lands first.
