---
id: T-1314
type: task
nature: improvement
title: A check scoped to a story leaves out a narrative.state finding on another story's narrative, so a close-out records none
status: backlog
parent: S-0323
owner: alex
created: 2026-10-08T00:26:35Z
updated: 2026-10-08T00:26:35Z
transitions: []
stream: S-0323
tags: [flai]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md]
after: [T-1313]
---
# T-1314 A check scoped to a story leaves out a narrative.state finding on another story's narrative, so a close-out records none

## Work

Build the remedy T-1313's ADR decides. It waits for T-1313, because the ADR settles what is left out and what is still recorded.

- In `ScopeToStory` (`flai/internal/check/scope.go`), drop a `narrative.state` finding whose path is not the story's own narrative, beside the `item.archive` and `wip.overlap` cases, with `res.drop`, which takes back the counts, the advisory count included. Reuse what S-0318 adds to tell another open story's narrative by `repo.NarrativePath(id)`; the story waits for S-0318 for that reason.
- Keep every other finding on another story's narrative an outside note, as it is now.
- In `recordOutside` (`flai/cmd/check.go`), change the comment that names what never reaches it. In the command's help, where it names `item.archive`, say that a `narrative.state` on another story's narrative is left out too. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Reproduce I-0109 in tests:
  - `flai/internal/check/scope_test.go`: another story in progress whose narrative's `## Current state` and `## Next steps` hold the placeholder, beside the story checked. The scoped result holds no `narrative.state`, and its warnings and advisory counts do not include it.
  - The story checked with its own narrative unwritten: its `narrative.state` stays inside.
  - The other story's narrative with a broken front matter: the finding is kept as an outside note.
  - `flai/cmd/check_test.go`: `flai check --story S-nnnn --record-issues` with the first fixture records no issue for `narrative.state`.
- An unscoped `flai check` still warns the other story's `narrative.state`; show it in a test.

## Done when

- The tests above pass under `flai test` on the changed paths, and the full Go tests pass.
- `docs/users/flai-reference.md` matches the help, as `make flai-reference` writes it.
- `flai check --strict` passes.

## Notes
