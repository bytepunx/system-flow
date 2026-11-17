---
id: T-3404
type: task
nature: remediation
title: Tell whether the running flai is behind main's flai sub-project, by its build commit or else by the newest flai release tag
status: done
parent: S-0797
owner: alex
created: 2026-11-16T01:44:19Z
updated: 2026-11-16T03:14:29Z
transitions:
  - to: ready
    at: 2026-11-16T03:03:04Z
    by: agent-S-0797
  - to: in-progress
    at: 2026-11-16T03:03:04Z
    by: agent-S-0797
  - to: done
    at: 2026-11-16T03:14:29Z
    by: agent-S-0797
stream: S-0797
tags: [flai, release]
touches: [flai/internal/release/outdated.go, flai/internal/release/outdated_test.go]
usage:
  source: log
  seconds: 685
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 4897
      cache_read: 1119457
      cache_write: 41197
      cost: 0.599
---
# T-3404 Tell whether the running flai is behind main's flai sub-project, by its build commit or else by the newest flai release tag

## Work

Add a function beside `FlaiOutdated` in `flai/internal/release/outdated.go` that says whether main's flai sub-project has commits the running flai was not built from. `FlaiOutdated` compares release tags only, so it cannot see unreleased commits.

- Input: an `execx.Runner`, the main checkout's root, the flai sub-project's path from the manifest, and the running build's version and commit (`buildinfo.Get()`).
- With a commit that is an ancestor of HEAD, behind means `git log <commit>..HEAD -- <path>` lists a commit.
- With the commit `unknown`, or not in main's history, fall back to the newest `flai/v*` tag the running version is not below (`buildinfo.Below`), and compare from that tag's commit.
- With a `dev` version and no commit, or with git unable to answer, say "cannot tell". Never say behind on a guess.
- Return the running version and commit, the count of commits main's sub-project has beyond it, and main's newest such commit, so the hold's words can name them.

This task waits for none: its input is `buildinfo` and git, and the preflight task that reads it waits for it.

## Done when

- Table tests in `flai/internal/release/outdated_test.go` cover, against a real git repository in a temporary folder:
  - an ancestor commit with later `flai/` changes (behind, with the count);
  - an ancestor commit with later changes only outside `flai/` (not behind);
  - an unknown commit with a tag fallback;
  - a commit not in the history;
  - a folder that is not a repository (cannot tell).
- `flai test flai/internal/release` passes.

## Notes
