---
id: T-1060
type: task
nature: remediation
title: The template package lists a git template's version tags, names the newest, and fetches a cached branch ref again
status: done
parent: S-0301
owner: alex
created: 2026-10-06T22:50:05Z
updated: 2026-10-06T23:00:09Z
transitions:
  - to: ready
    at: 2026-10-06T22:54:34Z
    by: agent-S-0301
  - to: in-progress
    at: 2026-10-06T22:54:34Z
    by: agent-S-0301
  - to: done
    at: 2026-10-06T23:00:09Z
    by: agent-S-0301
stream: S-0301
tags: [cli, template]
touches: [flai/internal/template/source.go, flai/internal/template/tags.go, flai/internal/template/tags_test.go, flai/internal/template/source_test.go]
usage:
  source: log
  seconds: 335
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 59
      output: 249
      cache_read: 1728410
      cache_write: 71324
      cost: 0.791
---
# T-1060 The template package lists a git template's version tags, names the newest, and fetches a cached branch ref again

## Work

Criteria 2 and 3 need a "latest available version" that flai cannot work out today. `template.Resolve` and `Source.Ensure` in `flai/internal/template/source.go` clone a repository at one ref into `<cache_dir>/templates/<hash of repo@ref>`. They reuse that clone for good after that. So a project on `main`, the config's default, keeps rendering whatever `main` was on the first clone, for example 1.0.18.

- In a new `flai/internal/template/tags.go`, add `Tags(r execx.Runner, repo string)`. It reads the remote's tags with one `git ls-remote --tags --refs`, and `Latest` returns the highest one that parses as a semantic version, with or without a leading `v`, as `flai template push --tag` writes `v<version>`. Pre-release tags are left out. Reuse `release.ParseVersion` if the template package can import it without a cycle; otherwise keep a small parser here. A repository with no version tags returns none, and callers fall back to the default branch.
- Add a way to match a ref the operator gives: `1.0.60` finds the tag `v1.0.60` when no ref `1.0.60` exists.
- Read the remote's default branch in the same call (`git ls-remote --symref`, the `ref: refs/heads/<branch> HEAD` line). Add `FollowsReleases(ref)`: true when the ref is empty, the default branch, or a version tag. Such a ref follows the template's releases and resolves to the newest tag. Any other branch, or a commit, is a deliberate choice and is used as given (agent-S-0301's review: every config and manifest written so far says `main`, so `main` must follow releases for criterion 2 to hold on existing installs).
- A local template directory has no tags to list. Its version is its `template.yaml`, as now.
- In `Source.Ensure`, keep a clone made at a version tag, since a tag does not move. Fetch a clone made at a branch, or at the default branch, again before it is used. When the fetch fails because the network is down, warn and use the clone, rather than fail.
- This task waits for nothing and shares no path with the other first-layer tasks.

## Done when

- [ ] Tests in `flai/internal/template/tags_test.go` cover the following against a local bare repository:
  - tags v1.0.9, v1.0.10 and 1.0.18 give 1.0.18 as the newest;
  - a pre-release tag and a tag that is not a version are skipped;
  - a repository with no tags gives none;
  - `1.0.60` resolves to `v1.0.60`;
  - the default branch is read, and `FollowsReleases` holds for empty, the default branch, and `v1.0.18`, and not for another branch or a commit.
- [ ] A test in `flai/internal/template/source_test.go` shows that a cached clone of a branch sees a commit pushed after it was made, and a cached clone of a tag is not fetched again.
- [ ] `scripts/flai-test.sh` passes for `flai/internal/template`.

## Notes

Drafted by the planner. `docs/contributors/template.md` already warns that a branch ref goes stale in the cache. That is why the default case reverts to an old template.
