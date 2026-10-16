---
id: T-1634
type: task
nature: feature
title: flai serve merges a story's pull request on the forge when asked, by the project's method, deletes the remote story branch, and answers whether it merged or why not
status: backlog
parent: S-0397
owner: alex
created: 2026-10-16T04:39:59Z
updated: 2026-10-16T04:39:59Z
transitions: []
stream: S-0397
tags: [cli]
touches: [flai/internal/serve/forgemerge.go, flai/internal/serve/forgemerge_test.go]
---
# T-1634 flai serve merges a story's pull request on the forge when asked, by the project's method, deletes the remote story branch, and answers whether it merged or why not

## Work

- In `flai/internal/serve/forgemerge.go`, add a function that, for one story, does the following:
  - `Find`s its pull request
  - answers "already merged" when it is
  - otherwise `Merge`s it with `acceptance.merge_method`, then deletes `story/S-nnnn` on `origin`
  - answers with the merge commit, or the forge's reason for refusing
- Take requests for it from the CLI as flai serve takes `flai serve agent restart`'s, gated on the `forge` host action.
- Never log the token.
- Test it against the fake forge: a merge by each method, already merged, refused for a conflict, refused for a rule, and the host action off.
- First in its layer: S-0392's contract and S-0393's host action are on main by the time this story starts.

## Done when

flai serve merges, reports an already merged request, and refuses with the forge's reason in each case above, and the serve tests pass.

## Notes
