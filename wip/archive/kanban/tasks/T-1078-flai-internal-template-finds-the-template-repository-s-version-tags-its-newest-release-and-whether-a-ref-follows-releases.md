---
id: T-1078
type: task
nature: remediation
title: flai/internal/template finds the template repository's version tags, its newest release, and whether a ref follows releases
status: cancelled
parent: S-0301
owner: alex
created: 2026-10-06T22:52:37Z
updated: 2026-10-06T22:53:26Z
transitions:
  - to: cancelled
    at: 2026-10-06T22:53:26Z
    by: agent-S-0301
stream: S-0301
tags: [cli]
touches: [flai/internal/template/releases.go, flai/internal/template/releases_test.go]
---
# T-1078 flai/internal/template finds the template repository's version tags, its newest release, and whether a ref follows releases

## Work

Criteria 2 and 3 need the template repository's newest release. Add it to `flai/internal/template`, in a new `releases.go` beside `source.go`:

- `Releases(r execx.Runner, repo string)` runs `git ls-remote --symref <repo>` once through the runner and returns the remote's default branch (from `ref: refs/heads/<b> HEAD`) and its version tags: `refs/tags/vX.Y.Z` (peeled `^{}` lines folded into the tag, pre-releases and other names left out), sorted by semantic version. An error says what was asked of which repository and what to do.
- `Latest()` on the result names the highest tag and its version (`v1.0.60`, `1.0.60`), and reports none when there is no version tag. A lookup by version (`1.0.60` to `v1.0.60`) answers whether that tag exists.
- `FollowsReleases(ref)` on the result: true when the ref is empty, the default branch, or a version tag. Such a ref follows the template's releases; any other branch or commit is a deliberate choice that is used as given.
- A git source whose ref is a branch must be fetched again rather than taken from the cache (the cache is keyed by repo and ref, so `main` stayed at its first clone, 1.0.18). A tag or commit clone may be reused. Expose what callers need for that, without changing `Ensure`'s contract for others.
- Semantic version comparison: reuse `internal/buildinfo.Semver` if it fits; add no dependency.

This task waits for nothing: it is a module with no caller yet, and shares no path with T-1057.

## Done when

- [ ] `flai/internal/template/releases_test.go` covers, with a fake runner and no network: tags sorted numerically (`v1.0.7` below `v1.0.60`), peeled and pre-release tags handled, the default branch read, none on a repository without version tags, the lookup by version, and `FollowsReleases` for empty, `main`, `v1.0.18`, another branch, and a commit.
- [ ] `go test -race ./internal/template/` passes.

## Notes
- 2026-10-06T22:53:26Z: moved to cancelled: duplicate of the planner's draft for S-0301, written before I saw it; the planner's tasks are kept and edited
