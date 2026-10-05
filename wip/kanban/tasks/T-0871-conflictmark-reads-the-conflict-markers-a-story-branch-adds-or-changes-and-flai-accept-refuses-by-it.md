---
id: T-0871
type: task
nature: remediation
title: conflictmark reads the conflict markers a story branch adds or changes, and flai accept refuses by it
status: backlog
parent: S-0276
owner: alex
created: 2026-10-05T04:06:33Z
updated: 2026-10-05T04:06:33Z
transitions: []
stream: S-0276
tags: [flai]
touches: [flai/internal/conflictmark, flai/cmd/branch.go]
---
# T-0871 conflictmark reads the conflict markers a story branch adds or changes, and flai accept refuses by it

## Work

Move the branch scan out of `flai/cmd/branch.go` into `flai/internal/conflictmark`, so that the acceptance and its preview find markers by one check (S-0276's second criterion). Today `(*app).conflictMarkers` in `flai/cmd/branch.go` diffs `main...story/S-nnnn` with `--diff-filter=AM` and reads each blob with `git cat-file --batch`. Only `cmd` can reach it, and `flai/internal/preview` cannot import `cmd`.

- Add a function to `conflictmark`, for example `Branch(r execx.Runner, mainRoot, branch string) ([]string, error)`. It finds the main branch with `storygit.MainBranch` and returns the markers as `path:line`, sorted by path and then line, exactly as `conflictMarkers` does now: submodules and binary files are skipped. Check that the new imports make no cycle with `flai/internal/check`, which imports `conflictmark`.
- Make `refuseConflictMarkers` call it, and delete `conflictMarkers` from `cmd`. Keep the refusal's text as it is.
- Test the function in `conflictmark` against a git repository the test builds: one branch that adds a marker and one that is clean.

This task is first because the preview task calls the function it adds.

## Done when

- `conflictmark` exports the branch scan, and `flai/cmd/branch.go` no longer reads blobs itself.
- `flai/cmd/accept_conflict_test.go` passes unchanged.
- A `conflictmark` test covers a branch with a marker and a clean branch.
- `scripts/flai-test.sh` passes.

## Notes
