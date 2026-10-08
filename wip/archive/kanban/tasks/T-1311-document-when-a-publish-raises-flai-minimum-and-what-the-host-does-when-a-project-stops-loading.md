---
id: T-1311
type: task
nature: remediation
title: Document when a publish raises flai.minimum and what the host does when a project stops loading
status: done
parent: S-0321
owner: alex
created: 2026-10-08T00:23:29Z
updated: 2026-10-08T08:11:22Z
transitions:
  - to: ready
    at: 2026-10-08T08:08:42Z
    by: agent-S-0321
  - to: in-progress
    at: 2026-10-08T08:08:42Z
    by: agent-S-0321
  - to: done
    at: 2026-10-08T08:11:22Z
    by: agent-S-0321
stream: S-0321
tags: [docs]
touches: [design/system/flai-cli.md, design/system/project-manifest.md, docs/users/flai.md, docs/operators/index.md, flai/internal/manifest/manifest.go, design/system/work-hierarchy.md]
after: [T-1309, T-1310]
usage:
  source: log
  seconds: 160
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 35
      output: 11569
      cache_read: 2320546
      cache_write: 54253
      cost: 1.049
---
# T-1311 Document when a publish raises flai.minimum and what the host does when a project stops loading

## Work

Bring the documents in line with what T-1309 and T-1310 built:

- `design/system/flai-cli.md` § Versions: the host's flai and the tree. Update the rule for the raise, the publish's output, and the drain of requests in flight when a served project stops loading.
- `design/system/project-manifest.md`. Update the `flai.minimum` comment in the example and the bullet on the key.
- `docs/users/flai.md`. Update the passages that say publishing a release raises `flai.minimum`.
- `docs/operators/index.md`. Update the minimum: it rises only to a flai that can already be installed.

Waits for T-1309 and T-1310, since it describes what they built.

## Done when

- Each document above says the minimum rises no higher than the flai that publishes, and that a raise waits for a publish from the upgraded flai.
- `design/system/flai-cli.md` says that requests in flight finish when a served project stops loading.
- The `updated` date of each changed design document is bumped.
- `flai test` passes on the changed paths, the markdown lint among them.

## Notes
