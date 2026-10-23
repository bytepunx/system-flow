---
id: T-2078
type: task
nature: improvement
title: flai verify's patch id leaves out the paths flai regenerates at sync, so a round after a rebase that regenerated design/issues/summary.md still reuses the tiers main's paths do not select
status: done
parent: S-0487
owner: alex
created: 2026-10-23T18:04:43Z
updated: 2026-10-23T18:20:36Z
transitions:
  - to: ready
    at: 2026-10-23T18:09:34Z
    by: agent-S-0487
  - to: in-progress
    at: 2026-10-23T18:09:35Z
    by: agent-S-0487
  - to: done
    at: 2026-10-23T18:20:36Z
    by: agent-S-0487
stream: S-0487
tags: [cli]
touches: [flai/internal/verify/paths.go, flai/internal/verify/paths_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/verify.go]
usage:
  source: log
  seconds: 661
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 61
      output: 15234
      cache_read: 3399517
      cache_write: 141192
      cost: 1.8976
---
# T-2078 flai verify's patch id leaves out the paths flai regenerates at sync, so a round after a rebase that regenerated design/issues/summary.md still reuses the tiers main's paths do not select

## Work

The cause of I-0196: `verify.BranchPatch` (`flai/internal/verify/paths.go`) takes `git patch-id --stable` of the whole diff from the merge base to `HEAD`. A story that records or closes an issue changes `design/issues/summary.md` on its branch. When main gains an issue meanwhile, the sync's rebase regenerates `summary.md` in the worktree (ADR-0098), so the branch's diff of that file differs after the rebase, the patch id differs, and `storyRun.reusable` (`flai/internal/verify/story.go`) reuses nothing: go-test runs again for a change that reaches no Go file.

The fix: the patch id covers the branch's own change without the generated paths, the files flai writes from others and regenerates at sync, `issues.Generated`'s `Path`s (`design/issues/summary.md` alone today). Their content is a function of the issue files, which stay in the diff, and of main's issue files, which `BetweenPaths` already brings to tier selection, so leaving them out loses nothing a tier's result depends on.

- Give `BranchPatch` the paths to leave out and exclude them from the `git diff` with a pathspec (`-- . ':(exclude,top)<path>'`, so it holds from any directory). When the diff without them is empty, record no patch, as today for an empty diff.
- Hand the paths in from the caller: an option on `verify.Options` that `flai/cmd/verify.go` fills from `issues.Generated`, so that `verify` need not import `issues` (no cycle either way; `issues` does not import `verify`). The MCP `verify` tool and `verify.run` go through the same options.
- Bring the doc comments on `BranchPatch`, `storyRun.reusable`, and the `patch` field, and the help text in `flai/cmd/verify.go`, in step: the patch id is of the branch's own change outside the generated paths.

The tests that reproduce the cause, before the fix fails them:

- `paths_test.go`: `BranchPatch` is the same across a rebase onto a main that changed the generated path, with the branch's copy of it regenerated after the rebase, and differs when another file of the branch changes.
- `story_test.go`: a branch that changes an issue file and the summary, verified once, then rebased onto a main that added an issue and regenerated the summary, reuses the tiers before the first that main's paths select (the markdown-like tier in the fixture), with the reuse note, and runs that tier and those after it again.

## Done when

The two tests above pass, the existing reuse tests in `story_test.go` and `paths_test.go` still pass, `flai test flai/internal/verify flai/cmd` passes, and `flai verify --help` says what the patch id leaves out.

## Notes

ADR-0156 keeps the one patch id as the gate (its alternatives reject a path-by-path comparison, confirmed on TH-0424): this task narrows what the id covers, it does not replace the gate. A round whose commit step committed recorded issues still changes the branch's own change through the issue files, and still reuses nothing, as ADR-0156's consequences say.
