---
id: T-3615
type: task
nature: remediation
title: Find the input that panics the section merge at S-0809's second stop and fix it, with a test that reproduces it
status: backlog
parent: S-0848
owner: alex
created: 2026-11-27T23:56:35Z
updated: 2026-11-27T23:56:35Z
transitions: []
stream: S-0848
tags: [flai, git]
touches: [flai/internal/storygit/mdmerge.go, flai/internal/storygit/mdmerge_test.go]
---
# T-3615 Find the input that panics the section merge at S-0809's second stop and fix it, with a test that reproduces it

## Work

I-0407's second stop was on `design/system/flai-cli.md` and `design/system/workflow.md`. Both are markdown under the design folder, which `ContinueOverResolved` merges by sections with `storygit.MergeMarkdown` (ADR-0276). That makes the section merge the likely panic: `flai/internal/storygit/mdmerge.go`, where S-0817 added the word merge (ADR-0314) two days before the instance. The panic's first frame was not recorded, so find it first.

- Rebuild the stop's three versions of the two files from git: the base (the merge base of S-0809's commit and main at S-0807's acceptance), ours (main then), and theirs (the S-0809 commit replayed after the agent's resolved stop). Look in `git log --grep S-0809` and `--grep S-0807` on main, and in the reflog of the S-0809 worktree if it is still there.
- Run `MergeMarkdown` on each triple at the revision of `mdmerge.go` S-0809's branch held, and at main's head.
- If the triples cannot be rebuilt, drive `MergeMarkdown` over three-way triples from the history of the design folder's markdown files until one panics.
- Write the panic's first frame, its cause, and the input in this task's `## Notes`, then propose the fix there before building it, as the story's goal asks.
- Add the smallest triple that panics as a test in `flai/internal/storygit/mdmerge_test.go`, and fix `mdmerge.go` so that the test passes. A triple the merge cannot settle must answer an error, which leaves the stop waiting, and never panic.

It waits for nothing. It shares no path with the guard in `sync.go`, so the two run together.

## Done when

- A test in `flai/internal/storygit/mdmerge_test.go` gives `MergeMarkdown` the input that panicked, and passes: the merge answers a merged file or an error.
- This task's notes name the function and line that panicked, the input that made it, and the fix.
- If no input could be found that panics, the notes say which inputs were tried, and no change is made to `mdmerge.go`.
- `flai test` on the two paths passes.

## Notes
