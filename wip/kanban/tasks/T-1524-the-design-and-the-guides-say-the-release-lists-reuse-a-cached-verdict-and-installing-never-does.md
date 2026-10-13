---
id: T-1524
type: task
nature: improvement
title: The design and the guides say the release lists reuse a cached verdict and installing never does
status: backlog
parent: S-0379
owner: alex
created: 2026-10-13T17:02:31Z
updated: 2026-10-13T17:02:31Z
transitions: []
stream: S-0379
tags: [flai, docs, release]
touches: [design/system/flai-cli.md, design/system/release-signing.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/runbooks/backup.md]
after: [T-1523]
---
# T-1524 The design and the guides say the release lists reuse a cached verdict and installing never does

## Work

Write down what T-1523 built. It waits for T-1523 so that the text describes the cache as built: its file, its key, and what it never stores.

- `design/system/flai-cli.md` § Verifying a release before it is installed, the bullet "The lists": a listing reuses the verdict cached for a release's `checksums.txt` and signature URLs and the trusted keys, and downloads only releases it has no verdict for. Cite S-0379.
- `design/system/release-signing.md` § How a release is verified before it is installed: the lists may answer from the cache, and an install verifies from the network every time.
- `docs/users/flai.md` § What an upgrade verifies: what a listing downloads now, and that deleting the cache file only makes the next listing slower.
- `docs/operators/settings.md` and `docs/users/flai.md`: `cache_dir` holds the release verdicts as well as the template clones.
- `docs/operators/runbooks/backup.md`: the verdicts file is fetched again like the rest of `~/.flai/cache`, and needs no backup.

## Done when

- [ ] Each document above says what the cache holds and that installing never reads it, with its `updated` date bumped.
- [ ] The markdown lint passes on the changed files.

## Notes

`docs/users/flai.md` and `design/system/flai-cli.md` are shared paths in `system-flow.yaml`, so they hold no other story.
