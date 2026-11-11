---
id: T-3101
type: task
nature: remediation
title: release-signing.md and flai.md say a tree build links its own commit, and I-0348 is closed
status: done
parent: S-0724
owner: alex
created: 2026-11-11T08:54:58Z
updated: 2026-11-11T14:42:07Z
transitions:
  - to: ready
    at: 2026-11-11T14:41:38Z
    by: agent-S-0724
  - to: in-progress
    at: 2026-11-11T14:41:38Z
    by: agent-S-0724
  - to: done
    at: 2026-11-11T14:42:07Z
    by: agent-S-0724
stream: S-0724
tags: [go, build, docs]
touches: [design/system/release-signing.md, docs/users/flai.md, design/issues/I-0348-a-go-build-in-a-story-s-worktree-stamps-vcs-revision-of-the-main-checkout-s.md, design/issues/summary.md]
after: [T-3100]
usage:
  source: log
  seconds: 29
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2686
      cache_read: 821528
      cache_write: 25923
      cost: 0.4116
---
# T-3101 release-signing.md and flai.md say a tree build links its own commit, and I-0348 is closed

## Work

Describe the fix and close the issue.

- In `design/system/release-signing.md`, section "How flai is released": replace the sentence that a build without `ldflags` falls back to `debug.ReadBuildInfo`. Say that `scripts/flai-build.sh` links the tree's commit and date and builds with `-buildvcs=false`, because Go's VCS lookup takes only a `.git` folder and so stamps a story worktree's build with the main checkout's HEAD (I-0348). The fallback stays for any other build without `ldflags`.
- In `docs/users/flai.md`, section "Version": say that a build from a tree by `scripts/flai.sh` prints that tree's commit, a story's worktree included.
- Close I-0348 with `flai issue close I-0348 --reason`, naming the build flags, T-3093's test, and the smoke tier that runs it. That regenerates `design/issues/summary.md`.

It waits for T-3100: it describes T-3100's fix, and the issue is closed only once that fix passes.

## Done when

- `release-signing.md` and `flai.md` say what the previous bullets say, and their `updated` dates are bumped.
- I-0348 is closed with a reason, and `design/issues/summary.md` no longer lists it as open.
- The markdown lint and `flai check --strict` pass.

## Notes
