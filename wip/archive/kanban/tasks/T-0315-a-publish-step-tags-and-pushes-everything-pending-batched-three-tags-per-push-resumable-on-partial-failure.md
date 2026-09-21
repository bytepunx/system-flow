---
id: T-0315
type: task
nature: feature
title: A publish step tags and pushes everything pending, batched three tags per push, resumable on partial failure
status: done
parent: S-0087
owner: alex
created: 2026-09-21T03:32:55Z
updated: 2026-09-21T03:58:22Z
transitions:
  - to: ready
    at: 2026-09-21T03:52:40Z
    by: system-flow
  - to: in-progress
    at: 2026-09-21T03:52:40Z
    by: system-flow
  - to: done
    at: 2026-09-21T03:58:22Z
    by: system-flow
stream: S-0087
tags: []
---
# T-0315 A publish step tags and pushes everything pending, batched three tags per push, resumable on partial failure

## Work
A publish command (`flai release --pending`, or a name that reads better once tried) applies T2's per-component plans: bump each component's version file, commit, create every tag, and push main and the tags together, three tags per push (`pending.Batches`, already used by accept's old push step, per I-0026). Resumable: a tag that already exists from a partial prior run is not recreated, and a push that partially succeeded (main pushed, not all tag batches) continues from what is left rather than starting over or erroring on an existing tag.

## Done when
- One command publishes everything T2 finds pending, in one run, creating exactly the tags the plan calls for
- A batch of more than three tags goes out three at a time, branch first, as accept's old push step did
- Killing the command mid-push and running it again finishes the rest without re-tagging or erroring on what already went out
- The changelog gets one entry per component release, naming every story ID the batch bundles
- Unit tests for the resumable case and the changelog content

## Notes
Depends on T2's plan. `flai release <id>`'s existing `--apply` path (`release.Apply` + `release.Tag`) is the model for one component; this generalises it to several components in one run and adds the push, which `flai release` has never done.
